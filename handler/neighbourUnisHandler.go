package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

const MIN_VALID_URL_PARTS = 5
const MAX_VALID_URL_PARTS = 6

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

	countryName, universityName, limit, err := urlHandler(r)

	if err != nil {
		if strings.Compare(err.Error(), "Invalid URL") == 0 {
			http.Error(w, "Malformed URL", http.StatusBadRequest)
			log.Println("Malformed URL in request.")
			return
		}

		fmt.Println("Error: ", err.Error())
		if strings.Compare(strings.Split(err.Error(), ":")[0], "Invalid limit given") == 0 {
			fmt.Println("Limit, is not convertable to type int:\n ", err)
			fmt.Fprint(w, "Limit, is not convertable to type int:\n", err)
			return
		}
	}
	fmt.Println(w, countryName, universityName, limit)

	countryInfo, err := http.Get(COUNTRY_URL + "name/" + countryName + "?fullText=true")

	var allUnis []UniInfo
	/*Retrieving all universities by the name given by the user as there is no good way to search
	for universities in a specific country in the api for universities*/
	uni, err := http.Get(UNI_URL + "search?name=" + universityName)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		fmt.Println("Something wrong when trying to get universities from the api: ", err)
		return
	}

	/*uniBody, err := ioutil.ReadAll(uni.Body)

	if err != nil {
		fmt.Println("Error when reading the get request from the university api: ", err)
		return
	}

	err = json.Unmarshal(uniBody, &allUnis)

	if err != nil {
		fmt.Println("Error when unmarshalling: ", err)
		return
	}*/

	err = json.NewDecoder(uni.Body).Decode(&allUnis)

	if err != nil {
		fmt.Println("Error during decoding: ", err)
		return
	}

	type borders struct {
		Borders []string `json:"borders"`
	}

	var borderingCountries []borders
	//borderingCountries := make(borders, 10)

	/*Placing the bordering countries to the country given by the user, in the borderingCountries array*/
	err = json.NewDecoder(countryInfo.Body).Decode(&borderingCountries)

	if err != nil {
		fmt.Fprintf(w, "Error during decoding: ", err)
		fmt.Println("Error during decoding: ", err)
	}
	//OR use unmarshal (see following code)
	/*countryInfoBody, err := ioutil.ReadAll(countryInfo.Body)
	if err != nil {
		fmt.Println("Error reading country info: ", err)
	}
	err = nil
	err = json.Unmarshal(countryInfoBody, &borderingCountries)
	if err != nil {
		fmt.Println("Error during unmarshalling: ", err)
	}
	json.Unmarshal(countryInfo.Body, &borderingCountries)*/

	//var allUnis []UniInfo
	oldUnisRange := 0
	var unisInBorderingCountries []UniInfo

	/*Looping through the bordering countries of the country given by the user*/
	/*for i := range borderingCountries[0].Borders {
	country, err := http.Get(COUNTRY_URL + "alpha/" + borderingCountries[0].Borders[i])
	if err != nil {
		fmt.Println("Error with http get method: ", err)
		return
	}

	var countryArr []Country
	err = json.NewDecoder(country.Body).Decode(&countryArr)

	if err != nil {
		fmt.Println("Error during decoding: ", err)
	}

	fmt.Println("Official: " + countryArr[0].Name.Official)      */

	/*Not a good idea to retrieve universities by country as country name in the api for universities can differ
	from the name in the country api*/
	/*uni, err := http.Get(UNI_URL + "search?name=" + universityName + "&country=" + countryArr[0].Name.Official)
	if err != nil {
		fmt.Fprintf(w, "Could not find any countries or universities with the parameters given")
		return
	}
	err = json.NewDecoder(uni.Body).Decode(&allUnis)

	if err != nil {
		fmt.Println("Error during decoding: ", err)
	}*/

	fmt.Println("Length: ", len(allUnis))
	fmt.Println("Old range: ", oldUnisRange)

	//unisInBorderingCountries := make([]UniInfo, 1000)
	//var countryArr []Country
	countryArr := make([]Country, len(borderingCountries[0].Borders))
	var tempCountryArr []Country

	for i := range borderingCountries[0].Borders {
		country, err := http.Get(COUNTRY_URL + "alpha/" + borderingCountries[0].Borders[i])
		if err != nil {
			fmt.Println("Error with http get method: ", err)
			return
		}

		err = json.NewDecoder(country.Body).Decode(&tempCountryArr)
		fmt.Println("Country ", i, ". ", tempCountryArr[0].Name.Official, " :: ", tempCountryArr[0].Languages)
		countryArr[i] = tempCountryArr[0]

		if err != nil {
			fmt.Println("Error during decoding: ", err)
		}
		tempCountryArr[0].Languages = nil
	}

	if len(allUnis) > 0 {
		unisAdded := 0
		for i := range allUnis {
			for j := range borderingCountries[0].Borders {
				//fmt.Println(i, ". Uni name: ", allUnis[i].UniName, " :: Uni iso: ", allUnis[i].Isocode, " :: Country iso: ", countryArr[j].Cca2)
				if strings.Compare(allUnis[i].Isocode, countryArr[j].Cca2) == 0 {
					unisInBorderingCountries = append(unisInBorderingCountries, allUnis[i])
					fmt.Println("Country: ", allUnis[i].Country, " :: Languages: ", countryArr[j].Languages)
					unisInBorderingCountries[unisAdded].CountryInfo.Languages = countryArr[j].Languages
					unisInBorderingCountries[unisAdded].CountryInfo.Map = countryArr[j].Map["openStreetMaps"]
					unisAdded++
				}
			}
			if limit != -1 && limit == unisAdded {
				break
			}
		}
	}
	if len(unisInBorderingCountries) > 0 {
		json.NewEncoder(w).Encode(unisInBorderingCountries)
	} else {
		fmt.Fprint(w, "No universities with university name '", addSpace(universityName),
			"', in the bordering countries to '", addSpace(countryName), "'.")
	}
}

func urlHandler(r *http.Request) (string, string, int, error) {
	urlParts := strings.Split(r.URL.String(), "/")

	if len(urlParts)-1 < MIN_VALID_URL_PARTS || len(urlParts)-1 > MAX_VALID_URL_PARTS ||
		strings.EqualFold(urlParts[1]+"/"+urlParts[2]+"/"+urlParts[3], NEIGHBOUR_UNIS_PATH) {
		return "", "", -1, errors.New("Invalid URL")
	}

	var err error
	countryName := urlParts[4]
	universityName := urlParts[5]
	limit := -1

	if len(urlParts)-1 == MAX_VALID_URL_PARTS {
		limit, err = strconv.Atoi(urlParts[6])
		if err != nil {
			return "", "", -1, errors.New("Invalid limit given: " + err.Error())
		}
	}

	return countryName, universityName, limit, nil
}

/*Changing the format from %20 to regular space*/
func addSpace(text string) string {
	return strings.ReplaceAll(text, "%20", " ")
}
