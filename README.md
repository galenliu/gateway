//git submodule update –init –recursive

### plugin API:

### App API说明:

***GET /things***  获取所有的things(json)

***协议： websocket router:  /new_thing***

```json
{
  "@context": ["https://webthings.io/schemas"],
  "title": "Light192.168.0.102",
  "id": "/things/Light192.168.0.102",
  "@type": ["Light"],
  "properties": {
    "level": {
      "@type": "BrightnessProperty",
      "title": "Brightness",
      "forms": [{
        "href": "/things/Light192.168.0.102/properties/level"
      }],
      "type": "integer",
      "unit": "percent",
      "writeOnly": false,
      "schema": {
        "minimum": 100
      },
      "name": "level",
      "readOnly": false,
      "visible": false,
      "minimum": 0,
      "maximum": 100,
      "thingId": ""
    },
    "on": {
      "@type": "OnOffProperty",
      "title": "On/Off",
      "forms": [{
        "href": "/things/Light192.168.0.102/properties/on"
      }],
      "type": "boolean",
      "writeOnly": false,
      "schema": null,
      "name": "on",
      "readOnly": false,
      "visible": false,
      "thingId": ""
    },
    "hue": {
      "@type": "ColorProperty",
      "title": "Hue",
      "forms": [{
        "href": "/things/Light192.168.0.102/properties/hue"
      }],
      "type": "string",
      "writeOnly": false,
      "schema": null,
      "name": "hue",
      "readOnly": false,
      "visible": false,
      "thingId": ""
    }
  },
  "forms": [{
    "rel": "properties",
    "href": " /things/Light192.168.0.102/properties"
  }, {
    "href": "/things/Light192.168.0.102",
    "rel": "alternate",
    "mediaType": "text/html"
  }, {
    "rel": "alternate",
    "href": "wss://localhost/things/Light192.168.0.102"
  }],
  "Pin": {
    "required": false
  },
  "selectedCapability": "Light",
  "connected": false
}
```