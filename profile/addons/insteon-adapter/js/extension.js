(function() {
  class InsteonExtension extends window.Extension {
    constructor() {
      super('insteon-adapter');
      this.addMenuEntry('INSTEON');

      this.content = '';
      fetch(`/extensions/${this.id}/views/content.html`)
        .then((res) => res.text())
        .then((text) => {
          this.content = text;
        })
        .catch((e) => console.error('Failed to fetch content:', e));
    }

    showStatus(status) {
      if (this._statusTimeout) {
        clearTimeout(this._statusTimeout);
      }

      const statusContainer = document.getElementById('insteon-status');
      statusContainer.innerText = status;
      this._statusTimeout = setTimeout(() => {
        statusContainer.innerText = '';
      }, 3000);
    }

    async show() {
      this.view.innerHTML = this.content;

      const scanButton = document.getElementById('button-scan');
      const hbDeviceSelect = document.getElementById('heartbeat-device-selection');
      const setupHbButton = document.getElementById('setup-heartbeat-button');
      const lowbatDeviceSelect = document.getElementById('lowbat-device-selection');
      const setupLowbatButton = document.getElementById('setup-lowbat-button');
      scanButton.addEventListener('click', async () => {
        this.showStatus('Scanning database and adding things...');
        try {
          scanButton.disabled = true;
          const result = await window.API.postJson(`/extensions/${this.id}/api/scan`);
          this.showStatus(`Scan complete. Added ${result.length} things.`);
        } finally {
          scanButton.disabled = false;
        }
      });

      const batteryDevices = (await window.API.getThings()).filter((thing) => {
        return thing.id.includes('/insteon-') && thing['@type'].includes('BatteryPowered');
      });

      for (const device of batteryDevices) {
        const opt = document.createElement('option');
        opt.value = device.id.split('/').pop();
        opt.textContent = device.title;
        hbDeviceSelect.add(opt.cloneNode(true));
        lowbatDeviceSelect.add(opt);
      }

      if (batteryDevices.length === 0) {
        const opt = document.createElement('option');
        opt.textContent = 'No Devices';
        hbDeviceSelect.add(opt);
        hbDeviceSelect.disabled = true;
      }

      setupHbButton.addEventListener('click', async () => {
        const deviceId = hbDeviceSelect.value;

        try {
          this.showStatus(`Setting up Heartbeat event for ${deviceId}...`);
          setupHbButton.disabled = true;
          await window.API.postJson(`/extensions/${this.id}/api/setupHeartbeat`, { deviceId });
          this.showStatus(`Successfuly set up Heartbeat event for ${deviceId}`);
        } catch (e) {
          this.showStatus(`Failed to setup Heartbeat event for ${deviceId}`);
        } finally {
          setupHbButton.disabled = false;
        }
      });

      setupLowbatButton.addEventListener('click', async () => {
        const deviceId = hbDeviceSelect.value;

        try {
          this.showStatus(`Setting up LowBattery event for ${deviceId}...`);
          setupHbButton.disabled = true;
          await window.API.postJson(`/extensions/${this.id}/api/setupLowBattery`, { deviceId });
          this.showStatus(`Successfuly set up LowBattery event for ${deviceId}`);
        } catch (e) {
          this.showStatus(`Failed to setup LowBattery event for ${deviceId}`);
        } finally {
          setupHbButton.disabled = false;
        }
      });
    }
  }

  new InsteonExtension();
})();
