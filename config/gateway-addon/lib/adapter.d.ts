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
import { Action } from './action';
import { AddonManagerProxy } from './addon-manager-proxy';
import { Device } from './device';
import { DeviceWithoutId as DeviceWithoutIdSchema, Preferences, UserProfile } from './schema';
export interface AdapterDescription {
    id: string;
    name: string;
    ready: boolean;
}
/**
 * Base class for adapters, which manage devices.
 * @class Adapter
 *
 */
export declare class Adapter {
    private manager;
    private id;
    private packageName;
    private verbose;
    private name;
    private devices;
    private actions;
    private ready;
    private gatewayVersion?;
    private userProfile?;
    private preferences?;
    constructor(manager: AddonManagerProxy, id: string, packageName: string, { verbose }?: Record<string, unknown>);
    dump(): void;
    /**
     * @method getId
     * @returns the id of this adapter.
     */
    getId(): string;
    getPackageName(): string;
    getDevice(id: string): Device;
    getDevices(): Record<string, Device>;
    getActions(): Record<string, Action>;
    getName(): string;
    isReady(): boolean;
    getManager(): AddonManagerProxy;
    isVerbose(): boolean;
    getGatewayVersion(): string | undefined;
    getUserProfile(): UserProfile | undefined;
    getPreferences(): Preferences | undefined;
    asDict(): AdapterDescription;
    /**
     * @method handleDeviceAdded
     *
     * Called to indicate that a device is now being managed by this adapter.
     */
    handleDeviceAdded(device: Device): void;
    /**
     * @method handleDeviceRemoved
     *
     * Called to indicate that a device is no longer managed by this adapter.
     */
    handleDeviceRemoved(device: Device): void;
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
    handleDeviceSaved(_deviceId: string, _device: DeviceWithoutIdSchema): void;
    startPairing(_timeoutSeconds: number): void;
    /**
     * Send a prompt to the UI notifying the user to take some action.
     *
     * @param {string} prompt - The prompt to send
     * @param {string} url - URL to site with further explanation or
     *                 troubleshooting info
     * @param {Object?} device - Device the prompt is associated with
     */
    sendPairingPrompt(prompt: string, url?: string, device?: Device): void;
    /**
     * Send a prompt to the UI notifying the user to take some action.
     *
     * @param {string} prompt - The prompt to send
     * @param {string} url - URL to site with further explanation or
     *                 troubleshooting info
     * @param {Object?} device - Device the prompt is associated with
     */
    sendUnpairingPrompt(prompt: string, url?: string, device?: Device): void;
    cancelPairing(): void;
    removeThing(device: Device): void;
    cancelRemoveThing(device: Device): void;
    /**
     * Unloads an adapter.
     *
     * @returns a promise which resolves when the adapter has finished unloading.
     */
    unload(): Promise<void>;
    /**
     * Set the PIN for the given device.
     *
     * @param {String} deviceId ID of device
     * @param {String} pin PIN to set
     *
     * @returns a promise which resolves when the PIN has been set.
     */
    setPin(deviceId: string, pin: string): Promise<void>;
    /**
     * Set the username and password for the given device.
     *
     * @param {String} deviceId ID of device
     * @param {String} username Username to set
     * @param {String} password Password to set
     *
     * @returns a promise which resolves when the credentials have been set.
     */
    setCredentials(deviceId: string, username: string, password: string): Promise<void>;
}
//# sourceMappingURL=adapter.d.ts.map