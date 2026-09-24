// Package services tests the business logic behind rental property filtering,
// transformation, lookup, and JSON loading.
//
// This file exercises the PropertyService with handcrafted sample data so the
// behavior is checked without depending on the real JSON dataset.
package services

import (
	"math"
	"os"
	"reflect"
	"testing"

	"rental-property-api/models"
)

// ------------------------------------------------------------
// Test helpers and sample fixture data
// ------------------------------------------------------------
// These helpers build a small, deterministic dataset that is easy to reason
// about when verifying filtering and transformation rules.

func strPtr(s string) *string {
	return &s
}

func testPropertyService(properties ...models.RentalPropertySource) *PropertyService {
	return &PropertyService{
		properties: properties,
	}
}

// testSourceProperty creates the primary dataset record used by most tests.
// It represents a published Hotel with moderate pricing and strong review
// metrics, making it a stable baseline for comparison and filtering checks.
func testSourceProperty() models.RentalPropertySource {
	return models.RentalPropertySource{
		ID:                   "BC-1000001",
		Feed:                 11,
		Country:              "United States",
		CountryCode:          "US",
		State:                "California",
		StateAbbr:            strPtr("CA"),
		City:                 "San Francisco",
		Display:              "San Francisco, CA",
		LocationID:           "loc-1",
		PropertyName:         "Ocean View Hotel",
		PropertySlug:         "ocean-view-hotel",
		PropertyTypeCategory: "Hotel",
		USDPrice:             150.50,
		Occupancy:            2,
		BedroomCount:         2,
		BathroomCount:        2,
		NumberOfReview:       120,
		ReviewScoreGeneral:   8.7,
		StarRating:           4,
		AmenityCategories:    []string{"WiFi", "Pool", "Parking"},
		LonLat:               models.LonLat{Coordinates: [2]float64{-122.4194, 37.7749}},
		Categories:           `[{"LocationID":"loc-1","Name":"San Francisco","Type":"City","Slug":"san-francisco","Display":["San Francisco","CA","US"]}]`,
		Published:            true,
		Images:               []string{"img1.jpg", "img2.jpg"},
	}
}

// testSecondSourceProperty creates a second record with different attributes to
// ensure filters such as price, published state, property type, and amenity
// matching behave correctly across multiple values.
func testSecondSourceProperty() models.RentalPropertySource {
	return models.RentalPropertySource{
		ID:                   "BC-1000002",
		Feed:                 12,
		Country:              "United States",
		CountryCode:          "US",
		State:                "New York",
		StateAbbr:            strPtr("NY"),
		City:                 "New York",
		Display:              "New York, NY",
		LocationID:           "loc-2",
		PropertyName:         "City Lite Apartments",
		PropertySlug:         "city-lite-apartments",
		PropertyTypeCategory: "Apartment",
		USDPrice:             300,
		Occupancy:            2,
		BedroomCount:         1,
		BathroomCount:        1,
		NumberOfReview:       20,
		ReviewScoreGeneral:   5.5,
		StarRating:           2,
		AmenityCategories:    []string{"Gym", "Spa"},
		LonLat:               models.LonLat{Coordinates: [2]float64{-73.935242, 40.73061}},
		Categories:           `[{"LocationID":"loc-2","Name":"New York","Type":"City","Slug":"new-york","Display":["New York","NY","US"]}]`,
		Published:            false,
		Images:               []string{"img3.jpg"},
	}
}

// testThirdSourceProperty creates a premium Villa record that differs in both
// price and amenities, allowing tests to validate boundary conditions and OR
// matching logic.
func testThirdSourceProperty() models.RentalPropertySource {
	return models.RentalPropertySource{
		ID:                   "BC-1000003",
		Feed:                 24,
		Country:              "United States",
		CountryCode:          "US",
		State:                "Florida",
		StateAbbr:            strPtr("FL"),
		City:                 "Miami",
		Display:              "Miami, FL",
		LocationID:           "loc-3",
		PropertyName:         "Palm Villa",
		PropertySlug:         "palm-villa",
		PropertyTypeCategory: "Villa",
		USDPrice:             500,
		Occupancy:            6,
		BedroomCount:         4,
		BathroomCount:        3,
		NumberOfReview:       250,
		ReviewScoreGeneral:   9.2,
		StarRating:           5,
		AmenityCategories:    []string{"Parking", "Garden"},
		LonLat:               models.LonLat{Coordinates: [2]float64{-80.19179, 25.76168}},
		Categories:           `[{"LocationID":"loc-3","Name":"Miami","Type":"City","Slug":"miami","Display":["Miami","FL","US"]}]`,
		Published:            true,
		Images:               []string{"img4.jpg", "img5.jpg", "img6.jpg"},
	}
}

