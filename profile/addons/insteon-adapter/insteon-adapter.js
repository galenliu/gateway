
'use strict';

const {
  Adapter,
  Device,
  Event,
  Property,
} = require('gateway-addon');

const { findProductInfo, MessageCommands } = require('insteon-plm');
const fs = require('fs');
const mkdirp = require('mkdirp');
const os = require('os');
const path = require('path');
const storage = require('node-persist');
const delay = require('delay');

const InsteonAPIHandler = require('./insteon-api-handler');
const { findCapabilities } = require('./insteon-capabilities');

const ID_PREFIX = 'insteon-';
const POLL_INTERVAL_MS = 60 * 60 * 1000;
const HEARTBEAT_GROUP = 4;
const LOW_BATTERY_GROUP = 3;

function getDataPath() {
  let profileDir;
  if (process.env.hasOwnProperty('MOZIOT_HOME')) {
    profileDir = process.env.MOZIOT_HOME;
  } else {
    profileDir = path.join(os.homedir(), '.mozilla-iot');
  }

  return path.join(profileDir, 'data', 'insteon');
}

function addressToId(address) {
  return ID_PREFIX + address;
}

function idToAddress(id) {
  return id.substring(ID_PREFIX.length);
}

class InsteonProperty extends Property {
  get hub() {
    return this.device.hub;
  }

  get messenger() {
    return this.device.messenger;
  }

  handleMessage(_message) {
  }

  async pollWithDelay(ms = 500) {
    await delay(ms);
    this.poll().catch((e) => {
      console.warn(`Failed to poll ${this.device.address}: ${e.message}`);
    });
  }

  poll() {
  }
}

class MotionProperty extends InsteonProperty {
  constructor(device) {
    super(device, 'motion', {
      '@type': 'MotionProperty',
      label: 'Motion',
      name: 'motion',
      type: 'boolean',
      value: false,
      readOnly: true,
    });
  }

  handleMessage(message) {
    const { cmd1 } = message;

    if (cmd1 === MessageCommands.ON) {
      this.setCachedValueAndNotify(true);
    } else if (cmd1 === MessageCommands.OFF) {
      this.setCachedValueAndNotify(false);
    }
  }
}

class OnOffProperty extends InsteonProperty {
  constructor(device, { readOnly = true } = {}) {
    super(device, 'on', {
      '@type': 'OnOffProperty',
      label: 'On/Off',
      name: 'on',
      type: 'boolean',
      value: false,
      readOnly,
    });
  }

  handleMessage(message) {
    const { flags: { ack }, cmd1 } = message;

    // Ignore ACKs. The `SettableOnOffProperty` updates the value
    // after a successful set.
    if (ack) {
      return;
    }

    // The IOLinc sends the same message whether it was the relay or sensor
    // that changed. This is problematic for us, but the relay message
    // is only sent if someone pushes the button on the device itself.
    // This should be a rare occurrence, so we'll treat on/off notifications
    // as coming from the sensor and ignore them here.
    if (this.device.isIOLinc && (cmd1 === MessageCommands.ON || cmd1 === MessageCommands.OFF)) {
      return;
    }

    if (cmd1 === MessageCommands.ON) {
      this.setCachedValueAndNotify(true);
    } else if (cmd1 === MessageCommands.ON_FAST) {
      this.setCachedValueAndNotify(true);
      this.device.eventNotify(new Event(this.device, 'FastOn'));
    } else if (cmd1 === MessageCommands.OFF) {
      this.setCachedValueAndNotify(false);
    } else if (cmd1 === MessageCommands.OFF_FAST) {
      this.setCachedValueAndNotify(false);
      this.device.eventNotify(new Event(this.device, 'FastOff'));
    }
  }
}

class SettableOnOffProperty extends OnOffProperty {
  constructor(device) {
    super(device, { readOnly: false });
  }

  async poll() {
    const level = await this.hub.status(this.device.address);
    this.setCachedValueAndNotify(level > 0);
  }

