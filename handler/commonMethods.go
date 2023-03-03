package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// GetCountryByName /*Returns the http response from the country api*/
func GetCountryByName(countryName string) (*http.Response, error) {
	return http.Get(COUNTRY_URL + "name/" + countryName + "?fullText=true")
}

// GetCountryByAlphaCode /*Returns the http response from the country api*/
func GetCountryByAlphaCode(alphaCode string) (*http.Response, error) {
	return http.Get(COUNTRY_URL + "alpha/" + alphaCode)
}

// GetUniByName /*Returns the http response from the university api*/
func GetUniByName(universityName string) (*http.Response, error) {
	fmt.Println("University name: ", universityName)
	return http.Get(UNI_URL + "search?name_contains=" + universityName)
}

// Decode /*A general function for decoding*/
func Decode(body io.ReadCloser, list any) error {
	return json.NewDecoder(body).Decode(&list)
}

//AddBorderingCountryToArr/*Adds a country to an array*/
func AddBorderingCountryToArr(countryCode string, countryArr *[]Country, placeToAddInArr int) {
	var tempCountryArr []Country

	country, err := GetCountryByAlphaCode(countryCode)

	//Todo: remember to return this error back
	if err != nil {
		fmt.Println("Error with http get method: ", err)
		return
	}
	err = Decode(country.Body, &tempCountryArr)

	//Todo: Return error back
	if err != nil {
		fmt.Println("Error during decoding: ", err)
	}

	//fmt.Println("Country ", i, ". ", tempCountryArr[0].Name.Official, " :: ", tempCountryArr[0].Languages)

	//Todo: Ask how to not treat tempCountryArr as an array but rather as a Country object
	/*Adding the country in the last place of the country array*/
	*countryArr = append(*countryArr, tempCountryArr[0])

	//Todo: I think i can remove this line
	tempCountryArr[0].Languages = nil
}
