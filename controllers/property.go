package controllers

import "rental-property-api/services"

type PropertyController struct {
	// controller will be able to use service
	Service *services.PropertyService
}