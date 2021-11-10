package model

//func toThing(addon *addon.Device) *thing.Thing {
//	thing := &thing.Thing{
//		Id:                  addon.Id,
//		AtContext:           addon.AtContext,
//		AtType:              addon.AtType,
//		Title:               addon.Title,
//		Description:         addon.Description,
//		Links:               nil,
//		BaseHref:            "",
//		Href:                fmt.Sprintf("/%s", addon.Id),
//		CredentialsRequired: false,
//		SelectedCapability:  "",
//		properties:          nil,
//		actions:             nil,
//		events:              nil,
//	}
//	var tps = make(map[string]*thing.Property)
//	for _, p := range addon.properties {
//		tp := toThingProperty(p, thing.Id)
//		tp.ThingId = addon.Id
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