  async setValue(value) {
    if (value) {
      await this.hub.turnOn(this.device.address);
    } else {
      await this.hub.turnOff(this.device.address);
    }

    this.setCachedValueAndNotify(value);

    if (value && this.device.isIOLinc) {
      // The IOLinc may be in a momentary mode, which means if we
      // just turned it on, it may turn off in a bit. Poll for
      // that change.
      this.pollWithDelay(3000);
    }

    return value;
  }
}

class SensorProperty extends InsteonProperty {
  constructor(device) {
    super(device, 'active', {
      '@type': 'BooleanProperty',
      label: 'Sensor',
      name: 'active',
      type: 'boolean',
      value: false,
      readOnly: true,
    });
  }

  handleMessage(message) {
    const { flags: { ack }, cmd1 } = message;

    if (ack) {
      return;
    }

    // OFF means open, which seems weird.
    if (cmd1 === MessageCommands.OFF) {
      this.setCachedValueAndNotify(true);
    } else if (cmd1 === MessageCommands.ON) {
      this.setCachedValueAndNotify(false);
    }
  }

  async poll() {
    if (this.device.isIOLinc) {
      const value = await this.hub.status(this.device.address, { cmd2: 1 });
      this.setCachedValueAndNotify(value > 0);
    }
  }
}

class DimmerProperty extends InsteonProperty {
  constructor(device) {
    super(device, 'level', {
      '@type': 'LevelProperty',
      label: 'Level',
      name: 'level',
      type: 'integer',
      minimum: 0,
      maximum: 100,
      value: 0,
    });
  }

  handleMessage(message) {
    const { cmd1 } = message;
    if (cmd1 === MessageCommands.ON || cmd1 === MessageCommands.ON_FAST) {
      this.pollWithDelay();
    } else if (cmd1 === MessageCommands.OFF || cmd1 === MessageCommands.OFF_FAST) {
      this.setCachedValueAndNotify(0);
    } else if (cmd1 === MessageCommands.STOP_MANUAL_CHANGE) {
      this.poll();
    }
  }

  async poll() {
    const level = await this.hub.status(this.device.address);
    this.setCachedValueAndNotify(Math.round((level / 255) * 100));
  }

  async setValue(level) {
    if (level > 0) {
      this.hub.turnOn(this.device.address, { level });
    } else {
      this.hub.turnOff(this.device.address);
    }

    this.setCachedValueAndNotify(level);

    return level;
  }
}

class OpenProperty extends InsteonProperty {
  constructor(device) {
    super(device, 'open', {
      '@type': 'OpenProperty',
      label: 'Open/Closed',
      name: 'open',
      type: 'boolean',
      readOnly: true,
      value: false,
    });
  }

  handleMessage(message) {
    const { cmd1 } = message;

    if (cmd1 === MessageCommands.ON) {
      this.setCachedValueAndNotify(true);
    } else if (cmd1 === MessageCommands.OFF) {
      this.setCachedValueAndNotify(false);
    }
  }
}

class InsteonDevice extends Device {
  constructor(adapter, id, deviceDescription) {
    super(adapter, id);
    this.title = deviceDescription.name || deviceDescription.title;
    this.description = deviceDescription.description;
    this.category = deviceDescription.category;
    this.subcategory = deviceDescription.subcategory;
    this._address = idToAddress(id);
    const type = this['@type'] = deviceDescription['@type'] || [];

    let isLight = false;
    let isPollable = false;
    if (type.includes('OnOffSwitch')) {
      if (type.includes('BatteryPowered')) {
        this.properties.set('on', new OnOffProperty(this));
      } else {
        this.properties.set('on', new SettableOnOffProperty(this));
        isPollable = true;
      }
      isLight = true;
    }

    if (type.includes('MultiLevelSwitch')) {
      this.properties.set('level', new DimmerProperty(this));
      isLight = true;
      isPollable = true;

      this.addAction('Fade', {
        '@type': 'FadeAction',
        input: {
          type: 'object',
          required: ['level', 'duration'],
          properties: {
            level: {
              type: 'integer',
              unit: 'percent',
              minimum: 0,
              maximum: 100,
            },
            duration: {
              type: 'integer',
              unit: 'second',
              minimum: 0,
            },
          },
        },
      });
    }

    if (type.includes('BatteryPowered')) {
      this.addEvent('Heartbeat', {
        '@type': 'HeartbeatEvent',
        description: 'Device checked in',
      });

      this.addEvent('LowBattery', {
        '@type': 'LowBatteryEvent',
        description: 'Device has notified of a low battery',
      });
    }

    if (type.includes('DoorSensor')) {
      this.properties.set('open', new OpenProperty(this));
    }

    if (type.includes('MotionSensor')) {
      this.properties.set('motion', new MotionProperty(this));
    }

    if (type.includes('BinarySensor')) {
      this.properties.set('active', new SensorProperty(this));
    }

    // Add events and actions for fast on/off
    if (isLight) {
      this.addEvent('FastOn', {
        '@type': 'DoublePressedEvent',
        description: 'Double tap on',
      });

      this.addEvent('FastOff', {
        '@type': 'DoublePressedEvent',
        description: 'Double tap off',
      });
    }

    if (isPollable) {
      this.addAction('Poll', {
        '@type': 'PollAction',
      });
    }
  }

