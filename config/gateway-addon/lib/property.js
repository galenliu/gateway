"use strict";
/**
 * Property.
 *
 * Object which decscribes a property, and its value.
 *
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/.
 */
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.Property = void 0;
const assert_1 = __importDefault(require("assert"));
class Property {
    constructor(device, name, propertyDescr) {
        var _a, _b, _c, _d;
        this.fireAndForget = false;
        this.device = device;
        this.name = name;
        // The propertyDescr argument used to be the 'type' string, so we add an
        // assertion here to notify anybody who has an older plugin.
        assert_1.default.equal(typeof propertyDescr, 'object', 'Please update plugin to use property description.');
        const legacyDescription = propertyDescr;
        this.title = propertyDescr.title || legacyDescription.label;
        this.type = propertyDescr.type;
        this['@type'] = propertyDescr['@type'];
        this.unit = propertyDescr.unit;
        this.description = propertyDescr.description;
        this.minimum = (_a = propertyDescr.minimum) !== null && _a !== void 0 ? _a : legacyDescription.min;
        this.maximum = (_b = propertyDescr.maximum) !== null && _b !== void 0 ? _b : legacyDescription.max;
        this.enum = propertyDescr.enum;
        this.readOnly = propertyDescr.readOnly;
        this.multipleOf = propertyDescr.multipleOf;
        this.links = (_c = propertyDescr.links) !== null && _c !== void 0 ? _c : [];
        this.visible = (_d = propertyDescr.visible) !== null && _d !== void 0 ? _d : true;
    }
    /**
     * @returns a dictionary of useful information.
     * This is primarily used for debugging.
     */
    asDict() {
        return {
            name: this.name,
            value: this.value,
            visible: this.visible,
            title: this.title,
            type: this.type,
            '@type': this['@type'],
            unit: this.unit,
            description: this.description,
            minimum: this.minimum,
            maximum: this.maximum,
            enum: this.enum,
            readOnly: this.readOnly,
            multipleOf: this.multipleOf,
            links: this.links,
        };
    }
    /**
     * @returns the dictionary as used to describe a property. Currently
     * this does not include the href field.
     */
    asPropertyDescription() {
        return {
            title: this.title,
            type: this.type,
            '@type': this['@type'],
            unit: this.unit,
            description: this.description,
            minimum: this.minimum,
            maximum: this.maximum,
            enum: this.enum,
            readOnly: this.readOnly,
            multipleOf: this.multipleOf,
            links: this.links,
            visible: this.visible,
        };
    }
    /**
     * @method isVisible
     * @returns true if this is a visible property, which is a property
     *          that is reported in the property description.
     */
    isVisible() {
        return this.visible;
    }
    /**
     * Make the property visible or invisible
     */
    setVisible(visible) {
        this.visible = visible;
    }
    isFireAndForget() {
        return this.fireAndForget;
    }
    /**
     * Sets the value and notifies the device if the value has changed.
     * @returns true if the value has changed
     */
    setCachedValueAndNotify(value) {
        const oldValue = this.value;
        this.setCachedValue(value);
        // setCachedValue may change the value, therefore we have to check
        // this.value after the call to setCachedValue
        const hasChanged = oldValue !== this.value;
        if (hasChanged) {
            this.device.notifyPropertyChanged(this);
        }
        return hasChanged;
    }
    /**
     * Sets this.value and makes adjustments to ensure that the value
     * is consistent with the type.
     */
    setCachedValue(value) {
        if (this.type === 'boolean') {
            // Make sure that the value is actually a boolean.
            this.value = !!value;
        }
        else {
            this.value = value;
        }
        return this.value;
    }
    /**
     * @method getValue
     * @returns a promise which resolves to the retrieved value.
     *
     * This implementation is a simple one that just returns
     * the previously cached value.
     */
    getValue() {
        return new Promise((resolve) => {
            if (this.value != this.prevGetValue) {
                this.prevGetValue = this.value;
            }
            resolve(this.value);
        });
    }
    /**
     * @method setValue
     * @returns a promise which resolves to the updated value.
     *
     * @note it is possible that the updated value doesn't match
     * the value passed in.
     *
     * It is anticipated that this method will most likely be overridden
     * by a derived class.
     */
    setValue(value) {
        return new Promise((resolve, reject) => {
            if (this.readOnly) {
                reject('Read-only property');
                return;
            }
            const numberValue = value;
            // eslint-disable-next-line no-undefined
            if (typeof this.minimum !== 'undefined' && numberValue < this.minimum) {
                reject(`Value less than minimum: ${this.minimum}`);
                return;
            }
            // eslint-disable-next-line no-undefined
            if (typeof this.maximum !== 'undefined' && numberValue > this.maximum) {
                reject(`Value greater than maximum: ${this.maximum}`);
                return;
            }
            // eslint-disable-next-line no-undefined
            if (typeof this.multipleOf !== 'undefined' &&
                numberValue / this.multipleOf -
                    Math.round(numberValue / this.multipleOf) !== 0) {
                // note that we don't use the modulus operator here because it's
                // unreliable for floating point numbers
                reject(`Value is not a multiple of: ${this.multipleOf}`);
                return;
            }
            if (this.enum && this.enum.length > 0 &&
                !this.enum.includes(`${value}`)) {
                reject('Invalid enum value');
                return;
            }
            this.setCachedValueAndNotify(value);
            resolve(this.value);
        });
    }
    getDevice() {
        return this.device;
    }
    getName() {
        return this.name;
    }
    setName(value) {
        this.name = value;
    }
    getTitle() {
        return this.title;
    }
    setTitle(value) {
        this.title = value;
    }
    getType() {
        return this.type;
    }
    setType(value) {
        this.type = value;
    }
    getAtType() {
        return this['@type'];
    }
    setAtType(value) {
        this['@type'] = value;
    }
    getUnit() {
        return this.unit;
    }
    setUnit(value) {
        this.unit = value;
    }
    getDescription() {
        return this.description;
    }
    setDescription(value) {
        this.description = value;
    }
    getMinimum() {
        return this.minimum;
    }
    setMinimum(value) {
        this.minimum = value;
    }
    getMaximum() {
        return this.maximum;
    }
    setMaximum(value) {
        this.maximum = value;
    }
    getEnum() {
        return this.enum;
    }
    setEnum(value) {
        this.enum = value;
    }
    getReadOnly() {
        return this.readOnly;
    }
    setReadOnly(value) {
        this.readOnly = value;
    }
    getMultipleOf() {
        return this.multipleOf;
    }
    setMultipleOf(value) {
        this.multipleOf = value;
    }
    getLinks() {
        return this.links;
    }
    setLinks(value) {
        this.links = value;
    }
}
exports.Property = Property;
//# sourceMappingURL=property.js.map