"use strict";
/**
 * @module Adapter base class.
 *
 * Manages Adapter data model and business logic.
 */
/**
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/.*
 */
Object.defineProperty(exports, "__esModule", { value: true });
exports.Adapter = void 0;
/**
 * Base class for adapters, which manage devices.
 * @class Adapter
 *
 */
class Adapter {
    constructor(manager, id, packageName, { verbose } = {}) {
        this.name = this.constructor.name;
        this.devices = {};
        this.actions = {};
        this.manager = manager;
        this.id = id;
        this.packageName = packageName;
        this.verbose = !!verbose;
        // We assume that the adapter is ready right away. If, for some reason
        // a particular adapter (like ZWave) needs some time, then it should
        // set ready to false in its constructor.
        this.ready = true;
        this.gatewayVersion = manager.getGatewayVersion();
        this.userProfile = manager.getUserProfile();
        this.preferences = manager.getPreferences();
    }
    dump() {
        if (this.verbose) {
            console.log('Adapter:', this.name, '- dump() not implemented');
        }
    }
    /**
     * @method getId
     * @returns the id of this adapter.
     */
    getId() {
        return this.id;
    }
    getPackageName() {
        return this.packageName;
    }
    getDevice(id) {
        return this.devices[id];
    }
    getDevices() {
        return this.devices;
    }
    getActions() {
        return this.actions;
    }
    getName() {
        return this.name;
    }
    isReady() {
        return this.ready;
    }
    getManager() {
        return this.manager;
    }
    isVerbose() {
        return this.verbose;
    }
    getGatewayVersion() {
        return this.gatewayVersion;
    }
    getUserProfile() {
        return this.userProfile;
    }
    getPreferences() {
        return this.preferences;
    }
    asDict() {
        return {
            id: this.getId(),
            name: this.getName(),
            ready: this.isReady(),
        };
    }
    /**
     * @method handleDeviceAdded
     *
     * Called to indicate that a device is now being managed by this adapter.
     */
    handleDeviceAdded(device) {
        this.devices[device.getId()] = device;
        this.manager.handleDeviceAdded(device);
    }
    /**
     * @method handleDeviceRemoved
     *
     * Called to indicate that a device is no longer managed by this adapter.
     */
    handleDeviceRemoved(device) {
        delete this.devices[device.getId()];
        this.manager.handleDeviceRemoved(device);
    }
    /**
     * @method handleDeviceSaved
     *
     * Called to indicate that the user has saved a device to their gateway. This
     * is also called when the adapter starts up for every device which has
     * already been saved.
     *
     * This can be used for keeping track of what devices have previously been
     * discovered, such that the adapter can rebuild those, clean up old nodes,
     * etc.
     *
     * @param {string} deviceId - ID of the device
     * @param {object} device - the saved device description
     */
    // eslint-disable-next-line @typescript-eslint/no-unused-vars
    handleDeviceSaved(_deviceId, _device) {
        // eslint-disable-next-line @typescript-eslint/no-empty-function
    }
    // eslint-disable-next-line @typescript-eslint/no-unused-vars
    startPairing(_timeoutSeconds) {
        if (this.verbose) {
            console.log('Adapter:', this.name, 'id', this.id, 'pairing started');
        }
    }
    /**
     * Send a prompt to the UI notifying the user to take some action.
     *
     * @param {string} prompt - The prompt to send
     * @param {string} url - URL to site with further explanation or
     *                 troubleshooting info
     * @param {Object?} device - Device the prompt is associated with
     */
    sendPairingPrompt(prompt, url, device) {
        this.manager.sendPairingPrompt(this, prompt, url, device);
    }
    /**
     * Send a prompt to the UI notifying the user to take some action.
     *
     * @param {string} prompt - The prompt to send
     * @param {string} url - URL to site with further explanation or
     *                 troubleshooting info
     * @param {Object?} device - Device the prompt is associated with
     */
    sendUnpairingPrompt(prompt, url, device) {
        this.manager.sendUnpairingPrompt(this, prompt, url, device);
    }
    cancelPairing() {
        if (this.verbose) {
            console.log('Adapter:', this.name, 'id', this.id, 'pairing cancelled');
        }
    }
    removeThing(device) {
        if (this.verbose) {
            console.log('Adapter:', this.name, 'id', this.id, 'removeThing(', device.getId(), ') started');
        }
        this.handleDeviceRemoved(device);
    }
    cancelRemoveThing(device) {
        if (this.verbose) {
            console.log('Adapter:', this.name, 'id', this.id, 'cancelRemoveThing(', device.getId(), ')');
        }
    }
    /**
     * Unloads an adapter.
     *
     * @returns a promise which resolves when the adapter has finished unloading.
     */
    unload() {
        if (this.verbose) {
            console.log('Adapter:', this.name, 'unloaded');
        }
        return Promise.resolve();
    }
    /**
     * Set the PIN for the given device.
     *
     * @param {String} deviceId ID of device
     * @param {String} pin PIN to set
     *
     * @returns a promise which resolves when the PIN has been set.
     */
    setPin(deviceId, pin) {
        const device = this.getDevice(deviceId);
        if (device) {
            if (this.verbose) {
                console.log('Adapter:', this.name, 'id', this.id, 'setPin(', deviceId, ',', pin, ')');
            }
            return Promise.resolve();
        }
        return Promise.reject('Device not found');
    }
    /**
     * Set the username and password for the given device.
     *
     * @param {String} deviceId ID of device
     * @param {String} username Username to set
     * @param {String} password Password to set
     *
     * @returns a promise which resolves when the credentials have been set.
     */
    setCredentials(deviceId, username, password) {
        const device = this.getDevice(deviceId);
        if (device) {
            if (this.verbose) {
                console.log('Adapter:', this.name, 'id', this.id, 'setCredentials(', deviceId, ',', username, ',', password, ')');
            }
            return Promise.resolve();
        }
        return Promise.reject('Device not found');
    }
}
exports.Adapter = Adapter;
//# sourceMappingURL=adapter.js.map