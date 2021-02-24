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
import { AddonManagerProxy } from './addon-manager-proxy';
import { Preferences, Request, UserProfile } from './schema';
/**
 * Class which holds an API request.
 */
export declare class APIRequest {
    private method;
    private path;
    private query;
    private body;
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
    constructor({ method, path, query, body }: Request);
    getMethod(): string;
    getPath(): string;
    getQuery(): Record<string, unknown>;
    getBody(): Record<string, unknown>;
}
export interface APIResponseOptions {
    status: number;
    contentType?: string;
    content?: string;
}
/**
 * Convenience class to build an API response.
 */
export declare class APIResponse {
    private status;
    private contentType?;
    private content?;
    /**
     * Build the response.
     *
     * @param {object} params - Response parameters, as such:
     *                   .status {number} (Required) Status code
     *                   .contentType {string} Content-Type of response content
     *                   .content {string} Response content
     */
    constructor({ status, contentType, content }?: APIResponseOptions);
    getStatus(): number;
    getContentType(): string | undefined;
    getContent(): string | undefined;
}
/**
 * Base class for API handlers, which handle sending alerts to a user.
 * @class Notifier
 */
export declare class APIHandler {
    private packageName;
    private verbose;
    private gatewayVersion?;
    private userProfile?;
    private preferences?;
    constructor(manager: AddonManagerProxy, packageName: string, { verbose }?: Record<string, unknown>);
    isVerbose(): boolean;
    getPackageName(): string;
    getGatewayVersion(): string | undefined;
    getUserProfile(): UserProfile | undefined;
    getPreferences(): Preferences | undefined;
    /**
     * @method handleRequest
     *
     * Called every time a new API request comes in for this handler.
     *
     * @param {APIRequest} request - Request object
     *
     * @returns {APIResponse} API response object.
     */
    handleRequest(request: APIRequest): Promise<APIResponse>;
    /**
     * Unloads the handler.
     *
     * @returns a promise which resolves when the handler has finished unloading.
     */
    unload(): Promise<void>;
}
//# sourceMappingURL=api-handler.d.ts.map