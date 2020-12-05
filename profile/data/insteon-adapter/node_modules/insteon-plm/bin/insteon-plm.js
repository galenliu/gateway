#!/usr/bin/env node

// Make things work locally and globally
let modulePath = 'insteon-plm';
try {
  require(modulePath);
} catch (e) {
  modulePath = '../lib';
}

const { Hub, Messenger, MessageCommands, Modem } = require(modulePath);
const program = require('commander');

const pkg = require('../package.json');

let modem = null;
let messenger = null;
let hub = null;

function wrapAction(work) {
  return async function(...args) {
    try {
      if (!modem) {
        modem = new Modem();
        await modem.open(program.device);
        await modem.cancelAllLinking();

        messenger = new Messenger(modem);
        hub = new Hub(messenger);
      }
    } catch (e) {
      console.error('Failed to open modem', e);
      process.exit(1);
    }

    const close = async () => {
      await modem.cancelAllLinking();
      await modem.close();
    }
      
    try {
      await work.apply(this, args);
      await close();
      process.exit(0);
    } catch (e) {
      console.error(e);

      await close();
      process.exit(1);
    }
  }
}

program.version(pkg.version);
program.requiredOption('-d, --device <path>', 'Device path, e.g. /dev/ttyUSB0');

program
  .command('info')
  .description("Show modem information")
  .action(wrapAction(async () => {
    const info = await modem.getInfo();
    console.log(info);
  }));

program
  .command('links [address]')
  .description("Show a device's link database, or that of the modem when address is omitted")
  .action(wrapAction(async address => {
    const links = await hub.links(address);
    console.log(links);
  }));

program
  .command('monitor')
  .description('Show incoming messages')
  .action(wrapAction(async () => {
    messenger.on('message', message => {
      console.log(message);
    });

    // Run forever
    return new Promise(() => {});
  }));

program
  .command('link [address]')
  .description("Link a new device. If address is omitted, press the set button on the desired device.")
  .option("-r, --responder", "Link modem as a responder")
  .option("-c, --controller", "Link modem as a controller")
  .option("-g, --group <number>", "The group number to link")
  .action(wrapAction(async (address, cmd) => {
    try {
      let controller = undefined;
      if (cmd.controller) {
        controller = true;
      } else if (cmd.responder) {
        controller = false;
      }

      let group = undefined;
      if (cmd.group) {
        group = parseInt(cmd.group);
      }

      const info = await hub.link(address, {
        controller, group,
        timeout: address ? 3000 : 30000 // 30s if we're waiting to press set button, 3s otherwise
      });
      console.log('Linking complete', info);
    } catch (e) {
      console.log('Linking failed', e);
    }
  }));

program
  .command('unlink [address]')
  .description("Unlink a device. If address is omitted, press the set button on the desired device.")
  .action(wrapAction(async address => {
    try {
      const info = await hub.unlink(address, {
        timeout: address ? 3000 : 30000 // 30s if we're waiting to press set button, 3s otherwise
      });
      console.log('Unlinking complete', info);
    } catch (e) {
      console.log('Unlinking failed', e);
    }
  }));

program
  .command('reset')
  .description("Performs a factory reset on the modem.")
  .action(wrapAction(async () => {
    await modem.factoryReset();
  }));

program
  .command('status <address>')
  .description('Request the current status of a device')
  .action(wrapAction(async address => {
    console.log(await hub.status(address));
  }));


program
  .command('ping <address>')
  .description('Ping a device')
  .action(wrapAction(async address => {
    await messenger.send({ to: address, cmd1: MessageCommands.PING });
  }));

program
  .command('turnOn <address> [level]')
  .description('Turn on a device with an optional level, 0-100')
  .action(wrapAction(async (address, level) => {
    await hub.turnOn(address, { level: parseInt(level || 100) });
  }));

program
  .command('turnOff <address>')
  .description('Turn off a device')
  .action(wrapAction(async address => {
    await hub.turnOff(address);
  }));

async function main() {
  try {
    program.parse(process.argv);
  } catch (e) {
    console.error(e);
  }
}

main();
