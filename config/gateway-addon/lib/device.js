"use strict";
/**
 * Device Model.
 *
 * Abstract base class for devices managed by an adapter.
 *
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/.
 */
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.Device = void 0;
const action_1 = require("./action");
const ajv_1 = __importDefault(require("ajv"));
const ajv = new ajv_1.default();
class Device {
    constructor(adapter, id) {
        this['@context'] = 'https://webthings.io/schemas';
        this['@type'] = [];
        this.name = '';
        this.title = '';
        this.description = '';
        this.properties = new Map();
        this.actions = new Map();
        this.events = new Map();
        this.links = [];
        this.pinRequired = false;
        this.credentialsRequired = false;
        this.adapter = adapter;
        this.id = `${id}`;
    }
    mapToDict(map) {
        const dict = {};
        map.forEach((property, propertyName) => {
            dict[propertyName] = Object.assign({}, property);
        });
        return dict;
    }
    mapToDictFromFunction(map) {
        const dict = {};
        map.forEach((property, propertyName) => {
            dict[propertyName] = property.asDict();
        });
        return dict;
    }
    asDict() {
        return {
            id: this.id,
            title: this.title || this.name,
            '@context': this['@context'],
            '@type': this['@type'],
            description: this.description,
            properties: this.mapToDictFromFunction(this.properties),
            actions: this.mapToDict(this.actions),
            events: this.mapToDict(this.events),
            links: this.links,
            baseHref: this.baseHref,
            pin: {
                required: this.pinRequired,
                pattern: this.pinPattern,
            },
            credentialsRequired: this.credentialsRequired,
        };
    }
    /**
     * @returns this object as a thing
     */
    asThing() {
        return {
            id: this.id,
            title: this.title || this.name,
            '@context': this['@context'],
            '@type': this['@type'],
            description: this.description,
            properties: this.mapToDictFromFunction(this.properties),
            actions: this.mapToDict(this.actions),
            events: this.mapToDict(this.events),
            links: this.links,
            baseHref: this.baseHref,
            pin: {
                required: this.pinRequired,
                pattern: this.pinPattern,
            },
            credentialsRequired: this.credentialsRequired,
        };
    }
    debugCmd(cmd, params) {
        console.log('Device:', this.name, 'got debugCmd:', cmd, 'params:', params);
    }
    getId() {
        return this.id;
    }
    /**
     * @deprecated Please use getTitle()
     */
    getName() {
        console.log('getName() is deprecated. Please use getTitle().');
        return this.getTitle();
    }
    getTitle() {
        if (this.name && !this.title) {
            this.title = this.name;
        }
        return this.title;
    }
    getPropertyDescriptions() {
        const propDescs = {};
        this.properties.forEach((property, propertyName) => {
            if (property.isVisible()) {
                propDescs[propertyName] = property.asPropertyDescription();
            }
        });
        return propDescs;
    }
    findProperty(propertyName) {
        return this.properties.get(propertyName);
    }
    addProperty(property) {
        this.properties.set(property.getName(), property);
    }
    /**
     * @method getProperty
     * @returns a promise which resolves to the retrieved value.
     */
    getProperty(propertyName) {
        return new Promise((resolve, reject) => {
            const property = this.findProperty(propertyName);
            if (property) {
                property.getValue().then((value) => {
                    resolve(value);
                });
            }
            else {
                reject(`Property "${propertyName}" not found`);
            }
        });
    }
    hasProperty(propertyName) {
        return this.properties.has(propertyName);
    }
    notifyPropertyChanged(property) {
        this.adapter.getManager().sendPropertyChangedNotification(property);
    }
    actionNotify(action) {
        this.adapter.getManager().sendActionStatusNotification(action);
    }
    eventNotify(event) {
        this.adapter.getManager().sendEventNotification(event);
    }
    connectedNotify(connected) {
        this.adapter.getManager().sendConnectedNotification(this, connected);
    }
    setDescription(description) {
        this.description = description;
    }
    /**
     * @deprecated Please use setName()
     */
    setName(name) {
        console.log('setName() is deprecated. Please use setTitle().');
        this.setTitle(name);
    }
    setTitle(title) {
        this.title = title;
    }
    /**
     * @method setProperty
     * @returns a promise which resolves to the updated value.
     *
     * @note it is possible that the updated value doesn't match
     * the value passed in.
     */
    setProperty(propertyName, value) {
        const property = this.findProperty(propertyName);
        if (property) {
            return property.setValue(value);
        }
        return Promise.reject(`Property "${propertyName}" not found`);
    }
    getAdapter() {
        return this.adapter;
    }
    /**
     * @method requestAction
     * @returns a promise which resolves when the action has been requested.
     */
    requestAction(actionId, actionName, input) {
        return new Promise((resolve, reject) => {
            if (!this.actions.has(actionName)) {
                reject(`Action "${actionName}" not found`);
                return;
            }
            // Validate action input, if present.
            const metadata = this.actions.get(actionName);
            if (metadata) {
                if (metadata.hasOwnProperty('input')) {
                    // eslint-disable-next-line @typescript-eslint/no-explicit-any
                    const valid = ajv.validate(metadata.input, input);
                    if (!valid) {
                        reject(`Action "${actionName}": input "${input}" is invalid`);
                    }
                }
            }
            else {
                reject(`Action "${actionName}" not found`);
            }
            const action = new action_1.Action(actionId, this, actionName, input);
            this.performAction(action).catch((err) => console.log(err));
            resolve();
        });
    }
    /**
     * @method removeAction
     * @returns a promise which resolves when the action has been removed.
     */
    removeAction(actionId, actionName) {
        return new Promise((resolve, reject) => {
            if (!this.actions.has(actionName)) {
                reject(`Action "${actionName}" not found`);
                return;
            }
            this.cancelAction(actionId, actionName).catch((err) => console.log(err));
            resolve();
        });
    }
    /**
     * @method performAction
     */
    // eslint-disable-next-line @typescript-eslint/no-unused-vars
    performAction(_action) {
        return Promise.resolve();
    }
    /**
     * @method cancelAction
     */
    // eslint-disable-next-line @typescript-eslint/no-unused-vars
    cancelAction(_actionId, _actionName) {
        return Promise.resolve();
    }
    /**
     * Add an action.
     *
     * @param {String} name Name of the action
     * @param {Object} metadata Action metadata, i.e. type, description, etc., as
     *                          an object
     */
    addAction(name, metadata) {
        metadata = metadata !== null && metadata !== void 0 ? metadata : {};
        if (metadata.hasOwnProperty('href')) {
            const metadataWithHref = metadata;
            delete metadataWithHref.href;
        }
        this.actions.set(name, metadata);
    }
    /**
     * Add an event.
     *
     * @param {String} name Name of the event
     * @param {Object} metadata Event metadata, i.e. type, description, etc., as
     *                          an object
     */
    addEvent(name, metadata) {
        metadata = metadata !== null && metadata !== void 0 ? metadata : {};
        if (metadata.hasOwnProperty('href')) {
            const metadataWithHref = metadata;
            delete metadataWithHref.href;
        }
        this.events.set(name, metadata);
    }
}
exports.Device = Device;
//# sourceMappingURL=device.js.map