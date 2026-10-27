package main

import (
	"eqray-eprint/bootstrap"
)

func main() {
	app := bootstrap.Boot()

	app.Start()
}
