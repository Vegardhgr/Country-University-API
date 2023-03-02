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

		if strings.Compare(strings.Split(err.Error(), ":")[0], "Invalid limit given") == 0 {
			fmt.Println("Limit, is not convertable to type int:\n ", err)
			fmt.Fprint(w, "Limit, is not convertable to type int:\n", err)
			return
		}
	}

	countryInfo, err := getCountryByName(countryName)

	var allUnis []UniInfo
	/*Retrieving all universities by the name given by the user as there is no good way to search
	for universities in a specific country in the api for universities*/
	uni, err := getUniByName(universityName) //http.Get(UNI_URL + "search?name=" + universityName)

	//Todo: ask if errors should be handle where they are happening or where the error happening are called from.
	//Todo: maybe better to handle errors where they are happening so there is no need to write multiple error handling
	//Todo: messages
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		fmt.Println("Something wrong when trying to get universities from the api: ", err)
		return
	}

	err = json.NewDecoder(uni.Body).Decode(&allUnis)

	if err != nil {
		fmt.Println("Error during decoding: ", err)
		return
	}

	var borderingCountries []Borders

	/*Placing the bordering countries to the country given by the user, in the borderingCountries array*/
	err = json.NewDecoder(countryInfo.Body).Decode(&borderingCountries)

	if err != nil {
		fmt.Fprintf(w, "Error during decoding: ", err, err.Error())
		fmt.Println("Error during decoding: ", err)
	}

	unisInBorderingCountries := make([]UniInfo, 0)

	countryArr := make([]Country, len(borderingCountries[0].Borders))
	//var tempCountryArr []Country

	for i := range borderingCountries[0].Borders {
		/*Gets the country and adds it to the array*/
		getAndAddCountryToArr(borderingCountries[0].Borders[i], &countryArr, i)
	}

	combineUniversityAndCountry(&unisInBorderingCountries, &allUnis, &borderingCountries,
		&countryArr, limit)

	//if len(unisInBorderingCountries) > 0 {
	err = json.NewEncoder(w).Encode(unisInBorderingCountries)
	if err != nil {
		log.Println("Error when encoding: ", err)
		_, err = fmt.Fprint(w, http.StatusInternalServerError)
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
	universityName := urlParts[5]
	limitString := r.URL.Query().Get("limit")
	limit := -1

	if !strings.EqualFold(limitString, "") {
		if limit, err = strconv.Atoi(limitString); err != nil {
			return "", "", 0, errors.New("Make sure that the limit is an int")
		}

		if limit <= 0 {
			return "", "", 0, errors.New("Limit must be greater than zero")
		}
	}
	fmt.Println(limitString)
	return countryName, universityName, limit, nil
}

/*Changing the format from %20 to regular space*/
/*func addSpace(text string) string {
	return strings.ReplaceAll(text, "%20", " ")
}*/

/*Combining universities with corresponding country information*/
func combineUniversityAndCountry(unisInBorderingCountries *[]UniInfo, allUnis *[]UniInfo,
	borderingCountries *[]Borders, countryArr *[]Country, limit int) {
	allUnisDereference := *allUnis
	if len(allUnisDereference) > 0 {
		unisAdded := 0
		for i := range allUnisDereference {
			for j := range (*borderingCountries)[0].Borders {
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

func getAndAddCountryToArr(countryCode string, countryArr *[]Country, placeToAddInArr int) {
	var tempCountryArr []Country

	country, err := getCountryByAlphaCode(countryCode) //http.Get(COUNTRY_URL + "alpha/" + countryCode)

	//Todo: remember to return this error back
	if err != nil {
		fmt.Println("Error with http get method: ", err)
		return
	}

	err = json.NewDecoder(country.Body).Decode(&tempCountryArr)

	fmt.Println("Length: ", len(*countryArr))

	//fmt.Println("Country ", i, ". ", tempCountryArr[0].Name.Official, " :: ", tempCountryArr[0].Languages)

	//Todo: Ask how to not treat tempCountryArr as an array but rather as a Country object
	/*Adding the country in the last place of the country array*/
	(*countryArr)[placeToAddInArr] = tempCountryArr[0]

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
	return http.Get(UNI_URL + "search?name=" + universityName)
}
