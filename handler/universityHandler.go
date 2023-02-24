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

			countryAsByteArr, err := ioutil.ReadAll(countryRetrievedFromUrl.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				log.Println("Error when trying to read using ioutil.ReadAll", err)
			}

			//var code []countryCode
			var countries []Country
			countryToBeAddedToUni := make([]CountryInfo, 1)
			//var countryToBeAddedToUni []CountryInfo

			json.Unmarshal(countryAsByteArr, &countries)
			fmt.Println(countryToBeAddedToUni[0])
			fmt.Println(countries[0].Languages)
			countryToBeAddedToUni[0].Languages = countries[0].Languages
			fmt.Println(2)

			countryToBeAddedToUni[0].Map = countries[0].Map["openStreetMaps"]
			fmt.Println(3)

			country[unis[i].Isocode] = countryToBeAddedToUni[0]
			fmt.Println(4)

			//unis[i].CountryInfo.Map = country
			/*fmt.Println(len(countries), " length")
			fmt.Println(countries[len(countries)-1])*/
			//json.Unmarshal(countryAsByteArr, &code)

			//fmt.Println(countryNumber, ". UNI Isocode: ", unis[i].Isocode, ", ", unis[i].Country)
			//fmt.Println("Country iso: ", code[0].Cca2)
			//countryNumber++

			/*Adding the country information to the map with isocode as key*/
			//country[code[0].Cca2] = countries[len(countries)-1]
		}

		//Adding the matching country into the university struct
		unis[i].CountryInfo = country[unis[i].Isocode]
	}
	//json.NewEncoder(w).Encode(unis)
	err = json.NewEncoder(w).Encode(unis)
	//error2 := encoder
	if err != nil {
		fmt.Println("Error during encoding: ", err)
	}
}
