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
		"<h6>Search by university:" +
		"<a href=\"" + UNI_INFO_PATH + "\">" + UNI_INFO_PATH + "</a></h6><br>" +
		"<h6>Search for universities in a neighbouring country:<br>" +
		"<a href=\"" + NEIGHBOUR_UNIS_PATH + "\">" + NEIGHBOUR_UNIS_PATH + "</a></h6><br>" +
		"<h6>For diagnostics:" +
		"<a href=\"" + DIAG_PATH + "\">" + DIAG_PATH + "</a></h6>"

	// Write output to client
	_, err := fmt.Fprintf(w, "%v", output)

	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Println("Error when returning output.")
	}
}
