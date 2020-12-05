import { WorkQueue } from '../lib/workQueue';

import chai from 'chai';
import asPromised from 'chai-as-promised';

chai.use(asPromised);
const expect = chai.expect;

function delay(ms) {
  return new Promise(resolve => {
    setTimeout(resolve, ms);
  });
}

describe('WorkQueue', function() {
  this.slow(1000);
  
  it('should do work', function() {
    const queue = new WorkQueue();

    let didWork = false;
    return queue.enqueue(() => {
      didWork = true;
    }).then(() => {
      expect(didWork).to.be.true;
    });
  });

  it('should execute sequentially', async function() {
    this.slow(500);

    const queue = new WorkQueue();
    let results = [];

    let promises = [];
    for (let i = 1; i < 5; i++) {
      promises.push(queue.enqueue(async () => {
        await delay(100 / i);

        results.push(i);
        return i;
      }));
    }

    expect(await Promise.all(promises)).to.deep.equal([1, 2, 3, 4]);
    expect(results).to.deep.equal([1, 2, 3, 4]);
  });

  it('should retry specified number of times', async function() {
    const queue = new WorkQueue({ retries: 1, retryErrors: [Error] });

    let count = 0;
    await expect(queue.enqueue(() => {
      count++;
      throw new Error('boom');
    })).rejectedWith(Error);

    expect(count).to.equal(2);
  });

  it('should retry on specified errors', async function() {
    const queue = new WorkQueue({ retries: 1, retryErrors: [] });

    let count = 0;
    await expect(queue.enqueue(() => {
      count++;
      throw new Error('boom');
    })).rejectedWith(Error);

    expect(count).to.equal(1);
  });
});
