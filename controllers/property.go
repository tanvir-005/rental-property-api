package controllers

import (
	"math"
	"rental-property-api/services"
	beego "github.com/beego/beego/v2/server/web"
	"strings"
)

// import "strconv"

type PropertyController struct {
	beego.Controller
	// controller will be able to use service
	Service *services.PropertyService
}

func (c *PropertyController) InvalidQueryError(message string) {
	c.Ctx.ResponseWriter.WriteHeader(400)

	c.Data["json"] = map[string]string{
		"Error": message,
	}

	c.ServeJSON()
}

func (c *PropertyController) Get() {
	queryParameters := c.Ctx.Request.URL.Query()

	minPrice, maxPrice, minStar, minReviewScore, minReviews, published, propertyType, feed, minBedroom := 0.0, math.MaxFloat64, 0, 0.0, 0, -1, "not-mentioned", -1, 0
	allgood := true

	if _, ok := queryParameters["min_price"]; ok {
		var err error
		minPrice, err = c.GetFloat("min_price")
		if err != nil || minPrice < 0{
			c.InvalidQueryError("Invalid min_price")
			allgood = false
		}
	} 
	
	if _, ok := queryParameters["max_price"]; ok {
		var err error
		maxPrice, err = c.GetFloat("max_price")
		if err != nil || maxPrice < minPrice{
			c.InvalidQueryError("Invalid max_price")
			allgood = false
		}
	}

	if _, ok := queryParameters["min_star_rating"]; ok {
		var err error
		minStar, err = c.GetInt("min_star_rating")
		if err != nil || minStar < 0{
			c.InvalidQueryError("invalid min_star_rating")
			allgood = false
		}
	}

	if _, ok := queryParameters["min_review_score"]; ok {
		var err error
		minReviewScore, err = c.GetFloat("min_review_score")
		if err != nil || minReviewScore < 0.0{
			c.InvalidQueryError("Invalid min_review_score")
			allgood = false
		}
	}

	if _, ok := queryParameters["min_reviews"]; ok {
		var err error	
		minReviews, err = c.GetInt("min_reviews")
		if err != nil || minReviews < 0{
			c.InvalidQueryError("Invalid min_reviews")
			allgood = false
		}
	}

	if _, ok := queryParameters["published"]; ok {
		var err error
		val, err := c.GetBool("published")
		if err != nil {
			c.InvalidQueryError("Invalid published")
			allgood = false
		}
		if val {
			published = 1
		} else {
			published = 0
		}
	}

	if _, ok := queryParameters["property_type"]; ok {
		propertyType = c.GetString("property_type")
		if propertyType == "" || (propertyType != "Hotel" && propertyType != "House" && propertyType != "Apartment" && propertyType != "Villa" && propertyType != "Resort" && propertyType != "Hostel") {
			c.InvalidQueryError("Invalid property_type")
			allgood = false
		}
	}

	if _, ok := queryParameters["feed"]; ok {
		var err error
		feed, err = c.GetInt("feed")
		if err != nil || (feed != 11 && feed != 12  && feed != 22 && feed != 24){
			c.InvalidQueryError("Invalid feed")
			allgood = false
		}
	}
	
	if _, ok := queryParameters["min_bedroom"]; ok {
		var err error
		minBedroom, err = c.GetInt("min_bedroom")
		if err != nil || minBedroom < 1{
			c.InvalidQueryError("Invalid min_bedroom")
			allgood = false
		}
	}

	var amenities []string
	if values, ok := queryParameters["amenities"]; ok {
		amenities = strings.Split(values[0], ",")
		if amenities[0] == "" {
			c.InvalidQueryError("No mentioned aminities")
			allgood = false
		}
	}

	limit := math.MaxInt
	if _, ok := queryParameters["limit"]; ok {
		var err error
		limit, err = c.GetInt("limit")
		if err != nil || limit < 1{
			c.InvalidQueryError("Invalid limit")
			allgood = false
		}
	}

	if allgood {
		response, err := c.Service.GetResponses(
			minPrice,
			maxPrice,
			minStar,
			minReviewScore,
			minReviews,
			published,
			propertyType,
			feed,
			minBedroom,
			amenities,
			limit,
		)
	
		if err != nil {
			c.Ctx.ResponseWriter.WriteHeader(400)
			c.Data["json"] = map[string]string{
				"error": err.Error(),
			}
			c.ServeJSON()
			return
		}
	
		c.Data["json"] = response
		c.ServeJSON()
	}
}

func (c *PropertyController) GetByID() {
	id := c.Ctx.Input.Param(":id")

	response, err := c.Service.GetResponseByID(id)

	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(404)

		c.Data["json"] = map[string]string{
			"Error": "The requested Property was not found!",
		}

		c.ServeJSON()
		return
	}

	c.Data["json"] = response
	c.ServeJSON()
}