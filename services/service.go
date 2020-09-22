package services

import "smartassistant/properties"

type Service struct {
	Description        string                `json:"description"`
	Type               string                `json:"type"`
	RequiredProperties []properties.Property `json:"required_properties"`
	OptionalProperties []properties.Property `json:"optional_properties"`
}
