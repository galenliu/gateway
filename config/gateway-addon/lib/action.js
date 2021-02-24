"use strict";
/**
 * High-level Action base class implementation.
 *
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/.
 */
Object.defineProperty(exports, "__esModule", { value: true });
exports.Action = void 0;
const utils_1 = require("./utils");
/**
 * An Action represents an individual action on a device.
 */
class Action {
    /**
    * Initialize the object.
    *
    * @param {String} id ID of this action
    * @param {Object} device Device this action belongs to
    * @param {String} name Name of the action
    * @param {unknown} input Any action inputs
    */
    constructor(id, device, name, input) {
        this.status = 'created';
        this.timeRequested = utils_1.timestamp();
        this.id = id;
        this.device = device;
        this.name = name;
        this.input = input;
    }
    /**
     * Get the action description.
     *
     * @returns {Object} Description of the action as an object.
     */
    asActionDescription() {
        const description = {
            id: this.id,
            name: this.name,
            timeRequested: this.timeRequested,
            status: this.status,
        };
        if (this.input !== null) {
            description.input = this.input;
        }
        if (this.timeCompleted !== null) {
            description.timeCompleted = this.timeCompleted;
        }
        return description;
    }
    /**
     * Get the action description.
     *
     * @returns {Object} Description of the action as an object.
     */
    asDict() {
        return {
            id: this.id,
            name: this.name,
            input: this.input,
            status: this.status,
            timeRequested: this.timeRequested,
            timeCompleted: this.timeCompleted,
        };
    }
    /**
     * Start performing the action.
     */
    start() {
        this.status = 'pending';
        this.device.actionNotify(this);
    }
    /**
     * Finish performing the action.
     */
    finish() {
        this.status = 'completed';
        this.timeCompleted = utils_1.timestamp();
        this.device.actionNotify(this);
    }
}
exports.Action = Action;
//# sourceMappingURL=action.js.map