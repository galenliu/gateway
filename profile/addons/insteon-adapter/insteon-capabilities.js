
// Map device category and subcategory to capabilities
const CAPABILITIES = {
  0: {
    // Mini remote, 1 scene
    17: ['OnOffSwitch', 'BatteryPowered'],
  },
  1: {
    // Dimmer switches
    default: ['MultiLevelSwitch', 'OnOffSwitch'],
  },
  2: {
    // Toggle switches
    default: ['OnOffSwitch'],
  },
  7: {
    // IOLinc
    0: ['OnOffSwitch', 'BinarySensor'],
  },
  16: {
    1: ['MotionSensor', 'BatteryPowered'],
    2: ['DoorSensor', 'BatteryPowered'],
    17: ['DoorSensor', 'BatteryPowered'],
    22: ['MotionSensor', 'BatteryPowered'],
  },
};

function findCapabilities(category, subcategory) {
  const categoryInfo = CAPABILITIES[category];
  if (!categoryInfo) {
    return [];
  }

  return categoryInfo[subcategory] || categoryInfo.default || [];
}

module.exports = { findCapabilities };
