package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

func UniAndCountryHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		handleUniAndCountryGet(w, r)
	default:
		log.Println("No implementation for method " + r.Method)
		_, err := fmt.Fprint(w, "No implementation for method "+r.Method)
		if err != nil {
			log.Println("Error using fmt.Fprint function: ", err)
		}
	}
}

//handleUniAndCountryGet
/*Gets uni and country from two api's, combines them, and sends it back to the user*/
func handleUniAndCountryGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	uniName := strings.Split(r.URL.String(), "/")[4]

	fmt.Println(uniName)
	uniInfoOutput, err := GetUniByName(uniName)
	//http.Get(UNI_URL + "search?name_contains=" + uniName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Println("Error getting the url:", err)
		return
	}

	var unis []UniInfo

	//Will store visited countries with isocode as key
	country := make(map[string]CountryInfo)

	err = Decode(uniInfoOutput.Body, &unis)

	if err != nil {
		log.Println("Error during decoding: ", err)
		fmt.Fprint(w, "Error status code: ", http.StatusInternalServerError)
		return
	}

	countries := make([]Country, 0)
	var countryToBeAddedToUni CountryInfo

	for i := range unis {
		//Checks if the country already exists in map
		if _, ok := country[unis[i].Isocode]; ok == false {
			//Country does not exist in map, so it must be added to it
			length := len(country)
			AddBorderingCountryToArr(unis[i].Isocode, &countries, length)

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
	err = json.NewEncoder(w).Encode(unis)
	if err != nil {
		fmt.Println("Error during encoding: ", err)
	}
}
