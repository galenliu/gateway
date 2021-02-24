/**
 * Outlet Model.
 *
 * Abstract base class for outlets managed by a notifier.
 *
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/.
 */
import { Notifier } from './notifier';
import { Level, OutletDescription } from './schema';
export declare class Outlet {
    private notifier;
    private id;
    private name;
    constructor(notifier: Notifier, id: string);
    asDict(): OutletDescription;
    getId(): string;
    getName(): string;
    setName(name: string): void;
    getNotifier(): Notifier;
    /**
     * Notify the user.
     *
     * @param {string} title Title of notification.
     * @param {string} message Message of notification.
     * @param {number} level Alert level.
     * @returns {Promise} Promise which resolves when the user has been notified.
     */
    notify(title: string, message: string, level: Level): Promise<void>;
}
//# sourceMappingURL=outlet.d.ts.map