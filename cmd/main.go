package main

import (
	"obsiTeleGo/cmd/app"
)

func main() {

	app := app.New()

	app.Log.Info("Starting obsiTeleGo")
}
