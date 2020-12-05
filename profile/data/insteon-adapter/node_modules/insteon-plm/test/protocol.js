import { Commands, START, ACK, frameFromBytes, frameToBytes, linkRecordFromBytes } from '../lib/protocol';

import { expect } from 'chai';

function readOneFrame(bytes) {
  const [frame, count] = frameFromBytes(bytes);
  expect(bytes.length - count).to.equal(0);
  return frame;
}

function testWrite(a, b) {
  expect(frameToBytes(a))
    .to.deep.equal(Buffer.from(b));
}

function testRead(a, b) {
  expect(readOneFrame(Buffer.from(a))).to.deep.equal(b);
}

function testSymmetric(frame) {
  const bytesWithAck = Buffer.concat([frameToBytes(frame), Buffer.from([ACK])]);
  expect(Object.assign({}, frame, { ack: true })).to.deep.equal(readOneFrame(bytesWithAck));
}

describe("PLMProtocol", function() {
  it('write GET_IM_INFO', function() {
    testWrite({ command: Commands.GET_IM_INFO },
      [START, Commands.GET_IM_INFO]);
  });

  it('read GET_IM_INFO', function() {
    testRead([START,
      Commands.GET_IM_INFO,
      0xf1, 0xf2, 0xf3, // address
      0x1, // category
      0x2, // subcategory
      0x3, // firmwareVersion
      ACK
    ], {
      ack: true,
      command: Commands.GET_IM_INFO,
      payload: {
        address: 'f1f2f3',
        category: 0x1,
        subcategory: 0x2,
        firmwareVersion: 0x3
      }
    });
  });

  it('write GET_IM_CONFIG', function() {
    expect(frameToBytes({ command: Commands.GET_IM_CONFIG }))
      .to.deep.equal(Buffer.from([START, Commands.GET_IM_CONFIG]));
  });

  it('read GET_IM_CONFIG', function() {
    testRead([
      START, Commands.GET_IM_CONFIG, 0xff, 0x00, 0x00, ACK
    ], {
      ack: true,
      command: Commands.GET_IM_CONFIG,
      payload: {
        automaticLinking: true,
        monitorMode: true,
        automaticLED: true,
        deadman: true
      }
    });
  });


  it('read/write SEND_INSTEON standard length', function() {
    testSymmetric({
      command: Commands.SEND_INSTEON,
      payload: {
        to: 'f1f2f3',
        flags: {
          broadcast: false,
          allLink: false,
          direct: true,
          ack: false,
          extended: false,
          hopsRemaining: 3,
          maxHops: 3
        },
        cmd1: 0x1,
        cmd2: 0x2
      }
    });
  });

  it('read/write SEND_INSTEON extended length', function() {
    const userdata = Buffer.alloc(14, 0xee);

    // Precomputed checksum value
    userdata[13] = 0xe7;

    testSymmetric({
      command: Commands.SEND_INSTEON,
      payload: {
        to: 'f1f2f3',
        flags: {
          broadcast: false,
          allLink: false,
          direct: true,
          ack: false,
          extended: true,
          hopsRemaining: 3,
          maxHops: 3
        },
        cmd1: 0x1,
        cmd2: 0x2,
        userdata
      }
    });
  });

  it('read STANDARD_INSTEON_RECEIVED', function() {
    testRead([START, Commands.STANDARD_INSTEON_RECEIVED,
      0xf1, 0xf2, 0xf3, // from
      0xe1, 0xe2, 0xe3, // to
      0xf, // flags
      0x1, // cmd1
      0x2, // cmd2
    ], {
      command: Commands.STANDARD_INSTEON_RECEIVED,
      payload: {
        from: 'f1f2f3',
        to: 'e1e2e3',
        flags: {
          broadcast: false,
          allLink: false,
          direct: true,
          ack: false,
          extended: false,
          hopsRemaining: 3,
          maxHops: 3
        },
        cmd1: 0x1,
        cmd2: 0x2
      }
    });
  });

  it('read EXTENDED_INSTEON_RECEIVED', function() {
    const userdata = [0xee, 0xee, 0xee, 0xee, 0xee, 0xee, 0xee, 0xee, 0xee, 0xee, 0xee, 0xee, 0xee, 0xd5];

    testRead([START, Commands.EXTENDED_INSTEON_RECEIVED,
      0xf1, 0xf2, 0xf3, // from
      0xe1, 0xe2, 0xe3, // to
      0xf, // flags
      0x1, // cmd1
      0x2, // cmd2
    ].concat(userdata), {
      command: Commands.EXTENDED_INSTEON_RECEIVED,
      payload: {
        from: 'f1f2f3',
        to: 'e1e2e3',
        flags: {
          broadcast: false,
          allLink: false,
          direct: true,
          ack: false,
          extended: false,
          hopsRemaining: 3,
          maxHops: 3
        },
        cmd1: 0x1,
        cmd2: 0x2,
        userdata: Buffer.from(userdata)
      }
    });
  });

  it('read link record data', function() {
    const userdata = [0x00, 0x01, 0x0f, 0xef, 0x02, 0xaa, 0x01, 0xaa, 0xbb, 0xcc, 0xff, 0x1c, 0x01, 0x69];

    expect(linkRecordFromBytes(Buffer.from(userdata))).to.deep.equal({
      address: "aabbcc",
      at: 4079,
      flags: {
        controller: false
      },
      group: 1
    });
  })
});
