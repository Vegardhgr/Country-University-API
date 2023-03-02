package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
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

func handleNeighbourCountryUnisGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")

	//Todo: add a close tage where u need it. Is it only necessary to close the r.body? And why does it still work,
	//todo: even when i do not defer it?
	err := r.Body.Close()
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

	countryInfo, err := getCountryByName(countryName)

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
	uni, err := getUniByName(universityName) //http.Get(UNI_URL + "search?name=" + universityName)

	//Todo: ask if errors should be handle where they are happening or where the error happening are called from.
	//Todo: maybe better to handle errors where they are happening so there is no need to write multiple error handling
	//Todo: messages
	if err != nil {
		http.Error(w, "Status code: ", http.StatusBadRequest)
		log.Println("Something wrong when trying to get universities from the api: ", err)
		return
	}

	err = decode(uni.Body, &allUnis)

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
	err = decode(countryInfo.Body, &borderingCountries) //json.NewDecoder(countryInfo.Body).Decode(&borderingCountries)

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
		addBorderingCountryToArr(borderingCountries[0].Borders[i], &countryArr, i)
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

/*Handles the url specified by the user*/
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

/*Combining universities with corresponding country information*/
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
					(*unisInBorderingCountries)[unisAdded].CountryInfo.Map = (*countryArr)[j].Map["openStreetMaps"]
					unisAdded++
				}
			}
			if limit != -1 && limit == unisAdded {
				break
			}
		}
	}
}

/*Adds a country to an array*/
func addBorderingCountryToArr(countryCode string, countryArr *[]Country, placeToAddInArr int) {
	var tempCountryArr []Country

	country, err := getCountryByAlphaCode(countryCode)

	//Todo: remember to return this error back
	if err != nil {
		fmt.Println("Error with http get method: ", err)
		return
	}
	err = decode(country.Body, &tempCountryArr)
	//err = json.NewDecoder(country.Body).Decode(&tempCountryArr)

	fmt.Println("Length: ", len(*countryArr))

	//fmt.Println("Country ", i, ". ", tempCountryArr[0].Name.Official, " :: ", tempCountryArr[0].Languages)

	//Todo: Ask how to not treat tempCountryArr as an array but rather as a Country object
	/*Adding the country in the last place of the country array*/
	(*countryArr)[placeToAddInArr] = tempCountryArr[0]

	//Todo: Return error back
	if err != nil {
		fmt.Println("Error during decoding: ", err)
	}

	//Todo: I think i can remove this line
	tempCountryArr[0].Languages = nil
}

/*Returns the http response from the country api*/
func getCountryByName(countryName string) (*http.Response, error) {
	return http.Get(COUNTRY_URL + "name/" + countryName + "?fullText=true")
}

/*Returns the http response from the country api*/
func getCountryByAlphaCode(alphaCode string) (*http.Response, error) {
	return http.Get(COUNTRY_URL + "alpha/" + alphaCode)
}

/*Returns the http response from the university api*/
func getUniByName(universityName string) (*http.Response, error) {
	fmt.Println("University name: ", universityName)
	return http.Get(UNI_URL + "search?name_contains=" + universityName)
}

/*A general function for decoding*/
func decode(body io.ReadCloser, list any) error {
	return json.NewDecoder(body).Decode(&list)
}
