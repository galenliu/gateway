'use strict';

const TuyaDevice = require('./tuya-device');
const PowerProperty = require('../properties/power-property');
const BrightnessProperty = require('../properties/brightness-property');

class PlugNightlightDevice extends TuyaDevice {
  constructor(adapter, cnf, cid) {
    super(adapter, cnf, cid);

    if (!('sockets' in this.ownconf)) {
      this.ownconf.sockets = 1;
    }

    this.name = cnf.name&&cnf.name!='' ? cnf.name : `Plug with night light`;
    this['@type'] = ['SmartPlug', 'Light', 'OnOffSwitch'];

    this.addProperty(new PowerProperty(this, {dps: this.ownconf.dps.on_light, default_dps: this.ownconf.sockets+1, num: '_light', label: 'Light'}));
    this.addProperty(new BrightnessProperty(this, {dps: this.ownconf.dps.brightness, default_dps: this.ownconf.sockets+2}));

    for (let i = 1; i <= this.ownconf.sockets; i++)
      this.addProperty(new PowerProperty(this, {dps: this.ownconf.dps[`on${i}`], default_dps: i, num: i}));
  }
}

module.exports = PlugNightlightDevice;
