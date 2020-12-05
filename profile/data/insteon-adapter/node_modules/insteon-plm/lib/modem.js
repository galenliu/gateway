import SerialPort from 'serialport';
import { EventEmitter2 } from 'eventemitter2';

const debug = require('debug')('insteon-plm:modem');

import { frameFromBytes, frameToBytes, Commands } from './protocol';
import { WorkQueue } from './workQueue';
import { CommandNAKError } from './errors';
import { timed } from './utils';

export { Commands };

const DEFAULT_TIMEOUT_MS = 5000; // ms

/**
 * The Modem represents the connection to an INSTEON PLM (PowerLinc Modem).
 */
export class Modem extends EventEmitter2 {
  constructor({ retries = 25 } = {}) {
    super({ wildcards: true });
    this._buffer = Buffer.alloc(0);
    this._queue = new WorkQueue({ retries, retryErrors: [CommandNAKError] });
  }

  /**
   * Whether the modem is open or not.
   * @type {boolean}
   */
  get isOpen() {
    return this._port && this._port.isOpen;
  }

  /**
   * Opens the modem. Throws if it has already been opened.
   *
   * @returns {Promise} Promise which is resolved when opening is completed.
   */
  open(path) {
    if (this.isOpen) {
      throw new Error('Already open');
    }

    debug(`Opening ${path}...`);
    this._port = new SerialPort(path, {
      baudRate: 19200,
      databits: 8,
      stopbits: 1,
      parity: 'none',
      autoOpen: false
    });

    return new Promise((resolve, reject) => {
      this._port.open(function(err) {
        if (err) {
          debug(`Failed to open ${path}...`, err);
          return reject(err);
        }

        debug(`Successfully opened ${path}...`);
        return resolve();
      });

      this._port.on('data', data => {
        this._handleData(data);
      });
    });
  }

  _drain() {
    return new Promise((resolve, reject) => {
      this._port.drain(err => {
        if (err) {
          debug('Error while draining', err);
          reject(err);
        } else {
          resolve();
        }
      });
    });
  }

  async _write(buf) {
    debug('Writing', buf);
    if (!this._port.isOpen) {
      return Promise.reject(new Error('Port is not open'));
    }

    this._port.write(buf);
    await this._drain();
  }

  /**
   * Closes the modem.
   *
   * @returns {Promise}
   */
  async close() {
    if (!this._port || !this._port.isOpen) {
      return Promise.resolve();
    }

    return new Promise((resolve, reject) => {
      return this._port.close(function(err) {
        if (err) {
          debug('Failed to close', err);
          reject(err);
        } else {
          debug('Closed.');
          resolve();
        }
      });
    });
  }

  /**
   * Sends a frame to the modem.
   *
   * @param {Object} frame
   * @param {number} frame.command - The command to send
   * @param {Object} frame.payload - The payload of the command
   * @param {Object} options - options
   * @param {number} options.timeout - Number of millis to wait.
   * @returns {Promise} The frame returned from the modem.
   */
  async send(frame, { timeout = DEFAULT_TIMEOUT_MS, retry = true } = {}) {
    debug('Queuing outgoing frame', frame);
    return timed(this._queue.enqueue(async attempt => {
      debug(`Sending frame (attempt ${attempt})`, frame);
      await this._write(frameToBytes(frame));
      const reply = await this.matchFrame(frame.command, { timeout });

      if (reply.ack) {
        debug('Frame ACK', reply);
        return reply;
      }

      debug('Frame NAK', reply);
      if (retry) {
        throw new CommandNAKError(reply);
      }

      return retry;
    }), timeout);
  }

  /**
   * Gets the modem information.
   */
  async getInfo() {
    const { payload } = await this.send({ command: Commands.GET_IM_INFO });
    return payload;
  }

  /**
   * Gets the modem configuration.
   */
  async getConfig() {
    const { payload } = await this.send({ command: Commands.GET_IM_CONFIG });
    return payload;
  }

  /**
   * Sets the modem configuration.
   */
  async setConfig(newConfig) {
    const payload = Object.assign({}, await this.getConfig(), newConfig);

    await this.send({
      command: Commands.SET_IM_CONFIG,
      payload
    });

    return this.getConfig();
  }

  /**
   * Exits linking mode previously started by {@link startAllLinking}.
   */
  async cancelAllLinking() {
    await this.send({ command: Commands.CANCEL_ALL_LINKING });
  }

  async factoryReset() {
    await this.send({ command: Commands.FACTORY_RESET }, { timeout: 30000 });
  }

  /**
   * Enters the modem into linking mode.
   */
  async startAllLinking({ controller = undefined, remove = false, group = 1 } = {}) {
    await this.send({ command: Commands.START_ALL_LINKING, payload: { controller, group, remove } });
  }

  /**
   * Waits for one frame with the specified command.
   */
  async matchFrame(command, { timeout = DEFAULT_TIMEOUT_MS } = {}) {
    let onMatchedFrame;
    const frameHandler = frame => {
      if (frame.command === command) {
        onMatchedFrame(frame);
      }
    }

    return timed(new Promise(resolve => {
      onMatchedFrame = resolve;
      this.on('frame', frameHandler);
    }), timeout).then(val => {
      this.off('frame', frameHandler);
      return val;
    }, err => {
      this.off('frame', frameHandler);
      throw err;
    });
  }

  /**
   * Returns the modem link database.
   */
  async getAllLinkDatabase() {
    let records = [];
    let frame;

    const firstReply = await this.send({ command: Commands.GET_FIRST_ALL_LINK }, { retry: false });
    if (!firstReply.ack) {
      return records;
    }

    frame = await this.matchFrame(Commands.ALL_LINK_RECORD_RESPONSE);
    records.push(frame.payload);

    while (true) {
      const reply = await this.send({ command: Commands.GET_NEXT_ALL_LINK }, { retry: false });
      if (!reply.ack) {
        // This indicates there are no more records.
        break;
      }

      frame = await this.matchFrame(Commands.ALL_LINK_RECORD_RESPONSE);
      records.push(frame.payload);
    }

    return records;
  }

  _handleData(data) {
    if (data) {
      this._buffer = Buffer.concat([this._buffer, data]);
    }

    debug('Buffer', this._buffer);

    try {
      const [frame, count] = frameFromBytes(this._buffer);
      this._buffer = this._buffer.slice(count);

      debug('Received Frame', frame);
      this.emit('frame', frame);

      // We want to be able to do work in between frames.
      setImmediate(() => this._handleData());
    } catch (error) {
      if (!(error instanceof RangeError)) {
        console.error('Error parsing modem frame', error);
      }
    }
  }
}
