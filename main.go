package main

import (
	"restapi/helpers"
	"restapi/routes"
)

func main() {

	helpers.Connect()
	helpers.Init()
	routes.Routes()
}
 