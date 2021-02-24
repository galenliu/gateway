/**
 * @module Notifier base class.
 *
 * Manages Notifier data model and business logic.
 */
/**
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/.*
 */
import { AddonManagerProxy } from './addon-manager-proxy';
import { Outlet } from './outlet';
import { Preferences, UserProfile } from './schema';
export interface NotifierDescription {
    id: string;
    name: string;
    ready: boolean;
}
/**
 * Base class for notifiers, which handle sending alerts to a user.
 * @class Notifier
 */
export declare class Notifier {
    private manager;
    private id;
    private packageName;
    private verbose;
    private name;
    private outlets;
    private ready;
    private gatewayVersion?;
    private userProfile?;
    private preferences?;
    constructor(manager: AddonManagerProxy, id: string, packageName: string, { verbose }?: Record<string, unknown>);
    dump(): void;
    /**
     * @method getId
     * @returns the id of this adapter.
     */
    getId(): string;
    getPackageName(): string;
    getOutlet(id: string): Outlet;
    getOutlets(): Record<string, Outlet>;
    getName(): string;
    setName(name: string): void;
    isReady(): boolean;
    setReady(ready: boolean): void;
    isVerbose(): boolean;
    getGatewayVersion(): string | undefined;
    getUserProfile(): UserProfile | undefined;
    getPreferences(): Preferences | undefined;
    asDict(): NotifierDescription;
    /**
     * @method handleOutletAdded
     *
     * Called to indicate that an outlet is now being managed by this notifier.
     */
    handleOutletAdded(outlet: Outlet): void;
    /**
     * @method handleOutletRemoved
     *
     * Called to indicate that an outlet is no longer managed by this notifier.
     */
    handleOutletRemoved(outlet: Outlet): void;
    /**
     * Unloads a notifier.
     *
     * @returns a promise which resolves when the notifier has finished unloading.
     */
    unload(): Promise<void>;
}
//# sourceMappingURL=notifier.d.ts.map