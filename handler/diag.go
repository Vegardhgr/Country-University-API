package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// DiagHandler
///*Handles different methods in diag path. Only functionality for get*/
func DiagHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		diagHandlerGet(w, r)
	default:
		fmt.Println("No implementation for method " + r.Method)
		http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
	}
}

//diagHandlerGet
//*Contains the functionality for diag path*/
func diagHandlerGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	/*urlHandlerDiag returns false if the url does not meet the required specifications*/
	if !urlHandlerDiag(w, r.URL.String()) {
		return
	}
	var statusCodeCountry int
	country, err := http.Get(COUNTRY_URL)
	if err != nil {
		statusCodeCountry = http.StatusServiceUnavailable
		log.Println("Error when getting the following url: ", COUNTRY_URL, ". Error: ", err)
	} else { /*Contains a body*/
		statusCodeCountry = country.StatusCode
	}

	var statusCodeUni int
	uni, err := http.Get(UNI_URL)
	if err != nil {
		statusCodeCountry = http.StatusServiceUnavailable
		log.Println("Error when getting the following url: ", UNI_URL, ". Error: ", err)
	} else { /*Contains a body*/
		statusCodeUni = uni.StatusCode
	}

	var diag Diag
	diag.CountriesApiStatus = strconv.Itoa(statusCodeCountry)
	diag.UnisApiStatus = strconv.Itoa(statusCodeUni)
	diag.Version = "V1"
	diag.Uptime = GetTime()
	enc := json.NewEncoder(w)
	//Set indent pretty prints as in postman
	enc.SetIndent("", "    ")
	err = enc.Encode(diag)

	if err != nil {
		log.Println("Error during encoding: ", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

//urlHandlerDiag
//*Validates the url for diag. Returns true if valid.*/
func urlHandlerDiag(w http.ResponseWriter, url string) bool {
	urlParts := strings.Split(url, "/")

	/*Validates the length of the url*/
	if len(urlParts)-1 != VALID_NUMBER_OF_URL_PARTS_DIAG || urlParts[len(urlParts)-1] != "" {
		http.Error(w, "Expecting format .../diag.", http.StatusNotFound)
		log.Println("Malformed URL in request for diag.")
		return false
	}
	return true
}