func absDiff(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

func smallDelta(a, b float64) bool {
	return absDiff(a-b) < 1e-9
}

// ------------------------------------------------------------
// Transformation tests
// ------------------------------------------------------------
// These tests verify that raw source data is converted into the API response
// format with the correct nested values, labels, and counts.

func TestGetResponses_Transform(t *testing.T) {
	service := testPropertyService(testSourceProperty())

	result, err := service.GetResponses(0, math.MaxFloat64, 0, 0, 0, -1, "not-mentioned", -1, 0, nil, 10)
	if err != nil {
		t.Fatalf("GetResponses returned error: %v", err)
	}
	if result.Result.Count != 1 {
		t.Fatalf("expected 1 item, got %d", result.Result.Count)
	}

	item := result.Result.Items[0]
	if item.ID != "BC-1000001" {
		t.Fatalf("expected ID BC-1000001, got %s", item.ID)
	}
	if item.Feed != 11 {
		t.Fatalf("expected Feed 11, got %d", item.Feed)
	}
	if !item.Published {
		t.Fatal("expected Published true")
	}
	if item.GeoInfo.City != "San Francisco" {
		t.Fatalf("expected city San Francisco, got %s", item.GeoInfo.City)
	}
	if item.GeoInfo.Country != "United States" {
		t.Fatalf("expected country United States, got %s", item.GeoInfo.Country)
	}
	if item.GeoInfo.CountryCode != "US" {
		t.Fatalf("expected country code US, got %s", item.GeoInfo.CountryCode)
	}
	if item.GeoInfo.State != "California" {
		t.Fatalf("expected state California, got %s", item.GeoInfo.State)
	}
	if item.GeoInfo.LocationID != "loc-1" {
		t.Fatalf("expected location id loc-1, got %s", item.GeoInfo.LocationID)
	}
	if len(item.GeoInfo.Breadcrumbs) != 1 {
		t.Fatalf("expected 1 breadcrumb, got %d", len(item.GeoInfo.Breadcrumbs))
	}
	if item.GeoInfo.Breadcrumbs[0].Name != "San Francisco" {
		t.Fatalf("expected breadcrumb name San Francisco, got %s", item.GeoInfo.Breadcrumbs[0].Name)
	}
	if item.GeoInfo.Breadcrumbs[0].Type != "City" {
		t.Fatalf("expected breadcrumb type City, got %s", item.GeoInfo.Breadcrumbs[0].Type)
	}
	if !smallDelta(item.GeoInfo.Lat, 37.7749) {
		t.Fatalf("expected lat 37.7749, got %f", item.GeoInfo.Lat)
	}
	if !smallDelta(item.GeoInfo.Lon, -122.4194) {
		t.Fatalf("expected lon -122.4194, got %f", item.GeoInfo.Lon)
	}
	if item.Property.Name != "Ocean View Hotel" {
		t.Fatalf("expected property name Ocean View Hotel, got %s", item.Property.Name)
	}
	if item.Property.PropertyType != "Hotel" {
		t.Fatalf("expected property type Hotel, got %s", item.Property.PropertyType)
	}
	if !smallDelta(item.Property.Price, 150.50) {
		t.Fatalf("expected price 150.5, got %f", item.Property.Price)
	}
	if !smallDelta(item.Property.ReviewScore, 8.7) {
		t.Fatalf("expected review score 8.7, got %f", item.Property.ReviewScore)
	}
	if item.Property.StarRating != 4 {
		t.Fatalf("expected star rating 4, got %d", item.Property.StarRating)
	}
	if item.Property.Counts.Bedroom != 2 {
		t.Fatalf("expected bedroom count 2, got %d", item.Property.Counts.Bedroom)
	}
	if item.Property.Counts.Bathroom != 2 {
		t.Fatalf("expected bathroom count 2, got %d", item.Property.Counts.Bathroom)
	}
	if item.Property.Counts.Reviews != 120 {
		t.Fatalf("expected review count 120, got %d", item.Property.Counts.Reviews)
	}
	if item.Property.Counts.Occupancy != 2 {
		t.Fatalf("expected occupancy 2, got %d", item.Property.Counts.Occupancy)
	}
	if !reflect.DeepEqual(item.Property.Amenities, []string{"WiFi", "Pool", "Parking"}) {
		t.Fatalf("expected amenities [WiFi Pool Parking], got %#v", item.Property.Amenities)
	}
	if item.Property.Image.Count != 2 {
		t.Fatalf("expected image count 2, got %d", item.Property.Image.Count)
	}
	if !reflect.DeepEqual(item.Property.Image.Images, []string{"img1.jpg", "img2.jpg"}) {
		t.Fatalf("expected image list [img1.jpg img2.jpg], got %#v", item.Property.Image.Images)
	}
}

// ------------------------------------------------------------
// Scalar AND filter tests
// ------------------------------------------------------------
// Each case checks that filters are combined using AND semantics: a property must
// satisfy every provided condition to appear in the result set.

func TestGetResponses_FilterAND(t *testing.T) {
	properties := []models.RentalPropertySource{
		testSourceProperty(),
		testSecondSourceProperty(),
		testThirdSourceProperty(),
	}

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
		{name: "minimum price boundary", minPrice: 150.50, maxPrice: math.MaxFloat64, published: -1, propertyType: "not-mentioned", feed: -1, wantCount: 3},
		{name: "maximum price boundary", minPrice: 0, maxPrice: 150.50, published: -1, propertyType: "not-mentioned", feed: -1, wantCount: 1},
		{name: "price range", minPrice: 100, maxPrice: 350, published: -1, propertyType: "not-mentioned", feed: -1, wantCount: 2},
		{name: "minimum star rating", minPrice: 0, maxPrice: math.MaxFloat64, minStar: 4, published: -1, propertyType: "not-mentioned", feed: -1, wantCount: 2},
		{name: "minimum review score", minPrice: 0, maxPrice: math.MaxFloat64, minReview: 8, published: -1, propertyType: "not-mentioned", feed: -1, wantCount: 2},
		{name: "minimum reviews", minPrice: 0, maxPrice: math.MaxFloat64, minReviews: 100, published: -1, propertyType: "not-mentioned", feed: -1, wantCount: 2},
		{name: "published true", minPrice: 0, maxPrice: math.MaxFloat64, published: 1, propertyType: "not-mentioned", feed: -1, wantCount: 2},
		{name: "published false", minPrice: 0, maxPrice: math.MaxFloat64, published: 0, propertyType: "not-mentioned", feed: -1, wantCount: 1},
		{name: "property type", minPrice: 0, maxPrice: math.MaxFloat64, published: -1, propertyType: "Hotel", feed: -1, wantCount: 1},
		{name: "feed", minPrice: 0, maxPrice: math.MaxFloat64, published: -1, propertyType: "not-mentioned", feed: 11, wantCount: 1},
		{name: "minimum bedroom", minPrice: 0, maxPrice: math.MaxFloat64, minBedroom: 2, published: -1, propertyType: "not-mentioned", feed: -1, wantCount: 2},
		{name: "multiple scalar filters use AND", minPrice: 100, maxPrice: 200, minStar: 4, minReview: 8, minReviews: 100, published: 1, propertyType: "Hotel", feed: 11, minBedroom: 2, wantCount: 1},
		{name: "no property satisfies filter", minPrice: 1000, maxPrice: 2000, published: -1, propertyType: "not-mentioned", feed: -1, wantCount: 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := testPropertyService(properties...)
			propertyType := tt.propertyType
			if propertyType == "" {
				propertyType = "not-mentioned"
			}
			feed := tt.feed
			if feed == 0 {
				feed = -1
			}
			published := tt.published
			if published == 0 && tt.name != "published false" {
				published = -1
			}
			maxPrice := tt.maxPrice
			if maxPrice == 0 {
				maxPrice = math.MaxFloat64
			}
			result, err := service.GetResponses(
				tt.minPrice,
				maxPrice,
				tt.minStar,
				tt.minReview,
				tt.minReviews,
				published,
				propertyType,
				feed,
				tt.minBedroom,
				nil,
				100,
			)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if result.Result.Count != tt.wantCount {
				t.Fatalf("expected %d items, got %d", tt.wantCount, result.Result.Count)
			}
		})
	}
}

