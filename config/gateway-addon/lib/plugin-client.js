"use strict";
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
Object.defineProperty(exports, "__esModule", { value: true });
exports.PluginClient = void 0;
const addon_manager_proxy_1 = require("./addon-manager-proxy");
const constants_1 = require("./constants");
const deferred_1 = require("./deferred");
const events_1 = require("events");
const ipc_1 = require("./ipc");
class PluginClient extends events_1.EventEmitter {
    constructor(pluginId, { verbose } = {}) {
        super();
        this.pluginId = pluginId;
        this.verbose = !!verbose;
        this.logPrefix = `PluginClient(${this.pluginId}):`;
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
    onMsg(genericMsg) {
        this.verbose &&
            console.log(this.logPrefix, 'rcvd ManagerMsg:', genericMsg);
        if (genericMsg.messageType === constants_1.MessageType.PLUGIN_REGISTER_RESPONSE) {
            const msg = genericMsg;
            this.gatewayVersion = msg.data.gatewayVersion;
            this.userProfile = msg.data.userProfile;
            this.preferences = msg.data.preferences;
            this.addonManager = new addon_manager_proxy_1.AddonManagerProxy(this);
            this.verbose &&
                console.log(this.logPrefix, 'registered with PluginServer');
            if (this.deferredReply) {
                const deferredReply = this.deferredReply;
                this.deferredReply = null;
                deferredReply.resolve(this.addonManager);
            }
        }
        else if (this.addonManager) {
            this.addonManager.onMsg(genericMsg);
        }
    }
    register(port) {
        var _a, _b;
        if (this.deferredReply) {
            console.error(this.logPrefix, 'Already waiting for registration reply');
            return Promise.resolve();
        }
        this.deferredReply = new deferred_1.Deferred();
        this.ipcSocket = new ipc_1.IpcSocket(false, port, this.onMsg.bind(this), `IpcSocket(${this.pluginId}):`, { verbose: this.verbose });
        (_b = (_a = this.ipcSocket) === null || _a === void 0 ? void 0 : _a.getConnectPromise()) === null || _b === void 0 ? void 0 : _b.then((ws) => {
            this.ws = ws;
            // Register ourselves with the server
            this.verbose &&
                console.log(this.logPrefix, 'Connected to server, registering...');
            this.sendNotification(constants_1.MessageType.PLUGIN_REGISTER_REQUEST);
        });
        return this.deferredReply.getPromise();
    }
    sendNotification(messageType, data = {}) {
        var _a;
        data.pluginId = this.pluginId;
        const jsonObj = JSON.stringify({ messageType, data });
        this.verbose && console.log(this.logPrefix, 'Sending:', jsonObj);
        (_a = this.ws) === null || _a === void 0 ? void 0 : _a.send(jsonObj);
    }
    unload() {
        var _a;
        (_a = this.ipcSocket) === null || _a === void 0 ? void 0 : _a.close();
        this.emit('unloaded', {});
    }
}
exports.PluginClient = PluginClient;
//# sourceMappingURL=plugin-client.js.map