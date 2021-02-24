/**
 * Device Model.
 *
 * Abstract base class for devices managed by an adapter.
 *
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/.
 */
import { Action } from './action';
import { Adapter } from './adapter';
import { Property } from './property';
import { Event } from './event';
import { Action as ActionSchema, Event as EventSchema, Device as DeviceSchema, Input, PropertyValue } from './schema';
export declare class Device {
    private adapter;
    private id;
    private '@context';
    private '@type';
    private name;
    private title;
    private description;
    private properties;
    private actions;
    private events;
    private links;
    private baseHref?;
    private pinRequired;
    private pinPattern?;
    private credentialsRequired;
    constructor(adapter: Adapter, id: string);
    mapToDict<V>(map: Map<string, V>): Record<string, V>;
    mapToDictFromFunction<V>(map: Map<string, {
        asDict: () => V;
    }>): Record<string, V>;
    asDict(): DeviceSchema;
    /**
     * @returns this object as a thing
     */
    asThing(): DeviceSchema;
    debugCmd(cmd: string, params: unknown): void;
    getId(): string;
    /**
     * @deprecated Please use getTitle()
     */
    getName(): string;
    getTitle(): string;
    getPropertyDescriptions(): Record<string, unknown>;
    findProperty(propertyName: string): Property<PropertyValue> | undefined;
    addProperty(property: Property<PropertyValue>): void;
    /**
     * @method getProperty
     * @returns a promise which resolves to the retrieved value.
     */
    getProperty(propertyName: string): Promise<unknown>;
    hasProperty(propertyName: string): boolean;
    notifyPropertyChanged(property: Property<PropertyValue>): void;
    actionNotify(action: Action): void;
    eventNotify(event: Event): void;
    connectedNotify(connected: boolean): void;
    setDescription(description: string): void;
    /**
     * @deprecated Please use setName()
     */
    setName(name: string): void;
    setTitle(title: string): void;
    /**
     * @method setProperty
     * @returns a promise which resolves to the updated value.
     *
     * @note it is possible that the updated value doesn't match
     * the value passed in.
     */
    setProperty(propertyName: string, value: PropertyValue): Promise<PropertyValue>;
    getAdapter(): Adapter;
    /**
     * @method requestAction
     * @returns a promise which resolves when the action has been requested.
     */
    requestAction(actionId: string, actionName: string, input: Input): Promise<void>;
    /**
     * @method removeAction
     * @returns a promise which resolves when the action has been removed.
     */
    removeAction(actionId: string, actionName: string): Promise<void>;
    /**
     * @method performAction
     */
    performAction(_action: Action): Promise<void>;
    /**
     * @method cancelAction
     */
    cancelAction(_actionId: string, _actionName: string): Promise<void>;
    /**
     * Add an action.
     *
     * @param {String} name Name of the action
     * @param {Object} metadata Action metadata, i.e. type, description, etc., as
     *                          an object
     */
    addAction(name: string, metadata?: ActionSchema): void;
    /**
     * Add an event.
     *
     * @param {String} name Name of the event
     * @param {Object} metadata Event metadata, i.e. type, description, etc., as
     *                          an object
     */
    addEvent(name: string, metadata?: EventSchema): void;
}
//# sourceMappingURL=device.d.ts.map