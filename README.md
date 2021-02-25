


### plugin API:


### App API说明:

***GET /things***  获取所有的things(json)


***协议： websocket  router:  /new_thing***
```json
{
  "id": "virtual-things-1",
  "title": "Virtual Multi-level Switch",
  "@context": "https://iot.mozilla.org/schemas",
  "@type": [
    "OnOffSwitch",
    "MultiLevelSwitch"
  ],
  "description": "",
  "properties": {
    "level": {
      "name": "level",
      "value": 0,
      "visible": true,
      "title": "Level",
      "type": "number",
      "@type": "LevelProperty",
      "unit": "percent",
      "minimum": 0,
      "maximum": 100,
      "readOnly": false,
      "links": [],
      "href": "/things/virtual-things-1/properties/level"
    },
    "on": {
      "name": "on",
      "value": false,
      "visible": true,
      "title": "On/Off",
      "type": "boolean",
      "@type": "OnOffProperty",
      "links": [],
      "href": "/things/virtual-things-1/properties/on"
    }
  },
  "actions": {},
  "events": {},
  "links": [],
  "baseHref": null,
  "pin": {
    "required": false,
    "pattern": ""
  },
  "credentialsRequired": false,
  "href": "/things/virtual-things-1"
}
```