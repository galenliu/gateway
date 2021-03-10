export const ThingType = {
    Alarm: "Alarm",
    AirQualitySensor: "AirQualitySensor",
    BarometricPressureSensor: "BarometricPressureSensor",
    BinarySensor: "BinarySensor",
    Camera: "Camera",
    ColorControl: "ColorControl",
    ColorSensor: "ColorSensor",
    DoorSensor: "DoorSensor",
    EnergyMonitor: "EnergyMonitor",
    HumiditySensor: "HumiditySensor",
    LeakSensor: "LeakSensor",
    Light: "Light",
    Lock: "Lock",
    MotionSensor: "MotionSensor",
    MultiLevelSensor: "MultiLevelSensor",
    MultiLevelSwitch: "MultiLevelSwitch",
    OnOffSwitch: "OnOffSwitch",
    PushButton: "PushButton",
    SmartPlug: "SmartPlug",
    SmokeSensor: "SmokeSensor",
    TemperatureSensor: "TemperatureSensor",
    Thermostat: "Thermostat",
    VideoCamera: "VideoCamera",

}

module.exports.CONNECTED = 'connected';
module.exports.DELETE_THING = 'deleteThing';
module.exports.DELETE_THINGS = 'deleteThings';
module.exports.EVENT_OCCURRED = 'eventOccurred';
module.exports.PROPERTY_STATUS = 'propertyStatus';
module.exports.REFRESH_THINGS = 'refreshThings';


export const ThingProperties = {
    OnOffProperty: "OnOffProperty",
}

export const AddonType = {
    Adapter: "adapter",
    Notifier: "notifier",
    Extension: "extension"
}

export const SettingsType = {
    Room: "room",
    Notifier: "notifier",
    Extension: "extension"
}

export const drawerWidth = 200;
