import { Hub } from '../lib';

import chai from 'chai';
import asPromised from 'chai-as-promised';

chai.use(asPromised);
const expect = chai.expect;

const modemPath = process.env['MODEM_PATH'];
let hub = new Hub();

const switchAddress = process.env['SWITCH_ADDRESS'];
describe('Hub', function() {
  // Some tests take a bit due to retries.
  this.timeout(10000);
  this.slow(5000);

  before(function() {
    if (!modemPath) {
      console.warn('Skipping Hub tests due to missing MODEM_PATH environment variable');
      this.skip();
    }

    if (!switchAddress) {
      console.warn('Skipping Hub tests due to missing SWITCH_ADDRESS environment variable');
      this.skip();
    }
  });

  beforeEach(function() {
    return hub.open(modemPath);
  });

  afterEach(function() {
    return hub.close();
  });

  it('constructor() should check arguments', function() {
    expect(() => new Hub(42)).to.throw();
    expect(() => new Hub()).to.be.ok;
  });

  it('open() should check arguments', function() {
    expect(new Hub().open(42)).rejected;
  });

  it('should return product information', async function() {
    const info = await hub.productInfo(switchAddress);

    expect(info).to.have.property('category');
    expect(info).to.have.property('subcategory');
    expect(info).to.have.property('firmwareVersion');
    expect(info).to.have.property('name');
    expect(info).to.have.property('product');
  });

  it('should return device status', async function() {
    const value = await hub.status(switchAddress);
    expect(value).to.be.a('number');
  });

  it('should link as a controller', async function() {
    const link = await hub.link(switchAddress, { controller: true });
    expect(link).to.deep.include({
      group: 1,
      address: switchAddress,
      controller: true
    });

    expect(link).to.have.property('category');
    expect(link).to.have.property('subcategory');
    expect(link).to.have.property('firmwareVersion');
  });

  it('should link as a responder', async function() {
    const link = await hub.link(switchAddress, { controller: false });
    expect(link).to.deep.include({
      group: 1,
      address: switchAddress,
      controller: false
    });

    expect(link).to.have.property('category');
    expect(link).to.have.property('subcategory');
    expect(link).to.have.property('firmwareVersion');
    expect(link).to.have.property('name');
    expect(link).to.have.property('product');
  });

  it('should return link database', async function() {
    this.timeout(10000);

    const links = await hub.links(switchAddress);
    expect(links).to.be.an('array');
    expect(links).to.have.lengthOf.above(0);

    const link = links.shift();
    expect(link).to.have.property('address');
    expect(link).to.have.property('at');
    expect(link).to.have.property('flags');
    expect(link.flags).to.have.property('controller');
  });

  it('should turn lights on', async function() {
    await hub.turnOn(switchAddress);
  });

  it('should turn lights off', async function() {
    await hub.turnOff(switchAddress);
  });

  it('should turn lights on (fast)', async function() {
    await hub.turnOn(switchAddress, { fast: true });
  });

  it('should turn lights off (fast)', async function() {
    await hub.turnOff(switchAddress, { fast: true });
  });

  it('should turn lights on (level)', async function() {
    await hub.turnOn(switchAddress, { level: 50 });
  });

  it('should turn lights on (level, duration)', async function() {
    await hub.turnOn(switchAddress, { level: 50, duration: 5 });
  });

  it('should turn lights off (duration)', async function() {
    await hub.turnOff(switchAddress, { duration: 5 });
  });
});
