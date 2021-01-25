const Utils = {


    sortCapabilities: (capabilities) => {
        // copy the array, as we're going to sort in place.
        const list = capabilities.slice();

        const priority = [
            'Lock',
            'Thermostat',
            'VideoCamera',
            'Camera',
            'SmartPlug',
            'Light',
            'MultiLevelSwitch',
            'OnOffSwitch',
            'ColorControl',
            'ColorSensor',
            'EnergyMonitor',
            'DoorSensor',
            'MotionSensor',
            'LeakSensor',
            'SmokeSensor',
            'PushButton',
            'TemperatureSensor',
            'HumiditySensor',
            'MultiLevelSensor',
            'Alarm',
            'BinarySensor',
            'BarometricPressureSensor',
            'AirQualitySensor',
        ];

        list.sort((a, b) => {
            if (!priority.includes(a) && !priority.includes(b)) {
                return 0;
            } else if (!priority.includes(a)) {
                return 1;
            } else if (!priority.includes(b)) {
                return -1;
            }

            return priority.indexOf(a) - priority.indexOf(b);
        });

        return list;
    },

    getClassFromCapability: (capability) => {
        switch (capability) {

            case 'Custom':
                return 'custom-thing';
            case 'OnOffSwitch':
                return 'on-off-switch';
            case 'MultiLevelSwitch':
                return 'multi-level-switch';
            case 'ColorControl':
                return 'color-control';
            case 'ColorSensor':
                return 'color-sensor';
            case 'EnergyMonitor':
                return 'energy-monitor';
            case 'BinarySensor':
                return 'binary-sensor';
            case 'MultiLevelSensor':
                return 'multi-level-sensor';
            case 'SmartPlug':
                return 'smart-plug';
            case 'Light':
                return 'light';
            case 'DoorSensor':
                return 'door-sensor';
            case 'MotionSensor':
                return 'motion-sensor';
            case 'LeakSensor':
                return 'leak-sensor';
            case 'SmokeSensor':
                return 'smoke-sensor';
            case 'PushButton':
                return 'push-button';
            case 'VideoCamera':
                return 'video-camera';
            case 'Camera':
                return 'camera';
            case 'TemperatureSensor':
                return 'temperature-sensor';
            case 'HumiditySensor':
                return 'humidity-sensor';
            case 'Alarm':
                return 'alarm';
            case 'Thermostat':
                return 'thermostat';
            case 'Lock':
                return 'lock';
            case 'BarometricPressureSensor':
                return 'barometric-pressure-sensor';
            case 'AirQualitySensor':
                return 'air-quality-sensor';
        }

        return '';
    }

}




module.exports = Utils
