"use strict";
/**
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/.*
 */
var __createBinding = (this && this.__createBinding) || (Object.create ? (function(o, m, k, k2) {
    if (k2 === undefined) k2 = k;
    Object.defineProperty(o, k2, { enumerable: true, get: function() { return m[k]; } });
}) : (function(o, m, k, k2) {
    if (k2 === undefined) k2 = k;
    o[k2] = m[k];
}));
var __setModuleDefault = (this && this.__setModuleDefault) || (Object.create ? (function(o, v) {
    Object.defineProperty(o, "default", { enumerable: true, value: v });
}) : function(o, v) {
    o["default"] = v;
});
var __importStar = (this && this.__importStar) || function (mod) {
    if (mod && mod.__esModule) return mod;
    var result = {};
    if (mod != null) for (var k in mod) if (k !== "default" && Object.prototype.hasOwnProperty.call(mod, k)) __createBinding(result, mod, k);
    __setModuleDefault(result, mod);
    return result;
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.getVersion = exports.Utils = exports.Property = exports.PluginClient = exports.Outlet = exports.Notifier = exports.IpcSocket = exports.Event = exports.Device = exports.Deferred = exports.Database = exports.Constants = exports.APIResponse = exports.APIRequest = exports.APIHandler = exports.AddonManagerProxy = exports.Adapter = exports.Action = void 0;
const action_1 = require("./action");
Object.defineProperty(exports, "Action", { enumerable: true, get: function () { return action_1.Action; } });
const adapter_1 = require("./adapter");
Object.defineProperty(exports, "Adapter", { enumerable: true, get: function () { return adapter_1.Adapter; } });
const addon_manager_proxy_1 = require("./addon-manager-proxy");
Object.defineProperty(exports, "AddonManagerProxy", { enumerable: true, get: function () { return addon_manager_proxy_1.AddonManagerProxy; } });
const api_handler_1 = require("./api-handler");
Object.defineProperty(exports, "APIHandler", { enumerable: true, get: function () { return api_handler_1.APIHandler; } });
Object.defineProperty(exports, "APIRequest", { enumerable: true, get: function () { return api_handler_1.APIRequest; } });
Object.defineProperty(exports, "APIResponse", { enumerable: true, get: function () { return api_handler_1.APIResponse; } });
const Constants = __importStar(require("./constants"));
exports.Constants = Constants;
const database_1 = require("./database");
Object.defineProperty(exports, "Database", { enumerable: true, get: function () { return database_1.Database; } });
const deferred_1 = require("./deferred");
Object.defineProperty(exports, "Deferred", { enumerable: true, get: function () { return deferred_1.Deferred; } });
const device_1 = require("./device");
Object.defineProperty(exports, "Device", { enumerable: true, get: function () { return device_1.Device; } });
const event_1 = require("./event");
Object.defineProperty(exports, "Event", { enumerable: true, get: function () { return event_1.Event; } });
const ipc_1 = require("./ipc");
Object.defineProperty(exports, "IpcSocket", { enumerable: true, get: function () { return ipc_1.IpcSocket; } });
const notifier_1 = require("./notifier");
Object.defineProperty(exports, "Notifier", { enumerable: true, get: function () { return notifier_1.Notifier; } });
const outlet_1 = require("./outlet");
Object.defineProperty(exports, "Outlet", { enumerable: true, get: function () { return outlet_1.Outlet; } });
const plugin_client_1 = require("./plugin-client");
Object.defineProperty(exports, "PluginClient", { enumerable: true, get: function () { return plugin_client_1.PluginClient; } });
const property_1 = require("./property");
Object.defineProperty(exports, "Property", { enumerable: true, get: function () { return property_1.Property; } });
const Utils = __importStar(require("./utils"));
exports.Utils = Utils;
function getVersion() {
    // eslint-disable-next-line @typescript-eslint/no-var-requires
    return require('./package.json').version;
}
exports.getVersion = getVersion;
//# sourceMappingURL=index.js.map