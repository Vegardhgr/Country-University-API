package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func UniversityHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")

	//uniInfoOutput, err := ioutil.ReadFile("handler/UniversityInfo.json")

	uniInfoOutput, err := http.Get(UNI_URL)
	if err != nil {
		log.Fatal("Error getting the url:", err)
		return
	}
	//countryInfoOutput, err := ioutil.ReadFile("handler/CountryInfo.json")

	countryInfoOutput, err := http.Get(COUNTRY_URL)

	if err != nil {
		log.Fatal("Error getting the url:", err)
	}

	var unis []UniInfo
	var country []CountryInfo

	uniInfo, _ := ioutil.ReadAll(uniInfoOutput.Body)

	json.Unmarshal(uniInfo, &unis)

	/*for i := range unis {

	}*/

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

		//json.NewEncoder(w).Encode(unis[i])
	}

	/*for i := range uniCountry {
		json.NewEncoder(w).Encode(uniCountry[i])
	}*/

	uniName := strings.ReplaceAll(strings.Split(r.URL.String(), "/")[4], "%20", " ")

	//retrieveUniInfoUrl, err := http.Get(UNI_URL + "search?name=" + uniName)

	var uunis []UniInfo
	//uunisIndex := 0
	for i := range unis {
		if strings.Contains(unis[i].UniName, uniName) {
			json.NewEncoder(w).Encode(unis[i])
		}

	}
	//infoFromSiteAsArray, _ := ioutil.ReadAll(retrieveUniInfoUrl.Body)

	//json.Unmarshal(infoFromSiteAsArray, &uunis)

	for i := range uunis {
		json.NewEncoder(w).Encode(uunis[i])
	}

	//json.NewEncoder(w).Encode(infoFromSiteUnmarshaled)

	/*if name != "" {
		for i := range unis {
			fmt.Println(unis[i].Name, " : ", name)
			if strings.Contains(unis[i].Name, name) {
				json.NewEncoder(w).Encode(unis[i])
			}
		}
	}
	urlSplit := strings.Split(r.URL.Path, "/")

	if len(urlSplit) < 3 || urlSplit[3] != UNI_INFO {
		http.Error(w, "URL not written correctly", http.StatusBadRequest)
		log.Println("Malformed URL in request.")
		return
	}

	/*for i := range unis {
		if unis[i].Name == "Norwegian University of Science and Technology" {
			encode := json.NewEncoder(w)
			error1 := encode.Encode(unis[i])
			if error1 != nil {
				errMessage, _ := "Something went wrong during encoding: ", error
				log.Println(errMessage)
			}
			//fmt.Println(unis[i])
		}
	}*/

	/*for i := range mappy {

		fmt.Println(i)
	}*/

	/*universityInfo := UniInfo{
		Name: "Norwegian University of Science and Technology",
		Coutry: "Norway",
		Isocode: "NO",
		Webpages: ["http://www.ntnu.no/"],
		Languages: Languages {
			NNO:"Norwegian Nynorsk",
			NOB:"Norwegian Nynorsk",
			SMI:"Sami"}
			Map: "https://www.openstreetmap.org/relation/2978650"
	}
	}*/

	//output := "This service works"
	//encode := json.NewEncoder(w)
	/*err2 := encode.Encode(mappy)

	if err2 != nil {
		http.Error(w, "Error using "+err.Error(), http.StatusInternalServerError)
		return
	}
	log.Println(mappy)*/
}
