package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func UniversityHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")

	uniInfoOutput, err := ioutil.ReadFile("handler/UniversityInfo.json")

	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	//mappy := make(map[interface{}]interface{})
	//var uni []interface{}
	var unis []UniInfo
	json.Unmarshal(uniInfoOutput, &unis)

	for i := range unis {
		json.NewEncoder(w).Encode(unis[i])
	}

	/*name := strings.ReplaceAll(strings.Split(r.URL.String(), "/")[4], "%20", " ")

	retrieveUniInfoUrl, _ := http.Get("http://universities.hipolabs.com/search?name=" + name)

	/*var uunis []UniInfo

	infoFromSiteAsArray, _ := ioutil.ReadAll(retrieveUniInfoUrl.Body)

	json.Unmarshal(infoFromSiteAsArray, &uunis)

	for i := range uunis {
		json.NewEncoder(w).Encode(uunis[i])
	}*/

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
