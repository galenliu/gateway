'use strict';

const { APIHandler, APIResponse } = require('gateway-addon');
const manifest = require('./manifest.json');

class InsteonAPIHandler extends APIHandler {
  constructor(addonManager, adapter) {
    super(addonManager, manifest.id);
    addonManager.addAPIHandler(this);

    this.adapter = adapter;
  }

  async handleRequest(request) {
    if (request.method !== 'POST') {
      return new APIResponse({ status: 404 });
    }

    let result = {};
    switch (request.path) {
      case '/scan':
        result = await this.adapter.scan();
        break;
      case '/setupHeartbeat': {
        const device = this.adapter.devices[request.body.deviceId];
        if (!device) {
          return new APIResponse({ status: 404 });
        }

        result = await device.link({ controller: false, group: 4 });
        break;
      }
      case '/setupLowBattery': {
        const device = this.adapter.devices[request.body.deviceId];
        if (!device) {
          return new APIResponse({ status: 404 });
        }

        result = await device.link({ controller: false, group: 3 });
        break;
      }
      default:
        return new APIResponse({ status: 404 });
    }

    return new APIResponse({
      status: 200,
      contentType: 'application/json',
      content: JSON.stringify(result),
    });
  }
}

module.exports = InsteonAPIHandler;
