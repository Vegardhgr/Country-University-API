package handler

import (
	"fmt"
	"net/http"
)

func NeighbourUnisHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		handleNeighbourCountryUnisGet(w, r)
	default:
		fmt.Fprintf(w, "No implementation for method "+r.Method)
		fmt.Println("No implementation for method " + r.Method)
	}
}

func handleNeighbourCountryUnisGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")

	countryName := r.URL.Path
	fmt.Println(countryName)
}
