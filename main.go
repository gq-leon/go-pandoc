package main

import (
	"fmt"
	"log"
	"net/http"

	"pandoc/pkg/logging"
	"pandoc/routers"

	_ "pandoc/pkg/utils"
)

func init() {
	logging.Setup()
}

func main() {
	initRouter := routers.InitRouter()
	endPoint := fmt.Sprintf(":%d", 8080)

	server := &http.Server{
		Addr:    endPoint,
		Handler: initRouter,
	}

	log.Printf("[info] start http server listening %s", endPoint)
	log.Fatal(server.ListenAndServe())
}
