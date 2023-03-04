package handler

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

// GetCountryByName
///*Returns the http response from the country api*/
func GetCountryByName(w http.ResponseWriter, countryName string) (*http.Response, bool) {
	countryInfo, err := http.Get(COUNTRY_URL + "name/" + countryName + "?fullText=true")
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Println("Error when getting country")
		return nil, false
	}
	return countryInfo, true
}

// GetCountryByAlphaCode
///*Returns the http response from the country api*/
func GetCountryByAlphaCode(w http.ResponseWriter, alphaCode string) (*http.Response, bool) {
	countryInfo, err := http.Get(COUNTRY_URL + "alpha/" + alphaCode)
	if err != nil {
		log.Println("Error with http get method: ", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return nil, false
	}
	return countryInfo, true
}

// GetUniByName
///*Returns the http response from the university api*/
func GetUniByName(w http.ResponseWriter, universityName string) (*http.Response, bool) {
	uniInfo, err := http.Get(UNI_URL + "search?name" + universityName)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Println("Error getting the url: ", err)
		return nil, false
	}
	return uniInfo, true
}

// Decode
///*A general function for decoding*/
func Decode(w http.ResponseWriter, body io.ReadCloser, list any) bool {
	err := json.NewDecoder(body).Decode(&list)
	if err != nil {
		log.Println("Error during decoding: ", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return false
	}
	return true
}

//AddBorderingCountryToArr
///*Adds a country to an array*/
func AddBorderingCountryToArr(w http.ResponseWriter, countryCode string, countryArr *[]Country) bool {
	var tempCountryArr []Country

	country, success := GetCountryByAlphaCode(w, countryCode)

	if !success {
		return false
	}
	if success = Decode(w, country.Body, &tempCountryArr); success == false {
		return false
	}

	//fmt.Println("Country ", i, ". ", tempCountryArr[0].Name.Official, " :: ", tempCountryArr[0].Languages)

	//Todo: Ask how to not treat tempCountryArr as an array but rather as a Country object
	/*Adding the country in the last place of the country array*/
	*countryArr = append(*countryArr, tempCountryArr[0])

	//Todo: I think i can remove this line
	tempCountryArr[0].Languages = nil
	return true
}
