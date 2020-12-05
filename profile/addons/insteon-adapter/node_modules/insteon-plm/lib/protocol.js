import Protocol from 'bin-protocol';

export const START = 0x02;
export const ACK = 0x06;
export const NAK = 0x15;

export const EXTENDED_USERDATA_LENGTH = 14;

export const Commands = Object.freeze({
  STANDARD_INSTEON_RECEIVED: 0x50,
  EXTENDED_INSTEON_RECEIVED: 0x51,
  ALL_LINKING_COMPLETED: 0x53,
  ALL_LINK_RECORD_RESPONSE: 0x57,
  GET_IM_INFO: 0x60,
  SEND_INSTEON: 0x62,
  START_ALL_LINKING: 0x64,
  CANCEL_ALL_LINKING: 0x65,
  FACTORY_RESET: 0x67,
  GET_FIRST_ALL_LINK: 0x69,
  GET_NEXT_ALL_LINK: 0x6a,
  SET_IM_CONFIG: 0x6b,
  GET_IM_CONFIG: 0x73
});

export const PLMProtocol = Protocol.createProtocol();
PLMProtocol.define('start', {
  read() {
    while (this.context.start != START) {
      this.UInt8('start');
    }

    // We just want to advance the stream, no need to store the value of the start byte.
    delete this.context.start;
  },
  write() {
    this.UInt8(START);
  }
});

const NO_ACK_COMMANDS = Object.freeze([
  Commands.STANDARD_INSTEON_RECEIVED,
  Commands.EXTENDED_INSTEON_RECEIVED,
  Commands.ALL_LINK_RECORD_RESPONSE,
  Commands.ALL_LINKING_COMPLETED
]);

PLMProtocol.define('end', {
  read() {
    this.path = [];

    if (NO_ACK_COMMANDS.indexOf(this.result.command) >= 0) {
      // No terminator expected for these, just return.
      return;
    }

    // In the case of an unknown frame or extended message, we may have some bytes
    // to read before getting to the terminator. Read those into a generic buffer first.
    const remainder = [];
    while (this.context.end !== ACK && this.context.end !== NAK) {
      this.UInt8('end');
      remainder.push(this.context.end);
    }

    // Remove the last byte from payload, which is actually the terminator
    remainder.pop();

    this.context.ack = (this.context.end == ACK);
    delete this.context.end;

    if (this.result.payload && this.result.payload.flags && this.result.payload.flags.extended) {
      if (remainder.length !== EXTENDED_USERDATA_LENGTH) {
        throw new RangeError(`Expected ${EXTENDED_USERDATA_LENGTH} bytes of userdata, have ${remainder.length}`);
      }

      this.context.payload.userdata = Buffer.from(remainder);
    } else if (remainder.length > 1) {
      this.context.payload = { data: Buffer.from(remainder) };
    }
  },

  write() {
    // We don't write ACK/NAK to the PLM.
  }
});

PLMProtocol.define('command', {
  read() {
    this.UInt8('command');
  },
  write(command) {
    this.UInt8(command);
  }
});

PLMProtocol.define('address', {
  read() {
    this.loop('bytes', this.UInt8, 3);
    return Buffer.from(this.context.bytes).toString('hex');
  },
  write(val) {
    this.loop(Buffer.from(val, 'hex'), this.UInt8);
  }
});

PLMProtocol.define('insteonFlags', {
  read() {
    this.UInt8('bits');
    const bits = this.context.bits;

    return {
      broadcast: !!(bits & (1 << 7)),
      allLink: !!(bits & (1 << 6)),
      direct: !(bits & (1 << 7)) && !(bits & (1 << 6)),
      ack: !!(bits & (1 << 5)),
      extended: !!(bits & (1 << 4)),
      hopsRemaining: ((bits & (1 << 3)) | (bits & (1 << 2))) >> 2,
      maxHops: (bits & (1 << 1)) | (bits & 1)
    }
  },

  write({ broadcast, allLink, ack, extended, hopsRemaining, maxHops }) {
    let bits = 0;

    if (broadcast) {
      bits |= (1 << 7);
    }

    if (allLink) {
      bits |= (1 << 6);
    }

    if (ack) {
      bits |= (1 << 5);
    }

    if (extended) {
      bits |= (1 << 4);
    }

    bits |= (hopsRemaining << 2);
    bits |= maxHops;

    this.UInt8(bits);
  }
});

PLMProtocol.define('linkRecordFlags', {
  read() {
    this.UInt8('flag');
    const { flag } = this.context;

    // There's other stuff in here but isn't really useful.
    return {
      controller: !!(flag & (1 << 6))
    }
  }
})

PLMProtocol.define('linkRecord', {
  read() {
    this.UInt8('unused');
    this.UInt8('byte');
    if (this.context.byte !== 1) {
      throw new Error('Unexpected value');
    }

    // this.loop('at', this.UInt16, 2);
    this.UInt16BE('at');
    this.UInt8('unused');
    this.linkRecordFlags('flags');
    this.UInt8('group');
    this.address('address');

    return {
      at: this.context.at,
      flags: this.context.flags,
      group: this.context.group,
      address: this.context.address
    };
  }
});

