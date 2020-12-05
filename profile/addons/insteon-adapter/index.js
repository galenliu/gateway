/**
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/.
 */

'use strict';

const { Hub } = require('insteon-plm');

const InsteonAdapter = require('./insteon-adapter');

module.exports = async (addonManager, manifest, errorCallback) => {
  const { path } = manifest.moziot.config;

  const hub = new Hub();
  try {
    await hub.open(path);
  } catch (e) {
    errorCallback(manifest.id, `Unable to open device: ${e}`);
    return;
  }

  new InsteonAdapter(addonManager, manifest, hub);
};
