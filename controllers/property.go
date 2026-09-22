package controllers

import "rental-property-api/services"
import beego "github.com/beego/beego/v2/server/web"

type PropertyController struct {
	beego.Controller
	// controller will be able to use service
	Service *services.PropertyService
}

func (c *PropertyController) Get() {
	properties := c.Service.GetProperties()

	c.Data["json"] = properties
	c.ServeJSON()
}