"use strict";
/**
 * Wraps up a promise in a slightly more convenient manner for passing
 * around, or saving.
 *
 * @module Deferred
 */
Object.defineProperty(exports, "__esModule", { value: true });
exports.Deferred = void 0;
let id = 0;
class Deferred {
    constructor() {
        this.id = ++id;
        this.promise = new Promise((resolve, reject) => {
            this.resolveFunc = resolve;
            this.rejectFunc = reject;
        });
    }
    resolve(arg) {
        if (this.resolveFunc) {
            return this.resolveFunc(arg);
        }
    }
    reject(arg) {
        if (this.rejectFunc) {
            return this.rejectFunc(arg);
        }
    }
    getId() {
        return this.id;
    }
    getPromise() {
        return this.promise;
    }
}
exports.Deferred = Deferred;
//# sourceMappingURL=deferred.js.map