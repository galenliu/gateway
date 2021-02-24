"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.IpcSocket = void 0;
const ajv_1 = __importDefault(require("ajv"));
const fs_1 = __importDefault(require("fs"));
const path_1 = __importDefault(require("path"));
const ws_1 = __importDefault(require("ws"));
class IpcSocket {
    constructor(isServer, port, onMsg, logPrefix, { verbose } = {}) {
        var _a;
        this.validators = {};
        this.isServer = isServer;
        this.port = port;
        this.onMsg = onMsg;
        this.logPrefix = logPrefix;
        this.verbose = !!verbose;
        // Build the JSON-Schema validator for incoming messages
        const baseDir = path_1.default.resolve(path_1.default.join(__dirname, '..', 'schema'));
        const schemas = [];
        // top-level schema
        schemas.push(JSON.parse(fs_1.default.readFileSync(path_1.default.join(baseDir, 'schema.json')).toString()));
        // individual message schemas
        for (const fname of fs_1.default.readdirSync(path_1.default.join(baseDir, 'messages'))) {
            const filePath = path_1.default.join(baseDir, 'messages', fname);
            schemas.push(JSON.parse(fs_1.default.readFileSync(filePath).toString()));
        }
        for (const schema of schemas) {
            if ((_a = schema === null || schema === void 0 ? void 0 : schema.properties) === null || _a === void 0 ? void 0 : _a.messageType) {
                const validate = new ajv_1.default({ schemas }).getSchema(schema.$id);
                if (validate) {
                    this.validators[schema.properties.messageType.const] = validate;
                }
            }
            else {
                console.debug(`Ignoring ${schema.$id} because it has no messageType`);
            }
        }
        if (this.isServer) {
            this.wss = new ws_1.default.Server({ host: '127.0.0.1', port: this.port });
            this.wss.on('connection', (ws) => {
                ws.on('message', (data) => {
                    this.onData(data, ws);
                });
            });
        }
        else {
            const ws = new ws_1.default(`ws://127.0.0.1:${this.port}/`);
            this.ws = ws;
            this.connectPromise = new Promise((resolve) => {
                ws.on('open', () => resolve(ws));
            });
            this.ws.on('message', this.onData.bind(this));
        }
    }
    getConnectPromise() {
        return this.connectPromise;
    }
    error(...args) {
        Array.prototype.unshift.call(args, this.logPrefix);
        console.error.apply(null, args);
    }
    log(...args) {
        Array.prototype.unshift.call(args, this.logPrefix);
        console.log.apply(null, args);
    }
    close() {
        var _a, _b;
        if (this.isServer) {
            (_a = this === null || this === void 0 ? void 0 : this.wss) === null || _a === void 0 ? void 0 : _a.close();
        }
        else {
            (_b = this === null || this === void 0 ? void 0 : this.ws) === null || _b === void 0 ? void 0 : _b.close();
        }
    }
    /**
     * @method onData
     * @param {Buffer} buf
     *
     * Called anytime a new message has been received.
     */
    onData(buf, ws) {
        const bufStr = buf.toString();
        let data;
        try {
            data = JSON.parse(bufStr);
        }
        catch (err) {
            this.error('Error parsing message as JSON');
            this.error(`Rcvd: "${bufStr}"`);
            this.error(err);
            return;
        }
        this.verbose && this.log('Rcvd:', data);
        // validate the message before forwarding to handler
        const messageType = data.messageType;
        if (typeof messageType !== 'undefined') {
            if (messageType in this.validators) {
                const validator = this.validators[messageType];
                if (!validator(data)) {
                    const dataJson = JSON.stringify(data, null, 2);
                    const errorJson = JSON.stringify(validator.errors, null, 2);
                    console.error(`Invalid message received: ${dataJson}`);
                    console.error(`Validation error: ${errorJson}`);
                }
            }
            else {
                console.error(`Unknown messageType ${messageType}`);
            }
        }
        else {
            console.error(`Message ${bufStr} has no messageType`);
        }
        this.onMsg(data, ws);
    }
}
exports.IpcSocket = IpcSocket;
//# sourceMappingURL=ipc.js.map