  async performAction(action) {
    switch (action.name) {
      case 'Fade': {
        const { duration, level } = action.input;
        if (level > 0) {
          await this.hub.turnOn(this.address, { level, duration });
        } else {
          await this.hub.turnOff(this.address, { duration });
        }

        setTimeout(() => this.poll(), duration * 1000);
        break;
      }
      case 'Poll': {
        await this.poll();
        break;
      }
      default:
        console.warn(`Unknown action ${action.name}`);
        break;
    }
  }

  get address() {
    return this._address;
  }

  get isBatteryPowered() {
    return this['@type'].includes('BatteryPowered');
  }

  get isIOLinc() {
    return this.category === 7 && this.subcategory === 0;
  }

  async handleMessage(message) {
    for (const prop of this.properties.values()) {
      await prop.handleMessage(message);
    }

    if (this.isBatteryPowered) {
      const { flags: { broadcast }, to } = message;
      if (broadcast && parseInt(to, 16) === HEARTBEAT_GROUP) {
        this.eventNotify(new Event(this, 'Heartbeat'));
      } else if (broadcast && parseInt(to, 16) === LOW_BATTERY_GROUP) {
        this.eventNotify(new Event(this, 'LowBattery'));
      }
    }
  }

  async poll() {
    for (const prop of this.properties.values()) {
      await prop.poll();
    }
  }

  async link({ controller, group = 1 }) {
    return this.hub.link(this.address, { controller, group });
  }

  get hub() {
    return this.adapter.hub;
  }

  get messenger() {
    return this.adapter.messenger;
  }
}

class InsteonAdapter extends Adapter {
  constructor(addonManager, manifest, hub) {
    super(addonManager, 'insteon', manifest.name);
    this.name = 'INSTEON';

    const dataDir = getDataPath();
    if (!fs.existsSync(dataDir)) {
      mkdirp.sync(dataDir, { mode: 0o755 });
    }

    this._hub = hub;

    storage.init({ dir: dataDir });
    addonManager.addAdapter(this);

    this.apiHandler = new InsteonAPIHandler(addonManager, this);

    this.messenger.on('message', (message) => this._handleMessage(message));

    // Poll devices periodically
    setInterval(() => this.poll(), POLL_INTERVAL_MS);
  }

  get hub() {
    return this._hub;
  }

  get messenger() {
    return this.hub.messenger;
  }

  async poll() {
    for (const device of Object.values(this.devices)) {
      try {
        await device.poll();
      } catch (e) {
        console.warn('Failed to poll ${device.id}', e);
      }
    }
  }

  _handleMessage(message) {
    const device = this.devices[addressToId(message.from)];
    if (device) {
      device.handleMessage(message);
    }
  }

  createDeviceByAddress(address, category, subcategory) {
    const info = findProductInfo(category, subcategory);
    const name = (info && info.name) || 'INSTEON Device';

    return this.createDevice(addressToId(address), {
      title: `${name} [${address}]`,
      category,
      subcategory,
    });
  }

