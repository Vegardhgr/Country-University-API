package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

//func GetUniInfo()
func UniAndCountryHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	/*case http.MethodPost:
	handleUniAndCountryPost(w, r)*/
	case http.MethodGet:
		handleUniAndCountryGet(w, r)
	}
}

func handleUniAndCountryGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")

	//uniName := strings.ReplaceAll(strings.Split(r.URL.String(), "/")[4], "%20", " ")
	uniName := strings.Split(r.URL.String(), "/")[4]
	log.Println(uniName)
	uniInfoOutput, err := http.Get(UNI_URL + "search?name=" + uniName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Println("Error getting the url:", err)
		return
	}

	//countryInfoOutput, err := http.Get(COUNTRY_URL)

	/*if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Println("Error getting the url:", err)
		return
	}*/

	var unis []UniInfo
	var countries []CountryInfo

	country := make(map[string]CountryInfo)

	//var country map[string]CountryInfo

	uniInfo, _ := ioutil.ReadAll(uniInfoOutput.Body)

	json.Unmarshal(uniInfo, &unis)

	fmt.Println(len(unis))

	/*countryInfo, _ := ioutil.ReadAll(countryInfoOutput.Body)

	json.Unmarshal(countryInfo, &country)*/
	//sum := 0
	for i := range unis {
		//var cNumber int = 0
		/*Gets the country information based on which country the university
		is located in*/
		//_, ok := country[unis[i].Country]
		//log.Println(S, " : ", ok)
		if _, ok := country[unis[i].Country]; ok == false {
			/*sum++
			fmt.Println(sum)*/
			countryRetrievedFromUrl, err := http.Get(COUNTRY_URL + "name/" + unis[i].Country)
			//fmt.Println(unis[i].Country)
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

			json.Unmarshal(countryAsByteArr, &countries)
			fmt.Println("Len: ", len(countries))
			fmt.Println(countries[len(countries)-1].Name.Common)
			country[countries[len(countries)-1].Name.Common] = countries[len(countries)-1]

			//_, ok = country[unis[i].Country]
		}

		//Adding the matching country into the university struct
		unis[i].CountryInfo = country[unis[i].Country]
	}
	//json.NewEncoder(w).Encode(unis)
	err = json.NewEncoder(w).Encode(unis)
	//error2 := encoder
	if err != nil {

	}
	/*for i := range unis {
		var cNumber int = 0
		for c := range country {
			if strings.Compare(country[c].Name.Common, unis[i].Country) == 0 {
				cNumber = c
				fmt.Println(country[c])
				break
			}
		}
		unis[i].CountryInfo = country[cNumber]
	}*/

	/*for i := range unis {
		if strings.Contains(unis[i].UniName, uniName) {
			json.NewEncoder(w).Encode(unis[i])
		}
	}*/
}

/*func UniversityHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")

	uniInfoOutput, err := http.Get(UNI_URL)
	if err != nil {
		log.Fatal("Error getting the url:", err)
		return
	}

	countryInfoOutput, err := http.Get(COUNTRY_URL)

	if err != nil {
		log.Fatal("Error getting the url:", err)
	}

	var unis []UniInfo
	var country []CountryInfo

	uniInfo, _ := ioutil.ReadAll(uniInfoOutput.Body)

	json.Unmarshal(uniInfo, &unis)

	countryInfo, _ := ioutil.ReadAll(countryInfoOutput.Body)

	json.Unmarshal(countryInfo, &country)

	for i := range unis {
		var cNumber int = 0
		for c := range country {
			if strings.Compare(country[c].Name.Common, unis[i].Country) == 0 {
				cNumber = c
				fmt.Println(country[c])
				break
			}
		}
		unis[i].CountryInfo = country[cNumber]

	}

	uniName := strings.ReplaceAll(strings.Split(r.URL.String(), "/")[4], "%20", " ")

	for i := range unis {
		if strings.Contains(unis[i].UniName, uniName) {
			json.NewEncoder(w).Encode(unis[i])
		}
	}
}

func UniversityHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")

	uniInfoOutput, err := http.Get(UNI_URL)
	if err != nil {
		log.Fatal("Error getting the url:", err)
		return
	}

	countryInfoOutput, err := http.Get(COUNTRY_URL)

	if err != nil {
		log.Fatal("Error getting the url:", err)
	}

	var unis []UniInfo
	var country []CountryInfo

	uniInfo, _ := ioutil.ReadAll(uniInfoOutput.Body)

	json.Unmarshal(uniInfo, &unis)

	countryInfo, _ := ioutil.ReadAll(countryInfoOutput.Body)

	json.Unmarshal(countryInfo, &country)

	for i := range unis {
		var cNumber int = 0
		for c := range country {
			if strings.Compare(country[c].Name.Common, unis[i].Country) == 0 {
				cNumber = c
				fmt.Println(country[c])
				break
			}
		}
		unis[i].CountryInfo = country[cNumber]

	}

	uniName := strings.ReplaceAll(strings.Split(r.URL.String(), "/")[4], "%20", " ")

	for i := range unis {
		if strings.Contains(unis[i].UniName, uniName) {
			json.NewEncoder(w).Encode(unis[i])
		}
	}
}*/
