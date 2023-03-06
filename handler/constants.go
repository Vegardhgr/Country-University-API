package handler

//Used for mock service:
/*const UNI_URL = "http://localhost:8081/unisearcher/v1/uniinfo/"
const COUNTRY_URL = "http://localhost:8081/unisearcher/v1/country/"*/

//Default paths to the APIs
const UNI_URL = "http://universities.hipolabs.com/"
const COUNTRY_URL = "https://restcountries.com/"

//Paths used for this service
const DEFAULT = "/"
const UNI_INFO = "uniinfo"
const NEIGHBOUR_UNIS = "neighbourunis"
const DIAG = "diag"
const UNI_INFO_PATH = "/unisearcher/v1/" + UNI_INFO + "/"
const NEIGHBOUR_UNIS_PATH = "/unisearcher/v1/" + NEIGHBOUR_UNIS + "/"
const DIAG_PATH = "/unisearcher/v1/" + DIAG + "/"

//Valid number of parts in each url
const VALID_NUMBER_OF_URL_PARTS_UNI_HANDLER = 4
const VALID_NUMBER_OF_URL_PARTS_NEIGHBOUR_UNIS_HANDLER = 5
const VALID_NUMBER_OF_URL_PARTS_DIAG = 4
