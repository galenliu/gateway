/**
 * Wrapper around the gateway's database.
 *
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/.
 */
/**
 * An Action represents an individual action on a device.
 */
export declare class Database {
    private packageName;
    private path;
    private conn?;
    /**
     * Initialize the object.
     *
     * @param {String} packageName The adapter's package name
     * @param {String?} path Optional database path
     */
    constructor(packageName: string, path: string);
    /**
     * Open the database.
     *
     * @returns Promise which resolves when the database has been opened.
     */
    open(): Promise<void>;
    /**
     * Close the database.
     */
    close(): void;
    /**
     * Load the package's config from the database.
     *
     * @returns Promise which resolves to the config object.
     */
    loadConfig(): Promise<Record<string, unknown>>;
    /**
     * Save the package's config to the database.
     */
    saveConfig(config: Record<string, unknown>): Promise<void>;
}
//# sourceMappingURL=database.d.ts.map