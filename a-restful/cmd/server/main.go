package main

import (
	"a-resetful/internal/router"
)

func main() {
	app := router.Setup()
	app.Listen(":3000")
}
