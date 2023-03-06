package handler

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

// GetCountryByName
//Returns the http response from the country API
func GetCountryByName(w http.ResponseWriter, countryName string) (*http.Response, bool) {
	countryInfo, err := http.Get(COUNTRY_URL + "v3.1/name/" + countryName + "?fullText=true")
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Println("Error when getting country by name: ", err)
		return nil, false
	}
	/*Checks the status code for the response*/
	if countryInfo.StatusCode != http.StatusOK {
		if countryInfo.StatusCode == http.StatusNotFound {
			/*The country API returns a 404 (http.StatusNotFound) if there is no match
			for a specific country name.*/
			return nil, true
		}
		http.Error(w, http.StatusText(countryInfo.StatusCode), countryInfo.StatusCode)
		log.Println("Invalid country name.")
		return nil, false
	}
	return countryInfo, true
}

// GetCountryByAlphaCode
//Returns the http response from the country API
func GetCountryByAlphaCode(w http.ResponseWriter, alphaCode string) (*http.Response, bool) {
	countryInfo, err := http.Get(COUNTRY_URL + "v3.1/alpha/" + alphaCode)
	if err != nil {
		log.Println("Error when getting country by alpha code: ", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return nil, false
	}

	/*Checks the status code for the response*/
	if countryInfo.StatusCode != http.StatusOK {
		http.Error(w, http.StatusText(countryInfo.StatusCode), countryInfo.StatusCode)
		log.Println("Invalid alpha code.")
		return nil, false
	}
	return countryInfo, true
}

// GetUniByName
//Returns the http response from the university API
func GetUniByName(w http.ResponseWriter, universityName string) (*http.Response, bool) {
	uniInfo, err := http.Get(UNI_URL + "search?name_contains=" + universityName)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Println("Error getting the university url: ", err)
		return nil, false
	}
	/*Checks the status code for the response*/
	if uniInfo.StatusCode != http.StatusOK {
		http.Error(w, http.StatusText(uniInfo.StatusCode), uniInfo.StatusCode)
		log.Println("Bad response from the uni api. Status code: ", uniInfo.StatusCode)
		return nil, false
	}
	return uniInfo, true
}

// Decode
//A general function for decoding
func Decode(w http.ResponseWriter, body io.ReadCloser, list any) bool {
	err := json.NewDecoder(body).Decode(&list)
	if err != nil {
		log.Println("Error during decoding: ", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return false
	}
	return true
}

//AddCountryToArr
//Adds a country to an array
func AddCountryToArr(w http.ResponseWriter, countryCode string, countryArr *[]Country) bool {
	var tempCountryArr []Country

	//Gets the country response. success is false if it fails
	country, success := GetCountryByAlphaCode(w, countryCode)

	if !success {
		return false
	}
	if success = Decode(w, country.Body, &tempCountryArr); success == false {
		return false
	}

	/*Adding the country in the last place of the country array*/
	*countryArr = append(*countryArr, tempCountryArr[0])

	return true
}
