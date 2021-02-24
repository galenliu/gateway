"use strict";
/*
 * WebThings Gateway Constants.
 *
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/.
 */
Object.defineProperty(exports, "__esModule", { value: true });
exports.DONT_RESTART_EXIT_CODE = exports.NotificationLevel = exports.MessageType = void 0;
const message_type_1 = require("./message-type");
Object.defineProperty(exports, "MessageType", { enumerable: true, get: function () { return message_type_1.MessageType; } });
var NotificationLevel;
(function (NotificationLevel) {
    NotificationLevel[NotificationLevel["LOW"] = 0] = "LOW";
    NotificationLevel[NotificationLevel["NORMAL"] = 1] = "NORMAL";
    NotificationLevel[NotificationLevel["HIGH"] = 2] = "HIGH";
})(NotificationLevel = exports.NotificationLevel || (exports.NotificationLevel = {}));
exports.DONT_RESTART_EXIT_CODE = 100;
//# sourceMappingURL=constants.js.map