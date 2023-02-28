package handler

import (
	"fmt"
	"net/http"
)

func DefaultHandler(w http.ResponseWriter, r *http.Request) {
	// Ensure interpretation as HTML by client (browser)
	w.Header().Set("content-type", "text/html")

	// Offer information for redirection to paths
	output := "This service does not provide any functionality on root path level. Please use the paths below\n" +
		"<a href=\"" + UNI_INFO_PATH + "\">" +
		"<a href=\"" + NEIGHBOUR_UNIS_PATH + "\">" +
		"<a href=\"" + DIAG_PATH + "\">"

	// Write output to client
	_, err := fmt.Fprintf(w, "%v", output)

	// Deal with error if any
	if err != nil {
		http.Error(w, "Error when returning output", http.StatusInternalServerError)
	}
}
