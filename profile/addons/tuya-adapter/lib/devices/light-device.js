'use strict';

const TuyaDevice = require('./tuya-device');
const PowerProperty = require('../properties/power-property');
const BrightnessProperty = require('../properties/brightness-property');
const LighttemperatureProperty = require('../properties/lighttemperature-property');

class LightDevice extends TuyaDevice {
  constructor(adapter, cnf, cid) {
    super(adapter, cnf, cid);

    this['@type'] = ['Light', 'OnOffSwitch'];

    this.addProperty(new PowerProperty(this, {dps: this.ownconf.dps.on, default_dps: 1}));
    this.addProperty(new BrightnessProperty(this, {dps: this.ownconf.dps.brightness, default_dps: 2}));
    this.addProperty(new LighttemperatureProperty(this, {dps: this.ownconf.dps.lighttemperature, default_dps: 3}));
  }
}

module.exports = LightDevice;
