package controllers

import "rental-property-api/services"
import beego "github.com/beego/beego/v2/server/web"

type PropertyController struct {
	beego.Controller
	// controller will be able to use service
	Service *services.PropertyService
}

func (c *PropertyController) Get() {
	responses := c.Service.GetResponses()

	c.Data["json"] = responses
	c.ServeJSON()
}
