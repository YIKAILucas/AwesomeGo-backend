package main

import (
	"../src/routes"
)

func main() {
	router:=routes.MyController()
	// Listen and Server in 0.0.0.0:8080
	router.Run(":8080")

}
