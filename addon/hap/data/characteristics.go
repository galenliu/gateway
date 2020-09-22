package data

type Characteristics struct {
	Characteristics []Characteristic `json:"characteristics"`
}

type Characteristic struct {
	AccessoryID      uint64      `json:"aid"`
	CharacteristicID uint64      `json:"iid"`
	Value            interface{} `json:"value"`

	Status interface{} `json:"status,omitempty"`

	Events interface{} `json:"ev,omitempty"`
}
