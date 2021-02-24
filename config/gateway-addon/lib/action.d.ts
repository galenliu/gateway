/**
 * High-level Action base class implementation.
 *
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/.
 */
import { Device } from './device';
import { ActionDescription, Input } from './schema';
/**
 * An Action represents an individual action on a device.
 */
export declare class Action {
    private status;
    private timeRequested;
    private timeCompleted?;
    private id;
    device: Device;
    private name;
    private input?;
    /**
    * Initialize the object.
    *
    * @param {String} id ID of this action
    * @param {Object} device Device this action belongs to
    * @param {String} name Name of the action
    * @param {unknown} input Any action inputs
    */
    constructor(id: string, device: Device, name: string, input?: Input);
    /**
     * Get the action description.
     *
     * @returns {Object} Description of the action as an object.
     */
    asActionDescription(): ActionDescription;
    /**
     * Get the action description.
     *
     * @returns {Object} Description of the action as an object.
     */
    asDict(): ActionDescription;
    /**
     * Start performing the action.
     */
    start(): void;
    /**
     * Finish performing the action.
     */
    finish(): void;
}
//# sourceMappingURL=action.d.ts.map