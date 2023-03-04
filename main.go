package main

import (
	"log"
	"net/http"
	"os"
	"server/handler"
)

func main() {
	//Sets start time
	handler.SetTime()

	port := os.Getenv("PORT")
	if port == "" {
		log.Println("$PORT has not been set. Default is 8080.")
		port = "8080"
	}

	http.HandleFunc(handler.DEFAULT, handler.DefaultHandler)
	http.HandleFunc(handler.UNI_INFO_PATH, handler.UniAndCountryHandler)
	http.HandleFunc(handler.NEIGHBOUR_UNIS_PATH, handler.NeighbourUnisHandler)
	http.HandleFunc(handler.DIAG_PATH, handler.DiagHandler)

	//http.HandleFunc(handler.DIAG_PATH, handler.UniversityHandler)

	log.Println("Server starts on port " + port)
	http.ListenAndServe(":"+port, nil)
}
