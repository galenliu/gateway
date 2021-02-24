/**
 * Proxy version of AddonManager used by plugins.
 *
 * @module AddonManagerProxy
 */
/**
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/.
 */
/// <reference types="node" />
import { Action } from './action';
import { Adapter } from './adapter';
import { APIHandler } from './api-handler';
import { Device } from './device';
import { Event } from './event';
import { Notifier } from './notifier';
import { Outlet } from './outlet';
import { PluginClient } from './plugin-client';
import { Property } from './property';
import { EventEmitter } from 'events';
import { Message, Preferences, PropertyValue, UserProfile } from './schema';
export declare class AddonManagerProxy extends EventEmitter {
    private pluginClient;
    private gatewayVersion?;
    private userProfile?;
    private preferences?;
    private verbose;
    private adapters;
    private notifiers;
    private apiHandlers;
    constructor(pluginClient: PluginClient, { verbose }?: Record<string, unknown>);
    getGatewayVersion(): string | undefined;
    getUserProfile(): UserProfile | undefined;
    getPreferences(): Preferences | undefined;
    /**
     * @method addAdapter
     *
     * Adds an adapter to the collection of adapters managed by AddonManager.
     */
    addAdapter(adapter: Adapter): void;
    /**
     * @method addNotifier
     *
     * Adds a notifier to the collection of notifiers managed by AddonManager.
     */
    addNotifier(notifier: Notifier): void;
    /**
     * @method addAPIHandler
     *
     * Adds a new API handler.
     */
    addAPIHandler(handler: APIHandler): void;
    /**
     * @method handleDeviceAdded
     *
     * Called when the indicated device has been added to an adapter.
     */
    handleDeviceAdded(device: Device): void;
    /**
     * @method handleDeviceRemoved
     * Called when the indicated device has been removed from an adapter.
     */
    handleDeviceRemoved(device: Device): void;
    /**
     * @method handleOutletAdded
     *
     * Called when the indicated outlet has been added to a notifier.
     */
    handleOutletAdded(outlet: Outlet): void;
    /**
     * @method handleOutletRemoved
     * Called when the indicated outlet has been removed from a notifier.
     */
    handleOutletRemoved(outlet: Outlet): void;
    /**
     * @method onMsg
     * Called whenever a message is received from the gateway.
     */
    onMsg(genericMsg: Message): void;
    /**
     * @method sendPairingPrompt
     * Send a prompt to the UI notifying the user to take some action.
     */
    sendPairingPrompt(adapter: Adapter, prompt: string, url?: string, device?: Device): void;
    /**
     * @method sendUnpairingPrompt
     * Send a prompt to the UI notifying the user to take some action.
     */
    sendUnpairingPrompt(adapter: Adapter, prompt: string, url?: string, device?: Device): void;
    /**
     * @method sendPropertyChangedNotification
     * Sends a propertyChanged notification to the gateway.
     */
    sendPropertyChangedNotification(property: Property<PropertyValue>): void;
    /**
     * @method sendActionStatusNotification
     * Sends an actionStatus notification to the gateway.
     */
    sendActionStatusNotification(action: Action): void;
    /**
     * @method sendEventNotification
     * Sends an event notification to the gateway.
     */
    sendEventNotification(event: Event): void;
    /**
     * @method sendConnectedNotification
     * Sends a connected notification to the gateway.
     */
    sendConnectedNotification(device: Device, connected: boolean): void;
    /**
     * @method unloadPlugin
     *
     * Unloads the plugin, and tells the server about it.
     */
    unloadPlugin(): void;
    sendError(message: string): void;
}
//# sourceMappingURL=addon-manager-proxy.d.ts.map