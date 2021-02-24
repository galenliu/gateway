/**
 * Property.
 *
 * Object which decscribes a property, and its value.
 *
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/.
 */
import { Device } from './device';
import { Link, Property as PropertySchema, PropertyValue, PropertyValuesEnum, PropertyValueType } from './schema';
export declare class Property<T extends PropertyValue> {
    private device;
    private name;
    private title?;
    private type;
    private '@type'?;
    private unit?;
    private description?;
    private minimum?;
    private maximum?;
    private enum?;
    private readOnly?;
    private multipleOf?;
    private links;
    private visible;
    private fireAndForget;
    private value?;
    private prevGetValue?;
    constructor(device: Device, name: string, propertyDescr: PropertySchema);
    /**
     * @returns a dictionary of useful information.
     * This is primarily used for debugging.
     */
    asDict(): PropertySchema;
    /**
     * @returns the dictionary as used to describe a property. Currently
     * this does not include the href field.
     */
    asPropertyDescription(): PropertySchema;
    /**
     * @method isVisible
     * @returns true if this is a visible property, which is a property
     *          that is reported in the property description.
     */
    isVisible(): boolean;
    /**
     * Make the property visible or invisible
     */
    setVisible(visible: boolean): void;
    isFireAndForget(): boolean;
    /**
     * Sets the value and notifies the device if the value has changed.
     * @returns true if the value has changed
     */
    setCachedValueAndNotify(value: T): boolean;
    /**
     * Sets this.value and makes adjustments to ensure that the value
     * is consistent with the type.
     */
    setCachedValue(value: T): T;
    /**
     * @method getValue
     * @returns a promise which resolves to the retrieved value.
     *
     * This implementation is a simple one that just returns
     * the previously cached value.
     */
    getValue(): Promise<T>;
    /**
     * @method setValue
     * @returns a promise which resolves to the updated value.
     *
     * @note it is possible that the updated value doesn't match
     * the value passed in.
     *
     * It is anticipated that this method will most likely be overridden
     * by a derived class.
     */
    setValue(value: T): Promise<T>;
    getDevice(): Device;
    getName(): string;
    setName(value: string): void;
    getTitle(): string | undefined;
    setTitle(value: string): void;
    getType(): string | undefined;
    setType(value: PropertyValueType): void;
    getAtType(): string | undefined;
    setAtType(value: string): void;
    getUnit(): string | undefined;
    setUnit(value: string): void;
    getDescription(): string | undefined;
    setDescription(value: string): void;
    getMinimum(): number | undefined;
    setMinimum(value: number): void;
    getMaximum(): number | undefined;
    setMaximum(value: number): void;
    getEnum(): PropertyValuesEnum[] | undefined;
    setEnum(value: string[]): void;
    getReadOnly(): boolean | undefined;
    setReadOnly(value: boolean): void;
    getMultipleOf(): number | undefined;
    setMultipleOf(value: number): void;
    getLinks(): Link[];
    setLinks(value: Link[]): void;
}
//# sourceMappingURL=property.d.ts.map