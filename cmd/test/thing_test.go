package test

import (
	"encoding/json"
	things "github.com/galenliu/gateway/api/models/container"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestThingMarshal(t *testing.T) {

	var data = `{"@context":"https://webthings.io/schemas","@type":["Thermostat"],"id":"virtual_thermostat","title":"virtual_thermostat","created":"2022-04-07T11:52:24.7399254+08:00","properties":{"coolingTargetTemperature":{"@type":"TargetTemperatureProperty","title":"Cooling Target","forms":[{"href":"things/virtual_thermostat/properties/coolingTargetTemperature","contentType":"application/json","op":["readproperty","writeproperty"]}],"unit":"degree celsius","type":"number","multipleOf":0.5},"heatingCooling":{"@type":"HeatingCoolingProperty","title":"Heating/Cooling","forms":[{"href":"things/virtual_thermostat/properties/heatingCooling","contentType":"application/json","op":["readproperty","writeproperty"]}],"enum":["off","heat","cool"],"type":"string"},"heatingTargetTemperature":{"@type":"TargetTemperatureProperty","title":"Heating Target","forms":[{"href":"things/virtual_thermostat/properties/heatingTargetTemperature","contentType":"application/json","op":["readproperty","writeproperty"]}],"unit":"degree celsius","type":"number","multipleOf":0.5},"temperature":{"@type":"TemperatureProperty","title":"Temperature","forms":[{"href":"things/virtual_thermostat/properties/temperature","contentType":"application/json","op":["readproperty"]}],"unit":"degree celsius","readOnly":true,"type":"number","multipleOf":0.1},"thermostatMode":{"@type":"ThermostatModeProperty","title":"Mode","forms":[{"href":"things/virtual_thermostat/properties/thermostatMode","contentType":"application/json","op":["readproperty","writeproperty"]}],"enum":["off","heat","cool","auto","dry","wind"],"type":"string"}},"forms":[{"href":"/things/virtual_thermostat/properties","contentType":"application/json","op":["readallproperties"]}],"selectedCapability":"Thermostat"}`

	var thing things.Thing
	err := json.Unmarshal([]byte(data), &thing)
	if err != nil {
		return
	}
	marshal, err := json.Marshal(thing)
	if err != nil {
		return
	}
	assert.Equal(t, string(marshal), data, "thing marshal ok")
	t.Log(string(marshal))
}
