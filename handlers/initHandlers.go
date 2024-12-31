package handlers

import (
	"httpServer/router"
)

func InitHandlers(router *router.Router) {
	router.AddRoute("GET", "/", WelcomeHandler)
}
