package main

import (
	"rental-property-api/controllers"
	"rental-property-api/routers"
	"rental-property-api/services"

	beego "github.com/beego/beego/v2/server/web"
)

func main() {

	service, err := services.LoadJSON()
	if err != nil {
		panic(err)
	}

	// giving controller pointer of the unmarshalled data
	controller := &controllers.PropertyController{
		Service: service,
	}

	// _ = controller // temporary to get out of error
	routers.RegisterRoutes(controller)

	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}
	beego.Run()
}
