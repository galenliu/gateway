package model

//func toThing(device *addon.Device) *thing.Thing {
//	thing := &thing.Thing{
//		Id:                  device.Id,
//		AtContext:           device.AtContext,
//		AtType:              device.AtType,
//		Title:               device.Title,
//		Description:         device.Description,
//		Links:               nil,
//		BaseHref:            "",
//		Href:                fmt.Sprintf("/%s", device.Id),
//		CredentialsRequired: false,
//		SelectedCapability:  "",
//		properties:          nil,
//		actions:             nil,
//		events:              nil,
//	}
//	var tps = make(map[string]*thing.Property)
//	for _, p := range device.properties {
//		tp := toThingProperty(p, thing.Id)
//		tp.ThingId = device.Id
//		tps[tp.name] = tp
//	}
//	thing.properties = tps
//	return thing
//}
//
//func toThings(devices []*addon.Device) []*thing.Thing {
//	if devices == nil {
//		return nil
//	}
//	var instance = make([]*thing.Thing, 100)
//	for _, dev := range devices {
//		if dev != nil {
//			var thing = toThing(dev)
//			instance = append(instance, thing)
//		}
//	}
//	return instance
//}
//
//func toThingProperty(p *addon.Property, deviceId string) *thing.Property {
//	prop := &thing.Property{
//		name:        p.name,
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
//		Href:        fmt.Sprintf("/instance/%s/properties/%s", deviceId, p.name),
//	}
//	return prop
//}
