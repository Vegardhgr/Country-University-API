package handler

import (
	"fmt"
	"log"
	"net/http"
)

func DefaultHandler(w http.ResponseWriter, r *http.Request) {
	// Ensure interpretation as HTML by client (browser)
	w.Header().Set("content-type", "text/html")

	// Offer information for redirection to paths
	output := "<h1>Welcome!<h1><h3>This service does not provide any functionality on root path level." +
		" Please try one of the paths below<h3>" +
		"<h5 style=\"background-color: lightblue; width: 250px;\">Search by university:<br>" +
		"<a href=\"" + UNI_INFO_PATH + "\">" + UNI_INFO_PATH + "</a></h5>" +
		"<h5 style=\"background-color: lightblue; width: 250px;\">Search for universities in a neighbouring country:<br>" +
		"<a href=\"" + NEIGHBOUR_UNIS_PATH + "\">" + NEIGHBOUR_UNIS_PATH + "</a></h5>" +
		"<h5 style=\"background-color: lightblue; width: 250px;\">For diagnostics:<br>" +
		"<a href=\"" + DIAG_PATH + "\">" + DIAG_PATH + "</a></h5>"

	// Write output to client
	_, err := fmt.Fprintf(w, "%v", output)

	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Println("Error when returning output.")
	}
}
