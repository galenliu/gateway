import ReopeningWebSocket from "./models/reopening-web-socket";
import GatewayModel from "./models/gateway-model";
import API from "./js/api";

const Core = {

    ORIGIN: window.location.origin,
    HOST: window.location.host,
    LANGUAGE: 'en-US',
    TIMEZONE: 'UTC',
    UNITS: {},

    init: function () {

        this.gatewayModel = new GatewayModel()
        //this.initWebSocket()
    },

    // showThings: function(context) {
    //     const events = context.pathname.split('/').pop() === 'events';
    //     ThingsScreen.show(context.params.thingId || null,
    //         context.params.actionName || null,
    //         events,
    //         context.querystring);
    //     this.selectView('things');
    // },

    initWebSocket() {
        const path = `${this.ORIGIN.replace(/^http/, 'ws')}/internal-logs?jwt=${API.jwt}`;
        this.ws = new ReopeningWebSocket(path);
        this.ws.addEventListener(
            'message',
            (msg) => {
                const message = JSON.parse(msg.data);
                if (message && message.message) {
                    this.showMessage(message.message, 5000, message.url);
                }
            }
        );
    }

}

export default Core;