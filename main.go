package main

import (
	"github.com/defskela/httpServer/router"
	"github.com/defskela/httpServer/server"
)

func main() {
	router := router.NewRouter()
	server.StartServ(router)
}
