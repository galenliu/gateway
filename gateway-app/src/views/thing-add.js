import React, {useEffect, useRef, useState} from "react";
import "../css/things.css";
import API from "../api";
import NewThing from "../components/new-thing";
import {Link} from "react-router-dom";
import {makeStyles} from "@material-ui/core/styles";
import useWebSocket, {ReadyState} from "react-use-websocket";
import Grid from "@material-ui/core/Grid";
import {Avatar, Paper} from "@material-ui/core";



const useStyles = makeStyles((theme) => ({
    root: {
        flexGrow: 1,
        overflow: 'auto',
        padding: theme.spacing(0, 3),
    },
    paper: {
        maxWidth: 800,
        margin: `${theme.spacing(1)}px auto`,
        padding: theme.spacing(2),
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
        try {
            if (thing) {
                API.addThing(thing).catch(e => {
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
                const newThing = <Paper className={classes.paper}>
                    <Grid container wrap="nowrap" spacing={2}>
                        <Grid item>
                            <Avatar>W</Avatar>
                        </Grid>
                        <Grid item xs>
                    <NewThing key={thing.id} thing={thing}
                              addNewThingRequest={addThingRequest}/>
                        </Grid>

                    </Grid>
                </Paper>
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
            <div className={classes.root}>

            <Grid container wrap="nowrap"  alignItems="center" justify="center" direction="column" spacing={1}>
                {RenderAvailableThings()}
            </Grid>

            </div>
        </>
    )
}



