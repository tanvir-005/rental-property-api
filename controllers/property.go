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

// func (c *PropertyController) Get() {
// 	responses := c.Service.GetResponses()

// 	c.Data["json"] = responses
// 	c.ServeJSON()
// }

func (c *PropertyController) Get() {
	queryParameters := c.Ctx.Request.URL.Query()

	minPrice, maxPrice, minStar, minReviewScore, minReviews, published, propertyType, feed, minBedroom := 0.0, math.MaxFloat64, 0, 0.0, 0, -1, "not-mentioned", -1, 0

	if _, ok := queryParameters["min_price"]; ok {
		minPrice, _ = c.GetFloat("min_price")
	} 
	if _, ok := queryParameters["max_price"]; ok {
		maxPrice, _ = c.GetFloat("max_price")
	}
	if _, ok := queryParameters["min_star_rating"]; ok {
		minStar, _ = c.GetInt("min_star_rating")
	}
	if _, ok := queryParameters["min_review_score"]; ok {
		minReviewScore, _ = c.GetFloat("min_review_score")
	}
	if _, ok := queryParameters["min_reviews"]; ok {
		minReviews, _ = c.GetInt("min_reviews")
	}
	if _, ok := queryParameters["published"]; ok {
		val, _ := c.GetBool("published")
		if val {
			published = 1
		} else {
			published = 0
		}
	}
	if _, ok := queryParameters["property_type"]; ok {
		propertyType = c.GetString("property_type")
	}
	if _, ok := queryParameters["feed"]; ok {
		feed, _ = c.GetInt("feed")
	}
	if _, ok := queryParameters["min_bedroom"]; ok {
		minBedroom, _ = c.GetInt("min_bedroom")
	}

	var amenities []string
	if values, ok := queryParameters["amenities"]; ok {
		amenities = strings.Split(values[0], ",")
	}

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
