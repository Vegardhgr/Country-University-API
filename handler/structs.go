package handler

// UniInfo
//Relevant university and country fields
type UniInfo struct {
	UniName  string   `json:"name"`
	Country  string   `json:"country"`
	Isocode  string   `json:"alpha_two_code"`
	WebPages []string `json:"web_pages"`
	CountryInfo
}

// CountryInfo
//Relevant country fields
type CountryInfo struct {
	Languages map[string]string `json:"languages"`
	StreetMap string            `json:"map"`
}

//Country
//Fields that are necessary from the country api, but these are not sent to the user
type Country struct {
	Name      CountryName       `json:"name"`
	Cca2      string            `json:"cca2"`
	Languages map[string]string `json:"languages"`
	StreetMap map[string]string `json:"maps"`
}

// Borders
//A struct for bordering countries
type Borders struct {
	Borders []string `json:"borders"`
}

// CountryName
//A struct for a countries official name
type CountryName struct {
	Official string `json:"official"`
}

// Diag
//A struct for diag
type Diag struct {
	UnisApiStatus      string  `json:"universitiesapi"`
	CountriesApiStatus string  `json:"countriesapi"`
	Version            string  `json:"version"`
	Uptime             float64 `json:"uptime"`
}
