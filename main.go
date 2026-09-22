package main

import (
	beego "github.com/beego/beego/v2/server/web"
	_ "rental-property-api/routers"
	"rental-property-api/services"
	"rental-property-api/controllers"
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

	_ = controller // temporary to get out of error

	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}
	beego.Run()
}
