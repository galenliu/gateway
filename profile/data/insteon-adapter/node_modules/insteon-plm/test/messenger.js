import { Messenger, MessageCommands } from '../lib';

import chai from 'chai';
import asPromised from 'chai-as-promised';

chai.use(asPromised);
const expect = chai.expect;

const modemPath = process.env['MODEM_PATH'];
let messenger = new Messenger();

const switchAddress = process.env['SWITCH_ADDRESS'];
describe('Messenger', function() {
  // Some tests take a bit due to retries.
  this.timeout(10000);
  this.slow(5000);

  before(function() {
    if (!modemPath) {
      console.warn('Skipping messenger tests due to missing MODEM_PATH environment variable');
      this.skip();
    }

    if (!switchAddress) {
      console.warn('Skipping messenger tests due to missing SWITCH_ADDRESS environment variable');
      this.skip();
    }
  });

  beforeEach(function() {
    return messenger.open(modemPath);
  });

  afterEach(function() {
    return messenger.close();
  });

  it('should require "to"', function() {
    expect(messenger.send({ cmd1: 0x12 })).rejectedWith(Error);
  });

  it('should require "cmd1"', function() {
    expect(messenger.send({ to: 'ffffff' })).rejectedWith(Error);
  });

  it('should require valid-looking "to"', function() {
    expect(messenger.send({ to: 5 })).rejectedWith(Error);
  });

  it('should turn on the switch', async function() {
    const reply = await messenger.send({ to: switchAddress, cmd1: MessageCommands.ON, cmd2: 0xff });
    expect(reply.flags.ack).to.be.true;
    expect(reply.flags.direct).to.be.true;
    expect(reply.cmd1).to.equal(MessageCommands.ON);
  });

  it('should emit "message" events', async function() {
    return new Promise(resolve => {
      messenger.once('message', message => {
        expect(message.cmd1).to.be.ok;
        resolve();
      });

      messenger.send({ to: switchAddress, cmd1: MessageCommands.ON, cmd2: 0xff });
    });
  });
});
