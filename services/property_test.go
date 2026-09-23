package services

import (
	"reflect"
	"testing"
	"rental-property-api/models"
)

func testPropertyService(properties ...models.RentalPropertySource) *PropertyService {
	return &PropertyService{
		properties: properties,
	}
}

func testSourceProperty() models.RentalPropertySource {
	stateAbbr := "TK"

	return models.RentalPropertySource{
		ID:                   "BC-1000001",
		Feed:                 11,
		Country:              "Japan",
		CountryCode:          "JP",
		State:                "Tokyo",
		StateAbbr:            &stateAbbr,
		City:                 "Shinjuku",
		Display:              "Shinjuku, Japan",
		LocationID:           "6100001",
		PropertyName:         "Shinjuku Grand Resort",
		PropertySlug:         "shinjuku-grand-resort",
		PropertyTypeCategory: "Hotel",
		USDPrice:             150.50,
		Occupancy:            4,
		BedroomCount:         2,
		BathroomCount:        1,
		NumberOfReview:       120,
		ReviewScoreGeneral:   8.7,
		StarRating:           4,
		AmenityCategories:    []string{"WiFi", "Pool", "Parking"},
		LonLat: models.LonLat{
			Coordinates: [2]float64{139.7000, 35.6900},
		},
		Categories: `[{"LocationID":"6100001","Name":"Shinjuku","Type":"city","Slug":"shinjuku","Display":["Shinjuku","Japan"]}]`,
		Published:  true,
		Images:     []string{"image1.jpg", "image2.jpg"},
	}
}

func TestGetResponses_Transform(t *testing.T) {
	service := testPropertyService(testSourceProperty())

	got, err := service.GetResponses(
		0,    // minPrice
		1000, // maxPrice
		0,    // minStar
		0,    // minReviewScore
		0,    // minReviews
		-1,   // published: ignored
		"not-mentioned",
		-1,  // feed: ignored
		0,   // minBedroom
		nil, // amenities
		10,  // limit
	)
	if err != nil {
		t.Fatalf("GetResponses() returned error: %v", err)
	}

	if got.Result.Count != 1 {
		t.Fatalf("Count = %d, want 1", got.Result.Count)
	}

	item := got.Result.Items[0]

	if item.ID != "BC-1000001" {
		t.Errorf("ID = %q, want %q", item.ID, "BC-1000001")
	}
	if item.GeoInfo.Lat != 35.6900 {
		t.Errorf("Lat = %v, want %v", item.GeoInfo.Lat, 35.6900)
	}
	if item.GeoInfo.Lon != 139.7000 {
		t.Errorf("Lon = %v, want %v", item.GeoInfo.Lon, 139.7000)
	}
	if len(item.GeoInfo.Breadcrumbs) != 1 {
		t.Fatalf("Breadcrumbs length = %d, want 1", len(item.GeoInfo.Breadcrumbs))
	}
	if item.GeoInfo.Breadcrumbs[0].Name != "Shinjuku" {
		t.Errorf("Breadcrumb name = %q, want %q",
			item.GeoInfo.Breadcrumbs[0].Name, "Shinjuku")
	}
	if item.Property.Image.Count != 2 {
		t.Errorf("Image.Count = %d, want 2", item.Property.Image.Count)
	}
	if !reflect.DeepEqual(item.Property.Image.Images, []string{"image1.jpg", "image2.jpg"}) {
		t.Errorf("Images = %v, want [image1.jpg image2.jpg]", item.Property.Image.Images)
	}
}

