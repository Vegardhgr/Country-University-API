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

const VALID_NUMBER_OF_URL_PARTS = 5

func NeighbourUnisHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		handleNeighbourCountryUnisGet(w, r)
	default:
		fmt.Println("No implementation for method " + r.Method)
		_, err := fmt.Fprintf(w, "No implementation for method "+r.Method)
		if err != nil {
			log.Println("Error using fmt.Fprint function: ", err)
		}
	}
}

//handleNeighbourCountryUnisGet
/*Gets uni and country from two api's, combines the universities with the corresponding
neighbouring country, and sends it back to the user*/
func handleNeighbourCountryUnisGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")

	//Todo: add a close tage where u need it. Is it only necessary to close the r.body? And why does it still work,
	//todo: even when i do not defer it?
	countryName, universityName, limit, err := urlHandler(r)

	if err != nil {
		if strings.Compare(err.Error(), "Invalid URL") == 0 {
			http.Error(w, "Malformed URL", http.StatusBadRequest)
			log.Println("Malformed URL in request.")
			return
		}

		if strings.Compare(strings.Split(err.Error(), ":")[0], "Limit is not an int") == 0 {
			log.Println("Limit, is not convertable to type int: ", err)
			_, err = fmt.Fprint(w, "Invalid limit given. Status code: ", http.StatusBadRequest)
			if err != nil {
				log.Println("Error using fmt.Fprint function: ", err)
			}
			return
		}
		if strings.Compare(strings.Split(err.Error(), ":")[0], "Limit must be greater than zero") == 0 {
			log.Println("Limit, must be greater than zero:\n ", err)
			_, err = fmt.Fprint(w, "Invalid limit given: Status code: ", http.StatusBadRequest)
			if err != nil {
				log.Println("Error using fmt.Fprint function: ", err)
			}
			return
		}
	}

	countryInfo, err := GetCountryByName(countryName)

	if err != nil {
		log.Println("Error when calling the api: ", err)
		_, err = fmt.Fprint(w, "Status code: ", http.StatusInternalServerError)
		if err != nil {
			log.Println("Error using fmt.Fprint function: ", err)
		}
	}

	var allUnis []UniInfo
	/*Retrieving all universities by the name given by the user as there is no good way to search
	for universities in a specific country in the api for universities*/
	uni, err := GetUniByName(universityName) //http.Get(UNI_URL + "search?name=" + universityName)

	//Todo: ask if errors should be handle where they are happening or where the error happening are called from.
	//Todo: maybe better to handle errors where they are happening so there is no need to write multiple error handling
	//Todo: messages
	if err != nil {
		http.Error(w, "Status code: ", http.StatusBadRequest)
		log.Println("Something wrong when trying to get universities from the api: ", err)
		return
	}

	err = Decode(uni.Body, &allUnis)

	if err != nil {
		log.Println("Error during decoding: ", err)
		_, err = fmt.Fprint(w, "Status code: ", http.StatusInternalServerError)
		if err != nil {
			log.Println("Error using fmt.Fprint function: ", err)
		}
		return
	}

	var borderingCountries []Borders

	/*Placing the bordering countries to the country given by the user, in the borderingCountries array*/
	err = Decode(countryInfo.Body, &borderingCountries) //json.NewDecoder(countryInfo.Body).Decode(&borderingCountries)

	if err != nil {
		_, err = fmt.Fprint(w, "Error during decoding: ", err, err.Error())
		if err != nil {
			log.Println("Error using fmt.Fprint function: ", err)
		}
		log.Println("Error during decoding: ", err)
	}

	/*A list for holding the unis in bordering countries*/
	unisInBorderingCountries := make([]UniInfo, 0)

	/*A list that will hold info about the bordering countries*/
	countryArr := make([]Country, len(borderingCountries[0].Borders))
	//var tempCountryArr []Country

	for i := range borderingCountries[0].Borders {
		/*Gets the country and adds it to the array*/
		AddBorderingCountryToArr(borderingCountries[0].Borders[i], &countryArr, i)
	}

	/*Combining university and the corresponding university*/
	combineUniversityAndCountry(&unisInBorderingCountries, &allUnis,
		&countryArr, limit)

	err = json.NewEncoder(w).Encode(unisInBorderingCountries)

	if err != nil {
		log.Println("Error when encoding: ", err)
		_, err = fmt.Fprint(w, "Status code: ", http.StatusInternalServerError)
		if err != nil {
			log.Println("Error when printing using fmt.Fprint: ", err)
		}
		return
	}
}

//urlHandler/*Handles the url specified by the user*/
//todo: pass the parameters as pointers instead, and then return just the error
func urlHandler(r *http.Request) (string, string, int, error) {
	urlParts := strings.Split(r.URL.String(), "/")

	if len(urlParts)-1 != VALID_NUMBER_OF_URL_PARTS ||
		strings.EqualFold(urlParts[1]+"/"+urlParts[2]+"/"+urlParts[3], NEIGHBOUR_UNIS_PATH) {
		return "", "", -1, errors.New("Invalid URL")
	}
	var err error
	countryName := urlParts[4]
	universityName := strings.Split(urlParts[5], "?")[0]

	println("countryName: ", countryName, "universityName: ", universityName)

	limitString := r.URL.Query().Get("limit")
	limit := -1

	if !strings.EqualFold(limitString, "") {
		if limit, err = strconv.Atoi(limitString); err != nil {
			return "", "", 0, errors.New("Limit is not an int")
		}

		if limit <= 0 {
			return "", "", 0, errors.New("Limit must be greater than zero")
		}
	}
	fmt.Println("Limit string: ", limitString, "Limit: ", limit)
	return countryName, universityName, limit, nil
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
					fmt.Println("Country: ", (*allUnis)[i].Country, " :: Languages: ", (*countryArr)[j].Languages)
					(*unisInBorderingCountries)[unisAdded].CountryInfo.Languages = (*countryArr)[j].Languages
					(*unisInBorderingCountries)[unisAdded].CountryInfo.StreetMap = (*countryArr)[j].StreetMap["openStreetMaps"]
					unisAdded++
				}
			}
			if limit != -1 && limit == unisAdded {
				break
			}
		}
	}
}
