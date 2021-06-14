package wot

type Link struct {
	Href   string      `json:"href,omitempty"`
	Type   string      `json:"type,omitempty"`
	Rel    string      `json:"rel,omitempty"`
	Anchor interface{} `json:"anchor,omitempty"`
}
