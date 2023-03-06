package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

// UniAndCountryHandler
//Handles different methods in university path. Only functionality for get
func UniAndCountryHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		handleUniAndCountryGet(w, r)
	default:
		log.Println("No implementation for method " + r.Method)
		http.Error(w, "No implementation for method "+r.Method, http.StatusNotImplemented)
	}
}

//handleUniAndCountryGet
/*Gets uni and country from two APIs, combines them, and sends it back to the user*/
func handleUniAndCountryGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")

	uniName := ""
	/*Validates the url provided by the user. Success will be false if the url is invalid*/
	if success := urlHandlerForUni(w, r.URL.String(), &uniName); success == false {
		return
	}

	//Gets university response. success is true if it was successful
	uniInfoOutput, success := GetUniByName(w, uniName)

	if !success {
		return
	}

	var unis []UniInfo

	/*Stores visited countries with isocode as key. Makes it easy to search and retrieve a country.*/
	country := make(map[string]CountryInfo)

	if success := Decode(w, uniInfoOutput.Body, &unis); success == false {
		return
	}

	//Stores all the countries
	countries := make([]Country, 0)

	/*countryToBeAddedToUni is used to add the last country that is stored in the list named countries,
	to the map named country.*/
	var countryToBeAddedToUni CountryInfo

	/*Loops through all the unis*/
	for i := range unis {
		//Checks if the country already exists in map
		if _, ok := country[unis[i].Isocode]; ok == false {
			//Country does not exist in map, so it must be added to it
			length := len(country)
			if success := AddCountryToArr(w, unis[i].Isocode, &countries); success == false {
				return
			}

			/*Need to go through countryToBeAddedToUni because to add country in the map,
			the object must be of same type as the value in the map.*/
			countryToBeAddedToUni.Languages = countries[length].Languages
			countryToBeAddedToUni.StreetMap = countries[length].StreetMap["openStreetMaps"]

			/*Adding the country information in the map with isocode as key*/
			country[unis[i].Isocode] = countryToBeAddedToUni
		}
		//Adding the matching country into the university struct
		unis[i].CountryInfo = country[unis[i].Isocode]
	}

	//Sending the uni and country info back to the user
	enc := json.NewEncoder(w)
	//SetIndent formats the output like pretty print in postman
	enc.SetIndent("", "    ")
	err := enc.Encode(unis)

	if err != nil {
		log.Println("Error during encoding: " + err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

/*Validates the url provided by the user*/
func urlHandlerForUni(w http.ResponseWriter, url string, uniName *string) bool {
	urlParts := strings.Split(url, "/")
	*uniName = strings.Split(url, "/")[4]

	/*Checks that the url has the required length, and that a university name is specified*/
	if len(urlParts)-1 != VALID_NUMBER_OF_URL_PARTS_UNI_HANDLER ||
		strings.Compare(*uniName, "") == 0 {
		http.Error(w, "Expecting format .../{uni name}. Please provide a university name.",
			http.StatusNotFound)
		log.Println("Malformed URL in request for universities.")
		return false
	}

	return true
}
