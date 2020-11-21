'use strict';

const TuyaDevice = require('./tuya-device');
const PowerProperty = require('../properties/power-property');
const ModeProperty = require('../properties/mode-property');
const BrightnessProperty = require('../properties/brightness-property');
const LighttemperatureProperty = require('../properties/lighttemperature-property');
const ColourProperty = require('../properties/colourlight/colour-property');
const HueProperty = require('../properties/colourlight/hue-property');
const SaturationProperty = require('../properties/colourlight/saturation-property');
const LightnessProperty = require('../properties/colourlight/lightness-property');

class ColourLightDevice extends TuyaDevice {
  constructor(adapter, cnf, cid) {
    super(adapter, cnf, cid);

    this.name = cnf.name&&cnf.name!='' ? cnf.name : 'Colour Light';
    this['@type'] = ['Light', 'OnOffSwitch'];

    this.addProperty(new PowerProperty(this, {dps: this.ownconf.dps.on, default_dps: 1}));

    const modecb = ((mode) => {
      switch (mode) {
        case 'white':
          this.setVisibility(['on', 'mode', 'brightness', 'lighttemperature']);
          break;
        case 'colour':
          this.setVisibility(['on', 'mode', 'colour', 'hue', 'saturation', 'lightness']);
          break;
      }
      this.adapter.handleDeviceUpdated(this);
    }).bind(this);
    this.addProperty(new ModeProperty(this, {dps: this.ownconf.dps.mode, default_dps: 2, modes: ['white', 'colour'], modestxt: ['White', 'Colour'], setcb: modecb, updatecb: modecb}));

    this.addProperty(new BrightnessProperty(this, {dps: this.ownconf.dps.brightness, default_dps: 3}));
    this.addProperty(new LighttemperatureProperty(this, {dps: this.ownconf.dps.lighttemperature, default_dps: 4}));
    this.addProperty(new ColourProperty(this, {dps: this.ownconf.dps.colour, default_dps: 5}));
    this.addProperty(new HueProperty(this));
    this.addProperty(new SaturationProperty(this));
    this.addProperty(new LightnessProperty(this));
  }
}

module.exports = ColourLightDevice;
