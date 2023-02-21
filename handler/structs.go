package handler

/*type UniInfo struct {
	Name      string    `json:"name"`
	Coutry    string    `json:"country"`
	Isocode   string    `json:"isocode"`
	Webpages  list.List `json:"webpages"`
	Languages Languages `json:"location"`
}*/
type UniInfo struct {
	WebPages []string `json:"web_pages,omitempty"`
	UniName  string   `json:"name,omitempty"`
	Country  string   `json:"country,omitempty"`
	Isocode  string   `json:"alpha_two_code"`
	CountryInfo
}

type CountryInfo struct {
	Name      CountryName       `json:"name"`
	Languages map[string]string `json:"languages"`
	Map       map[string]string `json:"maps"`
}

type UniAndCountryInfo struct {
	UniInfo
	CountryInfo
}

/*type Country struct {
	Name      CountryName       `json:"name"`
	Languages map[string]string `json:"languages"`
	Map       string            `json:"map"`
}*/

type CountryName struct {
	Common string `json:"common"`
}

type UniCountry struct {
	UniInfo
	CountryInfo
}

/*type UniCountry struct {
	WebPages  []string          `json:"web_pages,omitempty"`
	Name      string            `json:"name,omitempty"`
	Country   string            `json:"country,omitempty"`
	Isocode   string            `json:"isocode"`
	Webpages  []string          `json:"webpages"`
	Languages map[string]string `json:"languages"`
	Map       string            `json:"map"`
}*/
