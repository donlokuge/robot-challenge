// @title           Robot Warehouse API
// @version         1.0
// @description     REST API to control warehouse robots
// @host            localhost:3000
// @BasePath        /api/v1
package main

import (
	"a-resetful/internal/router"
)

func main() {
	app := router.Setup()
	app.Listen(":3000")
}
