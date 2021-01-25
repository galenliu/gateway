import React, {useEffect, useRef, useState} from "react";
import "../css/things.css";
import API from "../api";
import NewThing from "../components/new-thing";
import {Link} from "react-router-dom";
import Container from '@material-ui/core/Container';
import {makeStyles} from "@material-ui/core/styles";
import useWebSocket, {ReadyState} from "react-use-websocket";


const useStyles = makeStyles((theme) => ({
    root: {
        flexGrow: 1,
    },
    paper: {
        padding: theme.spacing(2),
        textAlign: 'center',
        color: theme.palette.text.secondary,
    },
}));


export default function ThingsAdd() {

    const classes = useStyles();

    const [availableThings, setAvailableThings] = useState([])

    const [socketUrl, setSocketUrl] = useState("ws://localhost:9090/new_things")
    const didUnmount = useRef(false);
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


    useEffect(() => {
            let proto = 'ws:';
            if (window.location.protocol === 'https:') {
                proto = 'wss:';
            }
            let host = window.location.host
            setSocketUrl(`ws://localhost:9090/new_things`);
            requestPairing()
            return () => {
                cancelPairing()
            }
        }, []
    )

    useEffect(() => {
            if (lastMessage != null) {
                try {
                    if (lastMessage.data) {
                        setAvailableThings([...availableThings])
                        const newThing = JSON.parse(lastMessage.data)
                        console.log("new thing:", newThing)
                        if (newThing.hasOwnProperty("id")) {
                            for (let t of availableThings) {
                                if (t.id === newThing.id) {
                                    return
                                }
                            }

                            let newList = []
                            newList = availableThings
                            newList.push(newThing)
                            setAvailableThings([...newList])
                            console.log("newList:", newList)
                            console.log("availableThings:", availableThings)
                        }
                    }
                } catch (e) {
                    console.log(e)
                }
            }
        }, [lastMessage]
    )


    function requestPairing() {
        API.startPairing(5000).then((json) => {
            console.log("action:", json)
        }).catch((err) => {
            console.log(err)
        })
    }

    function cancelPairing() {

        didUnmount.current = true;
        API.cancelPairing().then((json) => {
            console.log("action:", json)
        }).catch((err) => {
            console.log(err)
        })
    }

    function addThingRequest(thing) {
        console.log("addThingRequestaddThingRequestaddThingRequestaddThingRequestaddThingRequest")
        try {
            if(thing){
                API.addThing(thing).catch(e =>{
                    console.log(e)
                })
            }
        } catch (e) {
            console.log(e)
        }

    }

    function RenderAvailableThings() {
        let list = []
        for (let thing of availableThings) {
            if (thing.hasOwnProperty("id")) {
                console.log("render thing", thing)
                const newThing = <NewThing key={thing.id}
                                           thing={thing}
                                           addNewThingRequest={addThingRequest}
                />
                list.push(newThing)
            }
        }
        return list
    }


    return (
        <>
            <Link to="/things">
                <button id="back-button" className="icon-button"/>
            </Link>
            <div id="thing-add-view" className={classes.root}>
                <Container className="thing-add-view-container" maxWidth="sm">
                    {RenderAvailableThings()}
                </Container>
            </div>
        </>
    )
}



