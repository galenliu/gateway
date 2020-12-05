import delay from 'delay';

export class WorkQueue {
  constructor({ retries = 0, retryDelay = 100, retryErrors = [] } = {}) {
    this.retries = retries;
    this.retryDelay = retryDelay;
    this.retryErrors = retryErrors;

    this._queue = [];
    this._running = false;
  }

  async _runner() {
    try {
      if (!this._queue.length) {
        throw new Error('Runner scheduled with no work!');
      }

      const { work, resolve, reject } = this._queue.shift();
      let lastException = null;
      for (let i = 0; i < this.retries + 1; i++) {
        try {
          resolve(await work(i + 1));
          lastException = null;
          break;
        } catch (e) {
          lastException = e;

          if (!this.retryErrors.includes(e.constructor)) {
            break;
          }

          await delay(this.retryDelay);
        }
      }

      if (lastException) {
        reject(lastException);
      }
    } finally {
      this._running = false;
    }

    // Do more work if we need to
    this._ensureRunner();
  }

  _ensureRunner() {
    if (!this._running && this._queue.length > 0) {
      this._running = true;
      process.nextTick(() => this._runner());
    }
  }

  enqueue(work) {
    return new Promise((resolve, reject) => {
      this._queue.push({
        work, resolve, reject
      });

      this._ensureRunner();
    });
  }
}
