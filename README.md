# assignment1

## Introduction 

Combines information about university and country from two different APIs. <br>
The APIs used to retrieve information are listed below:
<ul>
  <li>University API: http://universities.hipolabs.com/</li>  
  <li>Country API: https://restcountries.com/</li>  
</ul>

## Instructions
### Search for university
>Path: /unisearcher/v1/uniinfo/{:partial_or_complete_university_name}
#### Example:
>/unisearcher/v1/uniinfo/Norwegian University of Science and Technology
#### Response:
```json
[
    {
        "name": "Norwegian University of Science and Technology",
        "country": "Norway",
        "alpha_two_code": "NO",
        "web_pages": [
            "http://www.ntnu.no/"
        ],
        "languages": {
            "nno": "Norwegian Nynorsk",
            "nob": "Norwegian Bokmål",
            "smi": "Sami"
        },
        "map": "https://www.openstreetmap.org/relation/2978650"
    }
]
```
### Search for universities in bordering countries
>Path: /unisearcher/v1/neighbourunis/{:country_name}/{:partial_or_complete_university_name}{?limit={:number}}
#### Example:
>/unisearcher/v1/neighbourunis/russia/univeristy?limit=2
#### Response:
```json
[
  {
    "name": "Kharkiv National University of Economics",
    "country": "Ukraine",
    "alpha_two_code": "UA",
    "web_pages": [
      "http://www.hneu.edu.ua/"
    ],
    "languages": {
      "ukr": "Ukrainian"
    },
    "map": "https://www.openstreetmap.org/relation/60199"
  },
  {
    "name": "University of Technology and Life Sciences",
    "country": "Poland",
    "alpha_two_code": "PL",
    "web_pages": [
      "http://www.utp.edu.pl/"
    ],
    "languages": {
      "pol": "Polish"
    },
    "map": "https://www.openstreetmap.org/relation/49715"
  }
]
```
### Diagnostics
>Path: /unisearcher/v1/diag<br>
#### Response:
```json
{
    "universitiesapi": "200",
    "countriesapi": "200",
    "version": "V1",
    "uptime": 370.679817207
}
```
Gets information about the APIs that are used, the version of this service, as well as time since last service restart.

## Decisions 
#### Country name issue
The country name in http://universities.hipolabs.com/, and the country name in https://restcountries.com/, are not always similar.
Therefore, it is not a good option to retrieve both universities and country based on one specific country name, as the name can vary. In the university API, it is 
not an option to search for universities by alpha code. For that reason, when retrieving universities, all universities are retrieved from the university API
based on the university name provided by the user. Then in the code, universities are filtered based on the country's alpha code.<br>

Another way to address this issue is to retrieve all the countries and all the universities from their respective APIs at every service restart.
If the alpha codes match, but the country names differ between the two APIs, add the country name from the country API as key, and add the country name from the
university API as value. This way, a mapping between a country with different names in the two APIs has been created.<br>

The reasoning for why the first solution was chosen, is that the second solution can be vulnerable for changes in the APIs.
If one of the APIs change a country name, and the map in this service has not yet been updated, the mapping will not work as intended anymore.
The first solution can therefore be more robust in terms of changes in the APIs that are used.<br>

The first solution is still not perfect. Since all universities with a specific name is retrieved, there is a great chance that a lot of
unnecessary data will be collected and processed. <br>

---

#### Error handling
In golang, it seems to be no prescribed approach for managing erros. If an error appears some where in the code, it could either be returned back to 
the top, based on the same concept as exceptions in java, or the error could be managed where it happens. Since exceptions is not a thing in golang, it 
does not seem to be any good ways to return and handle errors at the top of the code. Therefore, in this assignment, error handling is done where the error appears.
This may not be the best way to deal with errors, but it was more satisfying to just return a boolean based on whether there was an error and deal with it where it
appeared, compared to always needing to return an error and then check if the returned error was equal to a specific error string.