function readModemConfig(bits) {
  return {
    automaticLinking: !!(bits & (1 << 6)),
    monitorMode: !!(bits & (1 << 5)),
    automaticLED: !!(bits & (1 << 4)),
    deadman: !!(bits & (1 << 3)),
  };
}

PLMProtocol.define('payload', {
  read() {
    this.path = ['payload'];

    switch (this.result.command) {
    case Commands.STANDARD_INSTEON_RECEIVED:
      this.address('from');
      this.address('to');
      this.insteonFlags('flags');
      this.UInt8('cmd1');
      this.UInt8('cmd2');
      break;
    case Commands.EXTENDED_INSTEON_RECEIVED:
      this.address('from');
      this.address('to');
      this.insteonFlags('flags');
      this.UInt8('cmd1');
      this.UInt8('cmd2');
      this.loop('userdata', this.UInt8, EXTENDED_USERDATA_LENGTH);
      this.context.userdata = Buffer.from(this.context.userdata);
      break;
    case Commands.ALL_LINKING_COMPLETED:
      this.UInt8('code');
      this.UInt8('group');
      this.address('address');
      this.UInt8('category');
      this.UInt8('subcategory');
      this.UInt8('firmwareVersion');

      this.context.controller = this.context.code === 1;
      delete this.context.code;
      break;
    case Commands.ALL_LINK_RECORD_RESPONSE:
      this.UInt8('flags');
      this.UInt8('group');
      this.address('address');
      this.loop('linkData', this.UInt8, 3);

      this.context.controller = !!(this.context.flags & (1 << 6));
      delete this.context.flags;
      break;
    case Commands.GET_IM_INFO:
      this.address('address');
      this.UInt8('category');
      this.UInt8('subcategory');
      this.UInt8('firmwareVersion');
      break;
    case Commands.SEND_INSTEON:
      this.address('to');
      this.insteonFlags('flags');
      this.UInt8('cmd1');
      this.UInt8('cmd2');
      break;
    case Commands.GET_IM_CONFIG:
      this.UInt8('bits');
      this.UInt16LE('spare');
      return readModemConfig(this.context.bits);
    case Commands.SET_IM_CONFIG:
      this.UInt8('bits');
      return readModemConfig(this.context.bits);
    case Commands.START_ALL_LINKING:
      this.UInt8('code');
      this.UInt8('group');
      return {
        controller: this.context.code === 1,
        remove: this.context.code === 0xff,
        group: this.context.group
      };
    default:
      // Unknown frame. We'll read the remainder up to NAK/ACK in end()
      return null;
    }
  },

  write(frame) {
    const { payload, command } = frame;
    switch (command) {
    case Commands.SET_IM_CONFIG: {
      const { automaticLinking, monitorMode, automaticLED, deadman } = payload;
      let bits = 0;

      if (automaticLinking) {
        bits |= (1 << 6);
      }

      if (monitorMode) {
        bits |= (1 << 5);
      }

      if (automaticLED) {
        bits |= (1 << 4);
      }

      if (deadman) {
        bits |= (1 << 3);
      }

      this.UInt8(bits);
      break;
    }
    case Commands.START_ALL_LINKING: {
      const { controller, group, remove } = payload;
      if (remove) {
        this.UInt8(0xff);
      } else if (controller === undefined) {
        // controller/master relationship determined by order of devices entering link mode
        this.UInt8(0x03);
      } else {
        // adapter is master/slave according to `controller` value
        this.UInt8(controller ? 1 : 0);
      }
      this.UInt8(group);
      break;
    }
    case Commands.SEND_INSTEON: {
      const { to, flags, cmd1, cmd2 } = payload;
      let { userdata } = payload;

      if (flags.extended) {
        if (!userdata || userdata.length != EXTENDED_USERDATA_LENGTH) {
          throw new Error(`Expected extended message userdata to be of length ${EXTENDED_USERDATA_LENGTH}`);
        }
      }

      this.address(to);
      this.insteonFlags(flags);
      this.UInt8(cmd1);
      this.UInt8(cmd2);

      if (flags.extended) {
        this.loop(userdata, this.UInt8);

        let sum = 0;

        // Get the sum of cmd1 -> D13
        for (let i = 6; i < (this.result.length - 1); i++) {
          sum += this.result[i];
        }

        const checksum = (~(sum) + 1) & 255;

        // Put the result in D14
        this.result[this.result.length - 1] = checksum;
      }

      break;
    }
    }
  }
});

export function linkRecordFromBytes(data) {
  const prot = new PLMProtocol();

  const { result } = prot.read(data)
    .linkRecord()

  return result;
}

export function frameFromBytes(data) {
  const prot = new PLMProtocol();

  const result = prot.read(data)
    .start()
    .command()
    .payload()
    .end();

  return [result.result, result.offset];
}

export function frameToBytes(frame) {
  const prot = new PLMProtocol();

  const result = prot.write()
    .start()
    .command(frame.command)
    .payload(frame)
    .end();

  return result.buffer.slice(0, result.offset);
}
