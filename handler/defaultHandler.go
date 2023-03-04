package handler

import (
	"fmt"
	"net/http"
)

func DefaultHandler(w http.ResponseWriter, r *http.Request) {
	// Ensure interpretation as HTML by client (browser)
	w.Header().Set("content-type", "text/html")

	// Offer information for redirection to paths
	output := "<h1>Welcome!<h1><br><h3>This service does not provide any functionality on root path level." +
		" Please try one of the paths below<h3><br><br>" +
		"<h6>Search by university:<h6><br>" +
		"<a href=\"" + UNI_INFO_PATH + "\">" + UNI_INFO_PATH + "</a><br><br>" +
		"<h6>Search for universities in a neighbouring country:<h6><br>" +
		"<a href=\"" + NEIGHBOUR_UNIS_PATH + "\">" + NEIGHBOUR_UNIS_PATH + "</a><br><br>" +
		"<h6>For diagnostics:<h6><br>" +
		"<a href=\"" + DIAG_PATH + "\">" + DIAG_PATH + "</a>"

	// Write output to client
	_, err := fmt.Fprintf(w, "%v", output)

	// Deal with error if any
	if err != nil {
		http.Error(w, "Error when returning output", http.StatusInternalServerError)
	}
}
