package models

type Thing struct {
	AtContext    string   `json:"@context" valid:"url,required"`
	AtType       []string `json:"@type" valid:"url"`
	Title        string   `json:"title,required"`
	Titles       []string `json:"titles,omitempty"`
	Description  string   `json:"description,omitempty"`
	ID           string   `json:"id"`
	Descriptions []string `json:"descriptions,omitempty"`
	Version      string   `json:"version,omitempty"`
	Created      string   `json:"created,omitempty"`
	Modified     string   `json:"modified,omitempty"`
	Support      string   `json:"support,omitempty"`
	Properties   map[string]*Property
}

