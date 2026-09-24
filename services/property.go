package services

import (
	"encoding/json"
	"os"
	"rental-property-api/models"
	"errors"
)

type PropertyService struct {
	// json had the array, so we made slices of structs
	properties []models.RentalPropertySource
}

func LoadJSON(filepath string) (*PropertyService, error) {
	// Read the json
	data, err := os.ReadFile(filepath)
	if err != nil {
		// couldn't read
		return nil, err
	}

	// the structure of json
	var properties []models.RentalPropertySource
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

func buildResponse(property models.RentalPropertySource) models.RentalProperty {
	var breadcrumbs []models.Breadcrumb

	err := json.Unmarshal([]byte(property.Categories), &breadcrumbs)
	if err != nil {
		breadcrumbs = []models.Breadcrumb{}
	}

	response := models.RentalProperty{
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
	return response
}

func blocker(
	property models.RentalPropertySource,
	minPrice float64,
	maxPrice float64,
	minStar int,
	minReviewScore float64,
	minReviews int,
	published int,
	propertyType string,
	feed int,
	minBedroom int,
	amenities [] string,
) bool {
	// we will decide if we'd block this property or not
	// initially we're saying we won't block
	block := false

	if property.USDPrice < minPrice {
		block = true
	}
	if property.USDPrice > maxPrice {
		block = true
	}
	if property.StarRating < minStar {
		block = true
	}
	if property.ReviewScoreGeneral < minReviewScore {
		block = true
	}

	
	if property.NumberOfReview < minReviews {
		block = true
	}

	// argument published is int while property.published is bool
	var p bool = true 
	if published == 0 {
		p = false 
	}
	if (published != -1) && (property.Published != p) {
		block = true
	}
	
	if (propertyType != "not-mentioned") && (property.PropertyTypeCategory != propertyType) {
		block = true
	}
	
	if feed != -1 && (property.Feed != feed) {
		block = true
	}

	if property.BedroomCount < minBedroom {
		block = true
	}
	
	
	// amenities
	exists := make(map[string]bool)
	for _, aminity := range amenities {
		exists[aminity] = true
	}
	found := false
	for _, propertyAminity := range property.AmenityCategories {
		if exists[propertyAminity] == true {
			found = true 
			break
		}
	}
	if !found && len(amenities) > 0{
		block = true
	}

	return block
}

func (s *PropertyService) GetResponses(
	minPrice float64,
	maxPrice float64,
	minStar int,
	minReviewScore float64,
	minReviews int,
	published int,
	propertyType string,
	feed int,
	minBedroom int,
	amenities [] string,
	limit int,
) (models.RentalProperties, error) {

	totalProperty := len(s.properties)
	responses := make([]models.RentalProperty, 0)
	taken := 0
	for i := 0; i < totalProperty; i++ {
		property := s.properties[i]

		block := blocker(property, minPrice, maxPrice, minStar, minReviewScore, minReviews, published, propertyType, feed, minBedroom, amenities)
		// if block is true by any of the conditions, we should not include this property to response
		if block {
			continue
		}
		
		response := buildResponse(property)

		responses = append(responses, response)
		taken++

		if taken == limit {
			break
		}
	}

	return models.RentalProperties{
		Result: models.Result{
			Count: len(responses),
			Items: responses,
		},
	}, nil
}

func (s *PropertyService) GetResponseByID(id string) (models.RentalProperty, error) {
	for _, property := range s.properties {
		if property.ID == id {
			response := buildResponse(property)
			return response, nil
		}
	}
	return models.RentalProperty{}, errors.New("Property not found")
}