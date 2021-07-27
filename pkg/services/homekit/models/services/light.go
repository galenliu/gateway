package services

import (
	properties2 "github.com/galenliu/gateway/pkg/services/homekit/models/properties"
)

type Light struct {
	Properties map[string]properties2.Property
}

func NewLightService(data []byte) *Light {
	return nil
}
