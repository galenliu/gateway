import useWebSocket, {ReadyState} from "react-use-websocket";
import {useMemo, useRef, useState} from "react";


export const useThingsWebsocket = (props) => {
    const didUnmount = useRef(false);
    //Public API that will echo messages sent to it back to the client
    const [socketUrl, setSocketUrl] = useState(props);
    const messageHistory = useRef([]);

    const {
        sendMessage,
        lastMessage,
        readyState,
    } = useWebSocket(socketUrl);

    messageHistory.current = useMemo(() =>
        messageHistory.current.concat(lastMessage), [lastMessage]);


    const connectionStatus = {
        [ReadyState.CONNECTING]: 'Connecting',
        [ReadyState.OPEN]: 'Open',
        [ReadyState.CLOSING]: 'Closing',
        [ReadyState.CLOSED]: 'Closed',
        [ReadyState.UNINSTANTIATED]: 'Uninstantiated',
    }[readyState];

    return [lastMessage, sendMessage, setSocketUrl]

}