  createDevice(deviceId, { title, category, subcategory }) {
    const caps = findCapabilities(category, subcategory);
    if (caps.length === 0) {
      throw new Error(`Unsupported device: category=${category}, subcategory=${subcategory}`);
    }

    return new InsteonDevice(this, deviceId, {
      '@type': caps,
      title,
      category,
      subcategory,
    });
  }

  addDevice(device) {
    if (device.id in this.devices) {
      return this.devices[device.id];
    }

    this.devices[device.id] = device;
    this.handleDeviceAdded(device);

    device.poll().catch(() => {
      console.warn(`Failed initial poll of ${device.id}`);
    });
  }

  async removeDevice(deviceId) {
    if (deviceId in this.devices) {
      const device = this.devices[deviceId];
      delete this.devices[deviceId];

      await storage.removeItem(deviceId);

      this.handleDeviceRemoved(device);
    }
  }

  async handleDeviceSaved(deviceId, device) {
    if (!deviceId.startsWith(ID_PREFIX) || this.devices[deviceId]) {
      return;
    }

    const info = await storage.getItem(idToAddress(deviceId));
    if (!info) {
      return;
    }

    const { category, subcategory } = info;
    const newDevice = this.createDevice(deviceId, {
      title: device.title,
      category,
      subcategory,
    });

    if (newDevice) {
      this.addDevice(newDevice);
    }
  }

  async startPairing(timeoutSeconds) {
    super.startPairing(timeoutSeconds);

    this.link({ timeout: timeoutSeconds * 1000 }).then(async (device) => {
      if (device) {
        this.addDevice(device);
      }
    }).catch((e) => {
      console.error('Pairing failed', e);
      this.sendPairingPrompt(`Pairing failed: ${e.message}`);
    });
  }

  async cancelPairing() {
    super.cancelPairing();

    return this.hub.messenger.modem.cancelAllLinking();
  }

  async removeThing(device) {
    try {
      await this.removeDevice(device.id);
      console.log(`Device ${device.id} was removed.`);
    } catch (err) {
      console.error(`Failed to remove device ${device.id}`, err);
    }
  }

  async unload() {
    await this.hub.close();
  }

  async link({ timeout = 30000 } = {}) {
    this.sendPairingPrompt('Press the \'set\' button on the device you would like to pair.');
    const info = await this.hub.link(null, { timeout });
    if (!info) {
      console.warn('No link information received.');
      return null;
    }

    const existing = this.devices[addressToId(info.address)];

    if (!info.controller) {
      if (existing) {
        this.sendPairingPrompt('Paired existing device as responder', null, existing);
        console.log('Paired device as responder', existing);
        return existing;
      } else {
        this.sendPairingPrompt(`Paired device as responder: ${info.address}`);
        console.log(`Paired device as responder: ${info.address}`);
        return null;
      }
    }

    try {
      const device = this.createDeviceByAddress(info.address, info.category, info.subcategory);

      if (!device['@type'].includes('BatteryPowered')) {
        // For non-battery devices, also try to link as a responder. Not fatal if this fails.
        try {
          await this.hub.link(info.address, { controller: false });
        } catch (e) {
          console.warn('Unable to link device as responder', e);
        }
      }

      await storage.setItem(info.address, info);
      return device;
    } catch (e) {
      this.sendPairingPrompt(`Failed to pair: ${e.message}`);
      console.error('Failed to pair device', e);
    }

    return null;
  }

  async scan() {
    console.log('Attempting to add devices from adapter link database...');

    const addresses = new Set((await this.hub.links()).map((info) => info.address));

    console.log('Found addresses from database', addresses);

    const devices = [];
    for (const address of addresses) {
      if (this.devices[addressToId(address)]) {
        // Skip known devices
        continue;
      }

      try {
        const info = await this.hub.productInfo(address);
        await storage.setItem(address, info);

        const device = await this.createDeviceByAddress(address, info.category, info.subcategory);
        this.addDevice(device);
        devices.push(device.asDict());
      } catch (e) {
        console.error(`Failed to add device ${address}`, e);
      }
    }

    console.log(`Added ${devices.length} devices from database`);
    return devices;
  }
}

module.exports = InsteonAdapter;
