package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func DiagHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		diagHandlerGet(w, r)
	default:
		fmt.Println("No implementation for method " + r.Method)
		http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
	}
}

func diagHandlerGet(w http.ResponseWriter, r *http.Request) {
	country, err := http.Get(COUNTRY_URL)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Println("Error when getting the following url: ", COUNTRY_URL, ". Error: ", err)
		return
	}
	statusCodeCountry := country.StatusCode

	uni, err := http.Get(UNI_URL)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Println("Error when getting the following url: ", UNI_URL, ". Error: ", err)
		return
	}
	statusCodeUni := uni.StatusCode

	var diag Diag
	diag.CountriesApiStatus = strconv.Itoa(statusCodeCountry)
	diag.UnisApiStatus = strconv.Itoa(statusCodeUni)
	diag.Version = "V1"
	diag.Uptime = GetTime()
	err = json.NewEncoder(w).Encode(diag)

	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}

}
