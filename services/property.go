package services

import "rental-property-api/models"
import "os"
import "encoding/json"

type PropertyService struct {
	// json had the array, so we made slices of structs
	properties [] models.SourceProperty
}

func LoadJSON() *PropertyService {
	// Read the json
	data, err := os.ReadFile("data/rental_properties.json")
	if err != nil {
		// couldn't read
		return nil
	}

	// the structure of json
	var properties []models.SourceProperty
	err2 := json.Unmarshal(data, &properties)
	if err2 != nil {
		// couldn't unmarshal
		return nil
	}

	// create a struct and return its address
	return &PropertyService{
		properties: properties,
	}
}