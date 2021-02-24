/**
 * High-level Event base class implementation.
 *
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/.
 */
import { Device } from './device';
import { Any, EventDescription1 } from './schema';
/**
 * An Event represents an individual event from a device.
 */
export declare class Event {
    private device;
    private name;
    private data?;
    private timestamp;
    /**
     * Initialize the object.
     *
     * @param {Object} device Device this event belongs to
     * @param {String} name Name of the event
     * @param {*} data (Optional) Data associated with the event
     */
    constructor(device: Device, name: string, data?: Any);
    getDevice(): Device;
    /**
     * Get the event description.
     *
     * @returns {Object} Description of the event as an object.
     */
    asEventDescription(): EventDescription1;
    /**
     * Get the event description.
     *
     * @returns {Object} Description of the event as an object.
     */
    asDict(): EventDescription1;
}
//# sourceMappingURL=event.d.ts.map