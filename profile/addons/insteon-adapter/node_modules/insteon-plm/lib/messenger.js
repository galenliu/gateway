import { EventEmitter2 } from 'eventemitter2';
import { isMatch } from 'lodash/lang';
import ow from 'ow';

const debug = require('debug')('insteon-plm:messenger');

import { Modem, Commands } from './modem';
import { WorkQueue } from './workQueue';
import { EXTENDED_USERDATA_LENGTH } from './protocol';
import { MessageNAKError, UnexpectedMessageError } from './errors';
import { timed, checkAddress } from './utils';

const DEFAULT_TIMEOUT_MS = 3000;
const DEFAULT_FLAGS = {
  broadcast: false,
  allLink: false,
  direct: false,
  extended: false,
  hopsRemaining: 3,
  maxHops: 3
};

export const MessageCommands = {
  HEARTBEAT: 0x04,
  DIMMING_COMPLETE: 0x06,
  EXIT_LINKING: 0x08,
  ENTER_LINKING: 0x09,
  ENTER_UNLINKING: 0x0a,
  ID_REQUEST: 0x10,
  PING: 0x0f,
  ON: 0x11,
  ON_FAST: 0x12,
  OFF: 0x13,
  OFF_FAST: 0x14,
  START_MANUAL_CHANGE: 0x17,
  STOP_MANUAL_CHANGE: 0x18,
  STATUS_REQUEST: 0x19,
  GET_OP_FLAGS: 0x1f,
  SET_OP_FLAGS: 0x20,
  ON_WITH_RATE: 0x2e,
  OFF_WITH_RATE: 0x2f,
  READ_WRITE_LINKS: 0x2f,
};

export class Messenger extends EventEmitter2 {
  constructor(modem) {
    super({ wildcards: true });

    this._queue = new WorkQueue({ retries: 25, retryErrors: [MessageNAKError, UnexpectedMessageError] });
    this._handleFrame = this._handleFrame.bind(this);

    if (modem) {
      this.modem = modem;
    }
  }

  get modem() {
    return this._modem;
  }

  set modem(modem) {
    ow(modem, ow.object.instanceOf(Modem));

    if (this._modem) {
      this._modem.off('frame', this._handleFrame);
    }

    modem.on('frame', this._handleFrame);
    this._modem = modem;
  }

  async open(path) {
    ow(path, ow.string);

    if (!this._modem) {
      this._modem = new Modem();
    }

    await this._modem.open(path);
    this._modem.on('frame', this._handleFrame);
  }

  close() {
    this._modem.off('frame', this._handleFrame);
    return this._modem.close();
  }

  _handleFrame({ command, payload }) {
    if (command !== Commands.STANDARD_INSTEON_RECEIVED &&
        command !== Commands.EXTENDED_INSTEON_RECEIVED) {
      // We we only care about insteon frames.
      return;
    }

    debug('Received Message', payload);
    this.emit('message', payload);
  }

  async send(message, { timeout = DEFAULT_TIMEOUT_MS } = {}) {
    ow(message, ow.object.partialShape({
      to: checkAddress,
      cmd1: ow.number,
      cmd2: ow.optional.number,
      flags: ow.optional.object.partialShape({
        extended: ow.optional.boolean
      }),
      userdata: ow.optional.buffer.validate(value => ({
        validator: value.length === EXTENDED_USERDATA_LENGTH,
        message: `Expected userdata to be ${EXTENDED_USERDATA_LENGTH} bytes, got ${value.length}.`
      }))
    }));
    ow(timeout, ow.number.greaterThan(0));

    if (!message.to) {
      throw new Error(`Must specify 'to' in message`);
    } else if (typeof message.to !== 'string' || message.to.length !== 6) {
      throw new Error(`Message 'to' must be a string with 6 characters, e.g. ffffff`);
    }

    if (!message.cmd1) {
      throw new Error(`Must specify 'cmd1' in message`);
    }

    const fullMessage = Object.assign({}, { cmd2: 0 }, message);
    fullMessage.flags = Object.assign({}, DEFAULT_FLAGS, message.flags || {});

    if (fullMessage.userdata) {
      fullMessage.flags.extended = true;
    }

    if (fullMessage.flags.extended && !fullMessage.userdata) {
      fullMessage.userdata = Buffer.alloc(EXTENDED_USERDATA_LENGTH);
    }

    debug('Queueing Message', fullMessage);
    return this._queue.enqueue(async attempt => {
      debug(`Sending Message (attempt ${attempt})`, fullMessage);

      await this._modem.send({ command: Commands.SEND_INSTEON, payload: fullMessage });

      const reply = await this.matchMessage({
        from: fullMessage.to,
      }, { timeout });

      if (!reply.flags.ack) {
        debug('Received NAK Message', reply);
        throw new MessageNAKError(reply);
      }

      debug('Received ACK Message', reply);
      return reply;
    });
  }

  async matchMessage(criteria, { timeout = DEFAULT_TIMEOUT_MS } = {}) {
    ow(criteria, ow.object);
    ow(timeout, ow.number.greaterThan(0));

    let onMatchedMessage;
    const messageHandler = message => {
      if (isMatch(message, criteria)) {
        onMatchedMessage(message);
      }
    }

    return timed(new Promise(resolve => {
      onMatchedMessage = resolve;
      this.on('message', messageHandler);
    }), timeout).then(val => {
      this.off('message', messageHandler);
      return val;
    }, err => {
      this.off('message', messageHandler);
      throw err;
    });
  }
}
