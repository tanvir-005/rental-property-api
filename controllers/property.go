package controllers

import "rental-property-api/services"
import beego "github.com/beego/beego/v2/server/web"
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
	response, err := c.Service.GetResponses()

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