// ------------------------------------------------------------
// Amenities OR filter tests
// ------------------------------------------------------------
// Amenity checks are intentionally OR-based: a property should match if it has
// any of the requested amenities, while preserving case sensitivity rules.

func TestGetResponses_FilterAmenitiesOR(t *testing.T) {
	properties := []models.RentalPropertySource{
		testSourceProperty(),
		testSecondSourceProperty(),
		testThirdSourceProperty(),
	}

	tests := []struct {
		name      string
		amenities []string
		wantCount int
	}{
		{name: "first requested amenity matches", amenities: []string{"WiFi", "Spa"}, wantCount: 2},
		{name: "second requested amenity matches", amenities: []string{"Gym", "WiFi"}, wantCount: 2},
		{name: "none match", amenities: []string{"Gym", "Spa", "Garden"}, wantCount: 2},
		{name: "case sensitive", amenities: []string{"wifi"}, wantCount: 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := testPropertyService(properties...)
			result, err := service.GetResponses(0, math.MaxFloat64, 0, 0, 0, -1, "not-mentioned", -1, 0, tt.amenities, 100)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if result.Result.Count != tt.wantCount {
				t.Fatalf("expected %d items, got %d", tt.wantCount, result.Result.Count)
			}
		})
	}
}

// ------------------------------------------------------------
// Mixed filter tests
// ------------------------------------------------------------
// These test the combination of a scalar filter (such as feed) with an amenity
// OR condition, ensuring the overall logic still behaves as expected.

