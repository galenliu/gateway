package models

//func toThing(device *addon.Device) *thing.Thing {
//	thing := &thing.Thing{
//		ID:                  device.ID,
//		AtContext:           device.AtContext,
//		AtType:              device.AtType,
//		Title:               device.Title,
//		Description:         device.Description,
//		Links:               nil,
//		BaseHref:            "",
//		Href:                fmt.Sprintf("/%s", device.ID),
//		CredentialsRequired: false,
//		SelectedCapability:  "",
//		Properties:          nil,
//		Actions:             nil,
//		Events:              nil,
//	}
//	var tps = make(map[string]*thing.Property)
//	for _, p := range device.Properties {
//		tp := toThingProperty(p, thing.ID)
//		tp.ThingId = device.ID
//		tps[tp.Name] = tp
//	}
//	thing.Properties = tps
//	return thing
//}
//
//func toThings(devices []*addon.Device) []*thing.Thing {
//	if devices == nil {
//		return nil
//	}
//	var things = make([]*thing.Thing, 100)
//	for _, dev := range devices {
//		if dev != nil {
//			var thing = toThing(dev)
//			things = append(things, thing)
//		}
//	}
//	return things
//}
//
//func toThingProperty(p *addon.Property, deviceId string) *thing.Property {
//	prop := &thing.Property{
//		Name:        p.Name,
//		AtType:      p.AtType,
//		Type:        p.Type,
//		Title:       p.Title,
//		Description: p.Description,
//		Unit:        p.Unit,
//		ReadOnly:    p.ReadOnly,
//		Visible:     p.Visible,
//		Minimum:     p.Minimum,
//		Maximum:     p.Maximum,
//		Value:       p.Value,
//		Href:        fmt.Sprintf("/things/%s/properties/%s", deviceId, p.Name),
//	}
//	return prop
//}
