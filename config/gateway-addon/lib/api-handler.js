"use strict";
/**
 * @module API Handler base class.
 *
 * Allows add-ons to create generic REST API handlers without having to create
 * a full HTTP server.
 */
/**
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/.*
 */
var __awaiter = (this && this.__awaiter) || function (thisArg, _arguments, P, generator) {
    function adopt(value) { return value instanceof P ? value : new P(function (resolve) { resolve(value); }); }
    return new (P || (P = Promise))(function (resolve, reject) {
        function fulfilled(value) { try { step(generator.next(value)); } catch (e) { reject(e); } }
        function rejected(value) { try { step(generator["throw"](value)); } catch (e) { reject(e); } }
        function step(result) { result.done ? resolve(result.value) : adopt(result.value).then(fulfilled, rejected); }
        step((generator = generator.apply(thisArg, _arguments || [])).next());
    });
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.APIHandler = exports.APIResponse = exports.APIRequest = void 0;
/**
 * Class which holds an API request.
 */
class APIRequest {
    /**
     * Build the request.
     *
     * @param {object} params - Request parameters, as such:
     *                   .method {string} HTTP method, e.g. GET, POST, etc.
     *                   .path {string} Path relative to this handler, e.g.
     *                     '/mypath' rather than
     *                     '/extensions/my-extension/api/mypath'.
     *                   .query {object} Object containing query parameters
     *                   .body {object} Body content in key/value form. All
     *                     content should be requested as application/json or
     *                     application/x-www-form-urlencoded data in order for it
     *                     to be parsed properly.
     */
    constructor({ method, path, query, body }) {
        this.method = method;
        this.path = path;
        this.query = query !== null && query !== void 0 ? query : {};
        this.body = body !== null && body !== void 0 ? body : {};
    }
    getMethod() {
        return this.method;
    }
    getPath() {
        return this.path;
    }
    getQuery() {
        return this.query;
    }
    getBody() {
        return this.body;
    }
}
exports.APIRequest = APIRequest;
/**
 * Convenience class to build an API response.
 */
class APIResponse {
    /**
     * Build the response.
     *
     * @param {object} params - Response parameters, as such:
     *                   .status {number} (Required) Status code
     *                   .contentType {string} Content-Type of response content
     *                   .content {string} Response content
     */
    constructor({ status, contentType, content } = { status: 500 }) {
        this.status = Number(status);
        if (contentType) {
            this.contentType = `${contentType}`;
        }
        if (content) {
            this.content = `${content}`;
        }
    }
    getStatus() {
        return this.status;
    }
    getContentType() {
        return this.contentType;
    }
    getContent() {
        return this.content;
    }
}
exports.APIResponse = APIResponse;
/**
 * Base class for API handlers, which handle sending alerts to a user.
 * @class Notifier
 */
class APIHandler {
    constructor(manager, packageName, { verbose } = {}) {
        this.packageName = packageName;
        this.verbose = !!verbose;
        this.gatewayVersion = manager.getGatewayVersion();
        this.userProfile = manager.getUserProfile();
        this.preferences = manager.getPreferences();
    }
    isVerbose() {
        return this.verbose;
    }
    getPackageName() {
        return this.packageName;
    }
    getGatewayVersion() {
        return this.gatewayVersion;
    }
    getUserProfile() {
        return this.userProfile;
    }
    getPreferences() {
        return this.preferences;
    }
    /**
     * @method handleRequest
     *
     * Called every time a new API request comes in for this handler.
     *
     * @param {APIRequest} request - Request object
     *
     * @returns {APIResponse} API response object.
     */
    handleRequest(request) {
        return __awaiter(this, void 0, void 0, function* () {
            if (this.verbose) {
                console.log(`New API request for ${this.packageName}:`, request);
            }
            return new APIResponse({ status: 404 });
        });
    }
    /**
     * Unloads the handler.
     *
     * @returns a promise which resolves when the handler has finished unloading.
     */
    unload() {
        if (this.verbose) {
            console.log('API Handler', this.packageName, 'unloaded');
        }
        return Promise.resolve();
    }
}
exports.APIHandler = APIHandler;
//# sourceMappingURL=api-handler.js.map