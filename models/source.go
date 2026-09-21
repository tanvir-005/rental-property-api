package models

type LonLat struct {
	Coordinates [2]float64 `json:"coordinates"`
}

type PropertyNames struct {
	PropertyName string `json:"property_name"`
	PropertySlug string `json:"property_slug"`
}

type Location {
	LocationId int `json:"location_id"`
	Display string `json:"display"`
	City string `json:"city"`
	StateAbbr string `json:"state_abbr"`
	State string `json:"state"`
	CountryCode string `json:"country_code"`
	Country string `json:"country"`
	
	LonLat LonLat `json:"lonlat"`
}

type PropertyDetails struct {
	BedroomCount int `json:"bedroom_count"`
	BathroomCount int `json:"bathroom_count"`
}

type Reviews struct {
	NumberOfReview int `json:"number_of_review"`
	ReviewScoreGeneral float32 `json:"review_score_general"`
	StarRating int `json:"star_rating"`
}

type Property struct {
	Id string `json:"id"`

	PropertyNames PropertyNames
	Location Location
	PropertyDetails PropertyDetails
	Reviews Reviews
	
	Feed int `json:"feed"`
	PropertyTypeCategory string `json:"property_type_category"`
	UsdPrice float32 `json:"usd_price"`
	Occupancy int `json:"occupancy"`

	Images []string `json:"images"`
	AmenityCategories []string `json:"amenity_categories"`
	Categories string `json:"categories"`
	Published bool `json:"published"`
}