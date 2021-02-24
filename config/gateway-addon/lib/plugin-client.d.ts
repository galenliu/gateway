/**
 * @module PluginClient
 *
 * Takes care of connecting to the gateway for an adapter plugin
 */
/**
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/.
 */
/// <reference types="node" />
import { AddonManagerProxy } from './addon-manager-proxy';
import { EventEmitter } from 'events';
import { Message, Preferences, UserProfile } from './schema';
export declare class PluginClient extends EventEmitter {
    private pluginId;
    private verbose;
    private deferredReply?;
    private logPrefix;
    private gatewayVersion?;
    private userProfile?;
    private preferences?;
    private addonManager?;
    private ipcSocket?;
    private ws?;
    constructor(pluginId: string, { verbose }?: Record<string, unknown>);
    getGatewayVersion(): string | undefined;
    getUserProfile(): UserProfile | undefined;
    getPreferences(): Preferences | undefined;
    onMsg(genericMsg: Message): void;
    register(port: number): Promise<AddonManagerProxy | void>;
    sendNotification(messageType: number, data?: Record<string, unknown>): void;
    unload(): void;
}
//# sourceMappingURL=plugin-client.d.ts.map