func TestGetResponses_FilterAND(t *testing.T) {
	base := testSourceProperty()

	tests := []struct {
		name         string
		minPrice     float64
		maxPrice     float64
		minStar      int
		minReview    float64
		minReviews   int
		published    int
		propertyType string
		feed         int
		minBedroom   int
		wantCount    int
	}{
		{
			name:         "feed and published",
			minPrice:     0,
			maxPrice:     1000,
			published:    1,
			propertyType: "not-mentioned",
			feed:         11,
			wantCount:    1,
		},
		{
			name:         "price range",
			minPrice:     150.50,
			maxPrice:     150.50,
			published:    -1,
			propertyType: "not-mentioned",
			feed:         -1,
			wantCount:    1,
		},
		{
			name:         "price range excludes property",
			minPrice:     151,
			maxPrice:     1000,
			published:    -1,
			propertyType: "not-mentioned",
			feed:         -1,
			wantCount:    0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := testPropertyService(base)

			got, err := service.GetResponses(
				tt.minPrice,
				tt.maxPrice,
				tt.minStar,
				tt.minReview,
				tt.minReviews,
				tt.published,
				tt.propertyType,
				tt.feed,
				tt.minBedroom,
				nil,
				10,
			)
			if err != nil {
				t.Fatalf("GetResponses() returned error: %v", err)
			}

			if got.Result.Count != tt.wantCount {
				t.Errorf("Count = %d, want %d", got.Result.Count, tt.wantCount)
			}
		})
	}
}

func TestGetResponses_FilterAmenitiesOR(t *testing.T) {
	property := testSourceProperty()

	tests := []struct {
		name      string
		amenities []string
		wantCount int
	}{
		{
			name:      "matches one of multiple requested amenities",
			amenities: []string{"Gym", "WiFi"},
			wantCount: 1,
		},
		{
			name:      "no requested amenity matches",
			amenities: []string{"Gym", "Spa"},
			wantCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := testPropertyService(property)

			got, err := service.GetResponses(
				0, 1000, 0, 0, 0,
				-1, "not-mentioned", -1, 0,
				tt.amenities,
				10,
			)
			if err != nil {
				t.Fatalf("GetResponses() returned error: %v", err)
			}

			if got.Result.Count != tt.wantCount {
				t.Errorf("Count = %d, want %d", got.Result.Count, tt.wantCount)
			}
		})
	}
}

func TestGetResponses_FilterCombined(t *testing.T) {
	property := testSourceProperty()
	service := testPropertyService(property)

	got, err := service.GetResponses(
		100, // minPrice
		200, // maxPrice
		4,   // minStar
		8.0, // minReviewScore
		100, // minReviews
		1,   // published
		"Hotel",
		11,                      // feed
		2,                       // minBedroom
		[]string{"Gym", "WiFi"}, // OR within amenities
		10,
	)
	if err != nil {
		t.Fatalf("GetResponses() returned error: %v", err)
	}

	if got.Result.Count != 1 {
		t.Fatalf("Count = %d, want 1", got.Result.Count)
	}
}

func TestGetResponses_EmptyResult(t *testing.T) {
	service := testPropertyService(testSourceProperty())

	got, err := service.GetResponses(
		1000, 2000, 0, 0, 0,
		-1, "not-mentioned", -1, 0,
		nil,
		10,
	)
	if err != nil {
		t.Fatalf("GetResponses() returned error: %v", err)
	}

	if got.Result.Count != 0 {
		t.Errorf("Count = %d, want 0", got.Result.Count)
	}

	if got.Result.Items == nil {
		t.Error("Items is nil, want an empty non-nil slice")
	}

	if len(got.Result.Items) != 0 {
		t.Errorf("len(Items) = %d, want 0", len(got.Result.Items))
	}
}

func TestGetResponseByID(t *testing.T) {
	service := testPropertyService(testSourceProperty())

	t.Run("found", func(t *testing.T) {
		got, err := service.GetResponseByID("BC-1000001")
		if err != nil {
			t.Fatalf("GetResponseByID() returned error: %v", err)
		}

		if got.ID != "BC-1000001" {
			t.Errorf("ID = %q, want %q", got.ID, "BC-1000001")
		}
		if got.GeoInfo.Lat != 35.6900 {
			t.Errorf("Lat = %v, want %v", got.GeoInfo.Lat, 35.6900)
		}
		if got.Property.Image.Count != 2 {
			t.Errorf("Image.Count = %d, want 2", got.Property.Image.Count)
		}
	})

	t.Run("not found", func(t *testing.T) {
		_, err := service.GetResponseByID("does-not-exist")
		if err == nil {
			t.Fatal("GetResponseByID() returned nil error, want error")
		}
	})
}
