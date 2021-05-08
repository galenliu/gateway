package services

import "github.com/galenliu/gateway/homekit/models/properties"

type Light struct {
	Properties map[string]properties.Property
}

func NewLightService(data []byte) *Light {
	return nil
}
