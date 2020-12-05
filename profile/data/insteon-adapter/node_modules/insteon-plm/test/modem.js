import { Commands, Modem, CommandNAKError, TimeoutError } from '../lib';

import chai from 'chai';
import asPromised from 'chai-as-promised';

chai.use(asPromised);
const expect = chai.expect;

const modemPath = process.env['MODEM_PATH'];

const modem = new Modem({ retries: 3 });
describe('Modem', function() {
  // Some tests take a bit due to retries.
  this.timeout(10000);
  this.slow(5000);

  beforeEach(function() {
    if (!modemPath) {
      console.warn('Skipping modem tests due to missing MODEM_PATH environment variable');
      this.skip();
    }

    return modem.open(modemPath);
  });

  afterEach(function() {
    return modem.close();
  });

  it('should be open', function() {
    expect(modem.isOpen).to.be.true;
  });

  it('should send frames and wait for ack', async function() {
    const { ack } = await modem.send({ command: Commands.GET_IM_INFO });
    expect(ack).to.be.true;
  });

  it('should fail to send bogus/malformed frames', function() {
    return expect(modem.send({ command: 1 })).rejectedWith(CommandNAKError);
  });

  it(`should emit 'frame' events`, async function() {
    const waiting = new Promise(resolve => {
      modem.once('frame', frame => {
        expect(frame.command).to.be.ok;
        resolve();
      });
    });

    await modem.send({ command: Commands.GET_IM_INFO });
    await waiting;
  });

  it('send() should timeout', function() {
    // Timeout after 1ms
    return expect(modem.send({ command: 1 }, { timeout: 1 })).rejectedWith(TimeoutError);
  });

  it('send() should reject if closed', async function() {
    await modem.close();
    return expect(modem.send({ command: Commands.GET_IM_INFO })).rejectedWith(Error);
  });

  it('getInfo() should return info', async function() {
    const { address, category, subcategory, firmwareVersion } = await modem.getInfo();
    expect(address).to.be.ok;
    expect(category).to.be.ok;
    expect(subcategory).to.be.ok;
    expect(firmwareVersion).to.be.ok;
  });

  it('setConfig() should set config', async function() {
    const config = await modem.getConfig();

    const newConfig = await modem.setConfig({ automaticLED: !config.automaticLED });

    expect(config.automaticLED).to.not.equal(newConfig.automaticLED);
  });

  it('startAllLinking() should work', async function() {
    await modem.startAllLinking();
    await modem.cancelAllLinking();
  });

  it('getAllLinkDatabase() should work', async function() {
    const records = await modem.getAllLinkDatabase();
    expect(records).to.be.ok;
    expect(records.length).to.be.greaterThan(0);

    const { group, address, linkData, controller } = records.shift();
    expect(group).to.be.greaterThan(0);
    expect(address).to.be.ok;
    expect(linkData.length).to.equal(3);
    expect(controller).to.exist;
  });
});
