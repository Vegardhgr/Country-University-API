package handler

type Languages struct {
	NNO string `json:"nno"`
	NOB string `json:"nob"`
	SMI string `json:"smi"`
}

/*type UniInfo struct {
	Name      string    `json:"name"`
	Coutry    string    `json:"country"`
	Isocode   string    `json:"isocode"`
	Webpages  list.List `json:"webpages"`
	Languages Languages `json:"location"`
}*/
type UniInfo struct {
	WebPages      []string    `json:"web_pages,omitempty"`
	Name          string      `json:"name,omitempty"`
	AlphaTwoCode  string      `json:"alpha_two_code,omitempty"`
	StateProvince interface{} `json:"state-province,omitempty"`
	Domains       []string    `json:"domains,omitempty"`
	Country       string      `json:"country,omitempty"`
}

type UniAndContry struct {
	`"name`,
	`"country"`,
	`"isocode"`,
	`"webpages"`,
	`"languages"`,
}
