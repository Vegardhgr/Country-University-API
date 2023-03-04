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
		http.Error(w, "No implementation for method "+r.Method, http.StatusNotImplemented)
	}
}

//handleUniAndCountryGet
/*Gets uni and country from two api's, combines them, and sends it back to the user*/
func handleUniAndCountryGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	uniName := strings.Split(r.URL.String(), "/")[4]

	fmt.Println(uniName)
	uniInfoOutput, success := GetUniByName(w, uniName)
	//http.Get(UNI_URL + "search?name_contains=" + uniName)
	if !success {
		return
	}

	var unis []UniInfo

	//Will store visited countries with isocode as key
	country := make(map[string]CountryInfo)

	if success := Decode(w, uniInfoOutput.Body, &unis); success == false {
		return
	}

	countries := make([]Country, 0)
	var countryToBeAddedToUni CountryInfo

	counter := 0
	for i := range unis {
		//Checks if the country already exists in map
		if _, ok := country[unis[i].Isocode]; ok == false {
			counter++
			fmt.Println(counter, " : ", unis[i].Country)
			//Country does not exist in map, so it must be added to it
			length := len(country)
			if success := AddBorderingCountryToArr(w, unis[i].Isocode, &countries); success == false {
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
	err := json.NewEncoder(w).Encode(unis)
	if err != nil {
		log.Println("Error during encoding: " + err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}
