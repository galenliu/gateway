import WebSocket from 'ws';
import { Message } from './schema';
export declare class IpcSocket {
    private isServer;
    private port;
    private onMsg;
    private logPrefix;
    private verbose;
    private validators;
    private wss?;
    private ws?;
    private connectPromise?;
    constructor(isServer: boolean, port: number, onMsg: (_data: Message, _ws: WebSocket) => void, logPrefix: string, { verbose }?: Record<string, unknown>);
    getConnectPromise(): Promise<WebSocket> | undefined;
    error(...args: unknown[]): void;
    log(...args: unknown[]): void;
    close(): void;
    /**
     * @method onData
     * @param {Buffer} buf
     *
     * Called anytime a new message has been received.
     */
    onData(buf: WebSocket.Data, ws: WebSocket): void;
}
//# sourceMappingURL=ipc.d.ts.map