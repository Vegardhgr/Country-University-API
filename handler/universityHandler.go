package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func UniAndCountryHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		handleUniAndCountryGet(w, r)
	default:
		fmt.Fprintf(w, "No implementation for method "+r.Method)
		fmt.Println("No implementation for method " + r.Method)
	}
}

func handleUniAndCountryGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	uniName := strings.Split(r.URL.String(), "/")[4]

	fmt.Println(uniName)

	uniInfoOutput, err := http.Get(UNI_URL + "search?name=" + uniName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Println("Error getting the url:", err)
		return
	}

	var unis []UniInfo
	country := make(map[string]CountryInfo)

	uniInfo, _ := ioutil.ReadAll(uniInfoOutput.Body)

	err = json.Unmarshal(uniInfo, &unis)
	if err != nil {
		fmt.Println("Error during unmarshalling ", err)
		return
	}

	for i := range unis {
		//Checks if the country already exists in map
		if _, ok := country[unis[i].Isocode]; ok == false {
			/*Gets the country information based on which country the university
			is located in*/

			//Retrieving the country that matches with the country of the university
			countryRetrievedFromUrl, err := http.Get(COUNTRY_URL + "alpha/" + unis[i].Isocode + "?fullText=true")
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				log.Println("Error trying to get the url:", err)
				return
			}

			/*Reading the content of what was sent from the api*/
			countryAsByteArr, err := ioutil.ReadAll(countryRetrievedFromUrl.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				log.Println("Error when trying to read using ioutil.ReadAll function ", err)
			}

			var countries []Country
			countryToBeAddedToUni := make([]CountryInfo, 1)

			/*Unmarshalling the content into the countries list*/
			err = json.Unmarshal(countryAsByteArr, &countries)
			if err != nil {
				fmt.Println("Error during unmarshalling: ", err)
			}

			/*Adding the languages to the struct that will be used
			in the university struct*/
			countryToBeAddedToUni[0].Languages = countries[0].Languages

			/*Adding the open street maps link to the struct that will be used
			in the university struct*/
			countryToBeAddedToUni[0].Map = countries[0].Map["openStreetMaps"]

			/*Adding the country information in the map with isocode as key*/
			country[unis[i].Isocode] = countryToBeAddedToUni[0]
		}

		//Adding the matching country into the university struct
		unis[i].CountryInfo = country[unis[i].Isocode]
	}

	//Sending the universities back to the user
	err = json.NewEncoder(w).Encode(unis)
	if err != nil {
		fmt.Println("Error during encoding: ", err)
	}
}
