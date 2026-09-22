package controllers

import "rental-property-api/services"
import beego "github.com/beego/beego/v2/server/web"
import "strconv"

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
	page, err := strconv.Atoi(c.GetString("page"))
	if err != nil {
		page = 1
	}

	pageSize, err := strconv.Atoi(c.GetString("page_size"))
	if err != nil {
		pageSize = 10
	}

	response, err := c.Service.GetResponses(page, pageSize)

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
