package services

import "rental-property-api/models"
import "os"
import "encoding/json"

type PropertyService struct {
	// json had the array, so we made slices of structs
	properties []models.RentalProperty
}

func LoadJSON() (*PropertyService, error) {
	// Read the json
	data, err := os.ReadFile("data/rental_properties.json")
	if err != nil {
		// couldn't read
		return nil, err
	}

	// the structure of json
	var properties []models.RentalProperty
	err2 := json.Unmarshal(data, &properties)
	if err2 != nil {
		// couldn't unmarshal
		return nil, err2
	}

	// create a struct and return its address
	return &PropertyService{
		properties: properties,
	}, nil
}

func (s *PropertyService) GetProperties() []models.RentalProperty {
	// we are returning the slices of properties

	// return s.properties[:2]
	return s.properties
}

func (s *PropertyService) GetResponses() []models.Response {
	responses := make([]models.Response, 0, len(s.properties))

	for _, property := range s.properties {
		var breadcrumbs []models.Breadcrumb

		err := json.Unmarshal([]byte(property.Categories), &breadcrumbs)
		if err != nil {
			breadcrumbs = []models.Breadcrumb{}
		}

		response := models.Response{
			ID:        property.ID,
			Feed:      property.Feed,
			Published: property.Published,

			GeoInfo: models.GeoInfo{
				Breadcrumbs: breadcrumbs,
				City:        property.City,
				Country:     property.Country,
				CountryCode: property.CountryCode,
				Name:        property.Display,
				LocationID:  property.LocationID,
				Lat:         property.LonLat.Coordinates[1],
				Lon:         property.LonLat.Coordinates[0],
				State:       property.State,
				StateAbbr:   property.StateAbbr,
			},

			Property: models.Property{
				Amenities:    property.AmenityCategories,
				Name:         property.PropertyName,
				Slug:         property.PropertySlug,
				PropertyType: property.PropertyTypeCategory,
				Price:        property.USDPrice,
				ReviewScore:  property.ReviewScoreGeneral,
				StarRating:   property.StarRating,

				Counts: models.Counts{
					Bathroom:  property.BathroomCount,
					Bedroom:   property.BedroomCount,
					Reviews:   property.NumberOfReview,
					Occupancy: property.Occupancy,
				},

				Image: models.Image{
					Count:  len(property.Images),
					Images: property.Images,
				},
			},
		}

		responses = append(responses, response)
	}

	// return responses
	return responses[:2]
}
