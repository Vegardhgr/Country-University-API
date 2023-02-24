package handler

/*type UniInfo struct {
	Name      string    `json:"name"`
	Coutry    string    `json:"country"`
	Isocode   string    `json:"isocode"`
	Webpages  list.List `json:"webpages"`
	Languages Languages `json:"location"`
}*/
type UniInfo struct {
	UniName  string   `json:"name"`
	Country  string   `json:"country"`
	Isocode  string   `json:"alpha_two_code"`
	WebPages []string `json:"web_pages"`
	CountryInfo
}

type CountryInfo struct {
	//Languages map[string]string `json:"languages"`
	//Country Country
	Languages map[string]string `json:"languages"`
	Map       string            `json:"map"`
	//Map       map[string]string `json:"maps"`
}

type Country struct {
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
