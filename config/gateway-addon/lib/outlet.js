"use strict";
/**
 * Outlet Model.
 *
 * Abstract base class for outlets managed by a notifier.
 *
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/.
 */
Object.defineProperty(exports, "__esModule", { value: true });
exports.Outlet = void 0;
class Outlet {
    constructor(notifier, id) {
        this.name = '';
        this.notifier = notifier;
        this.id = `${id}`;
    }
    asDict() {
        return {
            id: this.id,
            name: this.name,
        };
    }
    getId() {
        return this.id;
    }
    getName() {
        return this.name;
    }
    setName(name) {
        this.name = name;
    }
    getNotifier() {
        return this.notifier;
    }
    /**
     * Notify the user.
     *
     * @param {string} title Title of notification.
     * @param {string} message Message of notification.
     * @param {number} level Alert level.
     * @returns {Promise} Promise which resolves when the user has been notified.
     */
    notify(title, message, level) {
        if (this.notifier.isVerbose()) {
            console.log(`Outlet: ${this.name} notify("${title}", "${message}", ${level})`);
        }
        return Promise.resolve();
    }
}
exports.Outlet = Outlet;
//# sourceMappingURL=outlet.js.map