package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func NeighbourUnisHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		handleNeighbourCountryUnisGet(w, r)
	default:
		fmt.Println("No implementation for method " + r.Method)
		http.Error(w, "No implementation for method "+r.Method, http.StatusNotImplemented)
	}
}

//handleNeighbourCountryUnisGet
/*Gets uni and country from two api's, combines the universities with the corresponding
neighbouring country, and sends it back to the user*/
func handleNeighbourCountryUnisGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")

	countryName := ""
	universityName := ""
	limit := -1
	success := urlHandler(w, r, &countryName, &universityName, &limit)

	if !success {
		return
	}

	countryInfo, success := GetCountryByName(w, countryName)

	if !success {
		return
	}

	var allUnis []UniInfo
	/*Retrieving all universities by the name given by the user as there is no good way to search
	for universities in a specific country in the api for universities*/
	uni, success := GetUniByName(w, universityName) //http.Get(UNI_URL + "search?name=" + universityName)

	if !success {
		return
	}

	if success = Decode(w, uni.Body, &allUnis); success == false {
		return
	}

	var borderingCountries []Borders

	/*Placing the bordering countries to the country given by the user, in the borderingCountries array*/
	if success = Decode(w, countryInfo.Body, &borderingCountries); success == false {
		return
	}

	/*A list for holding the unis in bordering countries*/
	unisInBorderingCountries := make([]UniInfo, 0)

	/*A list that will hold info about the bordering countries*/
	countryArr := make([]Country, len(borderingCountries[0].Borders))
	//var tempCountryArr []Country
	for i := range borderingCountries[0].Borders {
		/*Gets the country and adds it to the array*/
		if success := AddCountryToArr(w, borderingCountries[0].Borders[i], &countryArr); success == false {
			return
		}
	}
	/*Combining university and the corresponding country*/
	combineUniversityAndCountry(&unisInBorderingCountries, &allUnis,
		&countryArr, limit)

	enc := json.NewEncoder(w)
	//SetIndent formats the output like pretty print in postman
	enc.SetIndent("", "    ")
	err := enc.Encode(unisInBorderingCountries)

	if err != nil {
		log.Println("Error when encoding: ", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

//urlHandler
///*Handles the url specified by the user*/
func urlHandler(w http.ResponseWriter, r *http.Request, countryName,
	universityName *string, limit *int) bool {
	urlParts := strings.Split(r.URL.String(), "/")

	/*Validates the url*/
	if len(urlParts)-1 != VALID_NUMBER_OF_URL_PARTS_NEIGHBOUR_UNIS_HANDLER {
		http.Error(w, "Expecting format .../{country name}/{uni name}",
			http.StatusNotFound)
		log.Println("Malformed URL in request for neighbouring universities.")
		return false
	}

	*countryName = urlParts[4]
	*universityName = strings.Split(urlParts[5], "?")[0]

	limitString := r.URL.Query().Get("limit")

	if !strings.EqualFold(limitString, "") {
		var err error
		if *limit, err = strconv.Atoi(limitString); err != nil {
			log.Println("Limit, is not convertable to type int: ", err)
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return false
		}

		if *limit <= 0 {
			log.Println("Limit, must be greater than zero:\n ", err)
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return false
		}
	}
	return true
}

//combineUniversityAndCountry/*Combining universities with corresponding country information*/
func combineUniversityAndCountry(unisInBorderingCountries *[]UniInfo, allUnis *[]UniInfo,
	countryArr *[]Country, limit int) {
	if len(*allUnis) > 0 {
		unisAdded := 0
		for i := range *allUnis {
			for j := range *countryArr {
				//fmt.Println(i, ". Uni name: ", allUnis[i].UniName, " :: Uni iso: ", allUnis[i].Isocode, " :: Country iso: ", countryArr[j].Cca2)
				if strings.Compare((*allUnis)[i].Isocode, (*countryArr)[j].Cca2) == 0 {
					*unisInBorderingCountries = append(*unisInBorderingCountries, (*allUnis)[i])
					(*unisInBorderingCountries)[unisAdded].CountryInfo.Languages = (*countryArr)[j].Languages
					(*unisInBorderingCountries)[unisAdded].CountryInfo.StreetMap = (*countryArr)[j].StreetMap["openStreetMaps"]
					unisAdded++
				}
			}
			if limit == unisAdded {
				break
			}
		}
	}
}
