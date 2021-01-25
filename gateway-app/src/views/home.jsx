import React, {useCallback, useEffect, useReducer, useRef, useState} from "react";
import "../css/things.css";
import API from "../api";
import "../components/new-thing"
import {makeStyles} from '@material-ui/core/styles';
import {Link} from "react-router-dom";
import useWebSocket, {ReadyState} from "react-use-websocket";
import thingsReducer from "../components/use-things";
import Thing from "../components/thing";
import Grid from "@material-ui/core/Grid";
import {spacing} from "@material-ui/system";

const useStyles = makeStyles((theme) => ({
    root: {
        flexGrow: 1,
    },
    paper: {
        padding: theme.spacing(4),
        textAlign: 'center',
        color: theme.palette.text.secondary,
    },
}));

export const ThingContext = React.createContext()


export default function Home() {

    const [thingsContext, setThingsContext] = useReducer(thingsReducer, [])

    const [socketUrl, setSocketUrl] = useState("ws://localhost:9090/things/");

    const didUnmount = useRef(false);
    const classes = useStyles();

    const {
        sendMessage,
        lastMessage,
        readyState,
    } = useWebSocket(socketUrl,
        {
            shouldReconnect: (closeEvent) => {
                /*
                  useWebSocket will handle unmounting for you, but this is an example of a
                  case in which you would not want it to automatically reconnect
                */
                return didUnmount.current === false;
            },
            reconnectAttempts: 10,
            reconnectInterval: 3000,
        }
    );

    const connectionStatus = {
        [ReadyState.CONNECTING]: 'Connecting',
        [ReadyState.OPEN]: 'Open',
        [ReadyState.CLOSING]: 'Closing',
        [ReadyState.CLOSED]: 'Closed',
        [ReadyState.UNINSTANTIATED]: 'Uninstantiated',
    }[readyState];

    const handleSendMessage = useCallback((data) =>
            sendMessage(JSON.stringify(data)),
        []);


    useEffect(() => {

             setThingsContext([...thingsContext])
            API.getThings().then((getThings) => {
                if (getThings === null) {
                    return
                }
                console.log("getThings:", getThings)

                getThings.forEach(thing => {
                    setThingsContext({messageType: "addThing", thing: thing})
                })
            }).catch(e => {
                console.log("GET /things err:", e)
            })
            return () => {
                didUnmount.current = true;
            };
        }, []
    )

    useEffect(() => {
            if (lastMessage !== null) {
                console.log("lastMessage:",lastMessage)
                if (lastMessage.data) {
                    let message = JSON.parse(lastMessage.data)
                    if (message.messageType === "propertyStatus" && message.hasOwnProperty("id") && message.hasOwnProperty("data")) {

                    }
                }
            }
        }, [lastMessage]
    )


    function RenderThingsView() {
        let list = []
        for (let t of thingsContext) {
            console.log("item:----",t)
            const thing =
                <Thing key={t.id} thing={t}/>
            list.push(thing)
        }
        return list
    }


    return (
        <>
            {/*{addonThingShow && <ul id="new-things"> {RenderAvailableThings()}</ul>}*/}
            <Link to="/things/add">
                <button id="discover-addons-button"
                        className="icon-button right-top-button"/>
            </Link>
            <ThingContext.Provider value={{handleSendMessage}}>
                <div id="things-view" className={classes.root}>
                    <Grid container spacing={1}>
                        {RenderThingsView()}
                    </Grid></div>
            </ThingContext.Provider>
        </>
    )
};



