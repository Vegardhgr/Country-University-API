package handler

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// NeighbourUnisHandler
//Handles different methods in neighbour path. Only functionality for get
func NeighbourUnisHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		handleNeighbourCountryUnisGet(w, r)
	default:
		log.Println("No implementation for method " + r.Method)
		http.Error(w, "No implementation for method "+r.Method, http.StatusNotImplemented)
	}
}

//handleNeighbourCountryUnisGet
/*Gets uni and country from two API, combines the universities with the corresponding
neighbouring country, and sends it back to the user*/
func handleNeighbourCountryUnisGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")

	countryName := ""
	universityName := ""
	limit := -1

	if !urlHandler(w, r, &countryName, &universityName, &limit) {
		return
	}

	//Gets country info. success is true if it was successful
	countryInfo, success := GetCountryByName(w, countryName)

	if !success {
		return
	}

	/*A list for holding the unis in bordering countries*/
	unisInBorderingCountries := make([]UniInfo, 0)

	/*countryInfo can be nil while success is true when the country API returns a 404, as this means the
	country does not exist. For consistency, it is better to return an empty json object instead of a 404,
	as this is how the university API does it.*/
	if countryInfo != nil {
		var allUnis []UniInfo
		/*Retrieving all universities by the name given by the user as there is no good way to search
		for universities in a specific country in the API for universities*/
		uni, success := GetUniByName(w, universityName)

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

		/*A list that will hold info about the bordering countries*/
		countryArr := make([]Country, len(borderingCountries[0].Borders))
		for i := range borderingCountries[0].Borders {
			/*Gets the country and adds it to the array*/
			if success := AddCountryToArr(w, borderingCountries[0].Borders[i], &countryArr); success == false {
				return
			}
		}
		/*Combining university and the corresponding country*/
		combineUniversityAndCountry(&unisInBorderingCountries, &allUnis,
			&countryArr, limit)
	}

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
//Handles the url specified by the user
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

	//Gets the limit from the url if specified
	limitString := r.URL.Query().Get("limit")

	//Checks if a limit has been specified. If yes it must be validated
	if !strings.EqualFold(limitString, "") {
		var err error
		//Converts limit to int
		if *limit, err = strconv.Atoi(limitString); err != nil {
			log.Println("Limit, is not convertable to type int: ", err)
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return false
		}
		//Checks that limit is greater than zero
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
		/*If a limit is specified, the output will be randomized*/
		if limit != -1 {
			/*Shuffles the list so that the output will not be identical every time*/
			rand.Seed(time.Now().UnixNano())
			rand.Shuffle(len(*allUnis), func(i, j int) {
				(*allUnis)[i], (*allUnis)[j] = (*allUnis)[j], (*allUnis)[i]
			})
		}
		//Looping through universities
		for i := range *allUnis {
			//Lopping through countries
			for j := range *countryArr {
				//Comparing cca2 code. If equal, append it to unisInBorderingCountries
				if strings.Compare((*allUnis)[i].Isocode, (*countryArr)[j].Cca2) == 0 {
					*unisInBorderingCountries = append(*unisInBorderingCountries, (*allUnis)[i])
					(*unisInBorderingCountries)[unisAdded].CountryInfo.Languages = (*countryArr)[j].Languages
					(*unisInBorderingCountries)[unisAdded].CountryInfo.StreetMap = (*countryArr)[j].StreetMap["openStreetMaps"]
					unisAdded++
				}
			}
			//Breaks if limit is reached
			if limit == unisAdded {
				break
			}
		}
	}
}