func TestGetResponses_FilterCombined(t *testing.T) {
	service := testPropertyService(
		testSourceProperty(),
		testSecondSourceProperty(),
		testThirdSourceProperty(),
	)

	result, err := service.GetResponses(0, math.MaxFloat64, 0, 0, 0, -1, "not-mentioned", 11, 0, []string{"WiFi", "Parking"}, 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Result.Count != 1 {
		t.Fatalf("expected 1 item for feed==11 with WiFi or Parking, got %d", result.Result.Count)
	}
	if result.Result.Items[0].ID != "BC-1000001" {
		t.Fatalf("expected first property to match, got %s", result.Result.Items[0].ID)
	}

	result, err = service.GetResponses(0, math.MaxFloat64, 0, 0, 0, -1, "not-mentioned", 11, 0, []string{"Spa"}, 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Result.Count != 0 {
		t.Fatalf("expected zero results for feed==11 and non-matching amenity, got %d", result.Result.Count)
	}
}

// ------------------------------------------------------------
// Result limiting tests
// ------------------------------------------------------------
// These verify that limit values cut the result set correctly without changing
// the ordering or filter behavior.

func TestGetResponses_Limit(t *testing.T) {
	properties := []models.RentalPropertySource{
		testSourceProperty(),
		testSecondSourceProperty(),
		testThirdSourceProperty(),
	}

	tests := []struct {
		name      string
		limit     int
		wantCount int
	}{
		{name: "limit one", limit: 1, wantCount: 1},
		{name: "limit two", limit: 2, wantCount: 2},
		{name: "limit greater than result", limit: 10, wantCount: 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := testPropertyService(properties...)
			result, err := service.GetResponses(0, math.MaxFloat64, 0, 0, 0, -1, "not-mentioned", -1, 0, nil, tt.limit)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if result.Result.Count != tt.wantCount {
				t.Fatalf("expected %d items, got %d", tt.wantCount, result.Result.Count)
			}
		})
	}
}

// ------------------------------------------------------------
// Empty result handling
// ------------------------------------------------------------
// This checks that no-match searches return an empty result set rather than a
// nil or partially populated response.

func TestGetResponses_EmptyResult(t *testing.T) {
	service := testPropertyService(testSourceProperty(), testSecondSourceProperty(), testThirdSourceProperty())

	result, err := service.GetResponses(0, math.MaxFloat64, 0, 0, 0, -1, "not-mentioned", -1, 0, []string{"NoSuchAmenity"}, 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Result.Count != 0 {
		t.Fatalf("expected 0 items, got %d", result.Result.Count)
	}
	if result.Result.Items == nil {
		t.Fatal("expected Items slice to be initialized")
	}
	if len(result.Result.Items) != 0 {
		t.Fatalf("expected len(Items)==0, got %d", len(result.Result.Items))
	}
}

// ------------------------------------------------------------
// ID lookup tests
// ------------------------------------------------------------
// These validate the retrieval path for a single property by its unique ID, as
// well as the expected error case when the ID is missing.

func TestGetResponseByID(t *testing.T) {
	service := testPropertyService(
		testSourceProperty(),
		testSecondSourceProperty(),
		testThirdSourceProperty(),
	)

	t.Run("found", func(t *testing.T) {
		item, err := service.GetResponseByID("BC-1000001")
		if err != nil {
			t.Fatalf("expected property to be found, got error: %v", err)
		}
		if item.ID != "BC-1000001" {
			t.Fatalf("expected ID BC-1000001, got %s", item.ID)
		}
		if item.Property.Name != "Ocean View Hotel" {
			t.Fatalf("expected property name Ocean View Hotel, got %s", item.Property.Name)
		}
	})

	t.Run("not found", func(t *testing.T) {
		_, err := service.GetResponseByID("missing-id")
		if err == nil {
			t.Fatal("expected an error for a missing ID")
		}
	})
}

// ------------------------------------------------------------
// JSON loading tests
// ------------------------------------------------------------
// These tests cover successful parsing, missing file handling, and malformed
// JSON detection for the file-based loader used at startup.

func TestLoadJSON(t *testing.T) {
	t.Run("valid JSON", func(t *testing.T) {
		dir := t.TempDir()
		path := dir + "/properties.json"

		payload := `[
			{"id":"BC-1000001","feed":11,"country":"United States","country_code":"US","state":"California","state_abbr":"CA","city":"San Francisco","display":"San Francisco, CA","location_id":"loc-1","property_name":"Ocean View Hotel","property_slug":"ocean-view-hotel","property_type_category":"Hotel","usd_price":150.5,"occupancy":2,"bedroom_count":2,"bathroom_count":2,"number_of_review":120,"review_score_general":8.7,"star_rating":4,"amenity_categories":["WiFi","Pool","Parking"],"lonlat":{"coordinates":[-122.4194,37.7749]},"categories":"[]","published":true,"images":["img1.jpg","img2.jpg"]},
			{"id":"BC-1000002","feed":12,"country":"United States","country_code":"US","state":"New York","state_abbr":"NY","city":"New York","display":"New York, NY","location_id":"loc-2","property_name":"City Lite Apartments","property_slug":"city-lite-apartments","property_type_category":"Apartment","usd_price":300,"occupancy":2,"bedroom_count":1,"bathroom_count":1,"number_of_review":20,"review_score_general":5.5,"star_rating":2,"amenity_categories":["Gym","Spa"],"lonlat":{"coordinates":[-73.935242,40.73061]},"categories":"[]","published":false,"images":["img3.jpg"]}
		]`
		if err := os.WriteFile(path, []byte(payload), 0o600); err != nil {
			t.Fatalf("failed to write temp json: %v", err)
		}

		service, err := LoadJSON(path)
		if err != nil {
			t.Fatalf("LoadJSON returned error for valid file: %v", err)
		}
		if len(service.properties) != 2 {
			t.Fatalf("expected 2 properties, got %d", len(service.properties))
		}
	})

	t.Run("missing file", func(t *testing.T) {
		_, err := LoadJSON("/tmp/definitely-does-not-exist.json")
		if err == nil {
			t.Fatal("expected an error for missing file")
		}
	})

	t.Run("invalid JSON", func(t *testing.T) {
		dir := t.TempDir()
		path := dir + "/bad.json"
		if err := os.WriteFile(path, []byte("{not valid json"), 0o600); err != nil {
			t.Fatalf("failed to write invalid json: %v", err)
		}

		_, err := LoadJSON(path)
		if err == nil {
			t.Fatal("expected an error for invalid json")
		}
	})
}
