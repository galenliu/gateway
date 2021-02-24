"use strict";
/**
 * Wrapper around the gateway's database.
 *
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/.
 */
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.Database = void 0;
const fs_1 = __importDefault(require("fs"));
const os_1 = __importDefault(require("os"));
const path_1 = __importDefault(require("path"));
const sqlite3_1 = require("sqlite3");
const sqlite3 = sqlite3_1.verbose();
const DB_PATHS = [
    path_1.default.join(os_1.default.homedir(), '.webthings', 'config', 'db.sqlite3'),
];
if (process.env.WEBTHINGS_HOME) {
    // eslint-disable-next-line max-len
    DB_PATHS.unshift(path_1.default.join(process.env.WEBTHINGS_HOME, 'config', 'db.sqlite3'));
}
if (process.env.WEBTHINGS_DATABASE) {
    DB_PATHS.unshift(process.env.WEBTHINGS_DATABASE);
}
/**
 * An Action represents an individual action on a device.
 */
class Database {
    /**
     * Initialize the object.
     *
     * @param {String} packageName The adapter's package name
     * @param {String?} path Optional database path
     */
    constructor(packageName, path) {
        this.packageName = packageName;
        this.path = path;
        if (!this.path) {
            for (const p of DB_PATHS) {
                if (fs_1.default.existsSync(p)) {
                    this.path = p;
                    break;
                }
            }
        }
    }
    /**
     * Open the database.
     *
     * @returns Promise which resolves when the database has been opened.
     */
    open() {
        if (this.conn) {
            return Promise.resolve();
        }
        if (!this.path) {
            return Promise.reject(new Error('Database path unknown'));
        }
        return new Promise((resolve, reject) => {
            this.conn = new sqlite3.Database(this.path, (err) => {
                var _a;
                if (err) {
                    reject(err);
                }
                else {
                    (_a = this === null || this === void 0 ? void 0 : this.conn) === null || _a === void 0 ? void 0 : _a.configure('busyTimeout', 10000);
                    resolve();
                }
            });
        });
    }
    /**
     * Close the database.
     */
    close() {
        if (this.conn) {
            this.conn.close();
            this.conn = null;
        }
    }
    /**
     * Load the package's config from the database.
     *
     * @returns Promise which resolves to the config object.
     */
    loadConfig() {
        if (!this.conn) {
            return Promise.reject('Database not open');
        }
        const key = `addons.config.${this.packageName}`;
        return new Promise((resolve, reject) => {
            var _a;
            (_a = this === null || this === void 0 ? void 0 : this.conn) === null || _a === void 0 ? void 0 : _a.get('SELECT value FROM settings WHERE key = ?', [key], (error, row) => {
                if (error) {
                    reject(error);
                }
                else if (!row) {
                    resolve({});
                }
                else {
                    resolve(JSON.parse(row.value));
                }
            });
        });
    }
    /**
     * Save the package's config to the database.
     */
    saveConfig(config) {
        if (!this.conn) {
            return Promise.resolve();
        }
        const key = `addons.config.${this.packageName}`;
        return new Promise((resolve, reject) => {
            var _a;
            (_a = this === null || this === void 0 ? void 0 : this.conn) === null || _a === void 0 ? void 0 : _a.run('INSERT OR REPLACE INTO settings (key, value) VALUES (?, ?)', [key, JSON.stringify(config)], (error) => {
                if (error) {
                    reject(error);
                }
                else {
                    resolve();
                }
            });
        });
    }
}
exports.Database = Database;
//# sourceMappingURL=database.js.map