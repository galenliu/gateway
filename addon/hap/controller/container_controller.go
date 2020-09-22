package controller

import (
	"bytes"
	"encoding/json"
	"gateway/accessory"
	"io"
)

type AccessoriesController struct {
	container *accessory.Container
}

func NewAccessoriesController(container *accessory.Container) *AccessoriesController {
	return &AccessoriesController{container: container}
}

func (atr *AccessoriesController) HandleGetAccessories(i io.Reader) (io.Reader, error) {
	result, err := json.Marshal(atr.container)

	return bytes.NewBuffer(result), err
}
