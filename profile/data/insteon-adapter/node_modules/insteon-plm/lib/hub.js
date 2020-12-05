import { EventEmitter2 } from 'eventemitter2';
import ow from 'ow';

const debug = require('debug')('insteon-plm:hub');

import { Commands } from './modem';
import { Messenger, MessageCommands } from './messenger';
import { WorkQueue } from './workQueue';
import { linkRecordFromBytes } from './protocol';
import { findProductInfo } from './database';
import { checkAddress, checkOptionalAddress } from './utils';
import { rateForDuration } from './rates';

const DEFAULT_TIMEOUT_MS = 3000;

export class Hub extends EventEmitter2 {

  constructor(messenger) {
    super({ wildcards: true });

    ow(messenger, ow.optional.object.instanceOf(Messenger));

    this._queue = new WorkQueue();

    if (messenger) {
      this.messenger = messenger;
    }
  }

  async open(path) {
    ow(path, ow.string);

    if (!this.messenger) {
      this.messenger = new Messenger();
    }

    await this.messenger.open(path);
  }

  async close() {
    if (this._messenger) {
      await this._messenger.close();
    }
  }

  get messenger() {
    return this._messenger;
  }

  set messenger(messenger) {
    ow(messenger, ow.object.instanceOf(Messenger));

    this._messenger = messenger;
  }

  get modem() {
    return this.messenger.modem;
  }

  async status(address,  { cmd2 = 0, timeout = DEFAULT_TIMEOUT_MS } = {}) {
    checkAddress(address);
    ow(timeout, ow.number.greaterThan(0));

    const reply = await this.messenger.send({
      to: address,
      cmd1: MessageCommands.STATUS_REQUEST,
      cmd2
    });

    return reply.cmd2;
  }

  async productInfo(address, { timeout = DEFAULT_TIMEOUT_MS } = {}) {
    checkAddress(address);
    ow(timeout, ow.number.greaterThan(0));

    await this.messenger.send({
      to: address,
      cmd1: MessageCommands.ID_REQUEST,
    });

    const reply = await this.messenger.matchMessage({
      from: address,
      flags: { broadcast: true }
    }, { timeout });

    const category = parseInt(reply.to.substring(0, 2), 16);
    const subcategory = parseInt(reply.to.substring(2, 4), 16);
    const firmwareVersion = parseInt(reply.to.substring(4, 6), 16);
    const { name, product } = findProductInfo(category, subcategory) || {
      name: null,
      product: null
    };

    return {
      category, subcategory, firmwareVersion, name, product
    };
  }

  async link(address = null, { controller = undefined, group = 1, timeout = DEFAULT_TIMEOUT_MS } = {}) {
    checkOptionalAddress(address);

    if (address) {
      ow(controller, ow.boolean);
    }

    ow(group, ow.number.greaterThan(0));
    ow(timeout, ow.number.greaterThan(0));

    await this.modem.cancelAllLinking();

    const matching = this.modem.matchFrame(Commands.ALL_LINKING_COMPLETED, { timeout });

    try {
      await this.modem.startAllLinking({ controller, group });

      // If we are targeting a specific device put it into linking mode.
      if (address) {
        await this.messenger.send({
          to: address,
          cmd1: MessageCommands.ENTER_LINKING,
          cmd2: group,
          flags: {
            extended: true
          }
        });
      }

      const linkInfo = (await matching).payload;
      const { name, product } = findProductInfo(linkInfo.category, linkInfo.subcategory) || {
        name: null,
        product: null
      };

      return {
        name, product, ...linkInfo
      };
    } finally {
      this.modem.cancelAllLinking().catch(err => {
        debug('Failed to cancel linking', err);
      });

      if (address) {
        this.messenger.send({
          to: address,
          cmd1: MessageCommands.EXIT_LINKING,
          flags: { extended: true }
        }).catch(e => {
          debug('Failed to exit linking', e);
        });
      }
    }
  }

  async unlink(address = null, { group = 1, timeout = DEFAULT_TIMEOUT_MS }) {
    checkOptionalAddress(address);

    ow(group, ow.number.greaterThan(0));
    ow(timeout, ow.number.greaterThan(0));

    await this.modem.cancelAllLinking();

    const matching = this.modem.matchFrame(Commands.ALL_LINKING_COMPLETED, { timeout });

    try {
      await this.modem.startAllLinking({ remove: true, group });

      // If we are targeting a specific device put it into linking mode.
      if (address) {
        await this.messenger.send({
          to: address,
          cmd1: MessageCommands.ENTER_LINKING,
          cmd2: group,
          flags: {
            extended: true
          }
        });
      }

      const info = await matching;
      return info;
    } finally {
      this.modem.cancelAllLinking().catch(err => {
        debug('Failed to cancel linking', err);
      });
    }
  }

  async links(address = null, { timeout = DEFAULT_TIMEOUT_MS } = {}) {
    checkOptionalAddress(address);
    ow(timeout, ow.number.greaterThan(0));

    if (!address) {
      return this.modem.getAllLinkDatabase();
    }

    await this.messenger.send({
      to: address,
      cmd1: MessageCommands.READ_WRITE_LINKS,
      flags: { extended: true }
    }, { timeout });

    const records = [];

    while (true) {
      const message = await this.messenger.matchMessage({
        from: address,
        flags: { extended: true },
        cmd1: MessageCommands.READ_WRITE_LINKS
      }, { timeout });

      const record = linkRecordFromBytes(message.userdata);
      if (record.group === 0) {
        // Indicates last record
        break;
      }

      records.push(record);
    }

    return records;
  }

  async turnOn(address, { level = 100, fast = false, duration = undefined, timeout = DEFAULT_TIMEOUT_MS } = {}) {
    checkAddress(address);
    ow(level, ow.number.inRange(0, 100));
    ow(fast, ow.boolean);
    ow(duration, ow.optional.number.greaterThan(0));
    ow(timeout, ow.number.greaterThan(0));

    if (fast && (duration || level < 100)) {
      throw new Error('Cannot specify "fast" with "duration"');
    }

    if (!duration) {
      await this.messenger.send({
        to: address,
        cmd1: fast ? MessageCommands.ON_FAST : MessageCommands.ON,
        cmd2: Math.round((level / 100) * 255)
      });
    } else {
      // The level and rate are compressed into a single byte.
      // The high bits are the level, and the low bits are the ramp rate.
      const levelAndRate = (Math.round((level / 100) * 15) << 4) | rateForDuration(duration);

      await this.messenger.send({
        to: address,
        cmd1: MessageCommands.ON_WITH_RATE,
        cmd2: levelAndRate
      });
    }
  }

  async turnOff(address, { fast = false, duration = undefined, timeout = DEFAULT_TIMEOUT_MS } = {}) {
    checkAddress(address);
    ow(fast, ow.boolean);
    ow(duration, ow.optional.number.greaterThan(0));
    ow(timeout, ow.number.greaterThan(0));

    if (fast && duration) {
      throw new Error('Cannot specify both "fast" with "duration"');
    }

    if (!duration) {
      await this.messenger.send({
        to: address,
        cmd1: fast ? MessageCommands.OFF_FAST : MessageCommands.OFF,
        cmd2: 0
      });
    } else {
      await this.messenger.send({
        to: address,
        cmd1: MessageCommands.OFF_WITH_RATE,
        cmd2: rateForDuration(duration)
      });
    }
  }
}
