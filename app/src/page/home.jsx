import React, {useCallback, useContext, useEffect, useReducer, useRef, useState} from "react";
import Thing from "../component/thing.jsx";
import Grid from "@material-ui/core/Grid";
import API from "../js/api";
import NewThingsDialog from "./new-things";

import ThingsReducer, {Actions} from "../js/things-reducer";
import TopBar from "../component/topBar";
import {useTranslation} from "react-i18next";
import {makeStyles} from "@material-ui/core/styles";
import {AppContext} from "../App";
import {drawerWidth} from "../js/constant";
import clsx from "clsx";
import {CircularProgress} from "@material-ui/core";
import useWebSocket, {ReadyState} from "react-use-websocket";
import AddCircleIcon from "@material-ui/icons/AddCircle";
import ErrorOutlineIcon from '@material-ui/icons/ErrorOutline';
import {ThingPanel} from "../component/thing-panel";

const useStyles = makeStyles((theme) => ({
    containerGrid: {
        alignItems: "flex-start",
        height: "200",
        minHeight: "200",
        // background: 'linear-gradient(45deg, #FE6B8B 30%, #FF8E53 90%)',
        flexGrow: 1,
        padding: theme.spacing(3),
        transition: theme.transitions.create('margin', {
            easing: theme.transitions.easing.sharp,
            duration: theme.transitions.duration.leavingScreen,
        }),
        marginLeft: drawerWidth,
    },
    contentShift: {
        transition: theme.transitions.create('margin', {
            easing: theme.transitions.easing.easeOut,
            duration: theme.transitions.duration.enteringScreen,
        }),
        marginLeft: 0,
    },
    drawerHeader: {
        display: 'flex',
        alignItems: 'center',
        padding: theme.spacing(0, 1),
        // necessary for content to be below app bar
        ...theme.mixins.toolbar,
        justifyContent: 'flex-end',
    },

}))

const states = {
    fetching: "fetching",
    completed: "completed",
    connected: "connected",
    disconnected: "disconnected",
}

export default function Home(props) {

    const classes = useStyles()
    const {drawerOpen} = useContext(AppContext)
    const [addThingShow, setAddThingShow] = useState(false)
    const [thingPanelShow, setThingPanelShow] = useState(false)
    const [currentThing, setCurrentThing] = useState()
    const [things, dispatch] = useReducer(ThingsReducer, new Map())
    const {t, i18n} = useTranslation();
    const url = "ws://localhost:9090/things/"
    const [socketUrl, setSocketUrl] = useState(null);
    const didUnmount = useRef(false);
    const {
        sendMessage,
        lastMessage,
        readyState,
    } = useWebSocket(socketUrl
    );
    const [state, setState] = useState()
    const connectionStatus = {
        [ReadyState.CONNECTING]: 'Connecting',
        [ReadyState.OPEN]: 'Open',
        [ReadyState.CLOSING]: 'Closing',
        [ReadyState.CLOSED]: 'Closed',
        [ReadyState.UNINSTANTIATED]: 'Uninstantiated',
    }[readyState];

    const handleSendMessage = useCallback((data) => {
            console.log("sendMessage:", data)
            sendMessage(JSON.stringify(data))
        },
        []);

    // resolve=> return Promise(map)
    function fetchThings() {
        return new Promise(
            function (resolve, reject) {
                API.getThings().then((data) => {
                    if (!data) {
                        reject(new Error("fetch things empty"))
                    }
                    let map = new Map()
                    for (const t of data) {
                        map.set(t.id, t)
                    }
                    if (map) {
                        return resolve(map)
                    }
                }).catch((e) => {
                    reject(e)
                })
            }
        )

    }

    useEffect(() => {
        if (connectionStatus === 'Open' && state === states.completed) {
            setState(states.connected)
        }
        if (connectionStatus === 'Closed') {
            setState(states.disconnected)
        }

    }, [connectionStatus])


    useEffect(() => {

        setState(states.fetching)
        fetchThings().then((mapThings) => {
            console.log("mapThings:", mapThings)
            if (mapThings) {
                dispatch({type: Actions.initial, initialState: mapThings})
            }
            setState(states.completed)
        }).catch((e) => {
                console.log(e, "------------------------------------------")
                setState(states.disconnected)
            }
        )
        return () => {
            didUnmount.current = true;
        };
    }, [])


    useEffect(() => {

        if (lastMessage != null) {
            console.log("websocket rev message:", lastMessage.data)
        }

    }, [lastMessage])

    useEffect(() => {
        if (state === states.completed) {
            try {
                setSocketUrl(url)
            } catch (e) {
                console.error(e, "++++++++++++++++++++++++")
            }
        }
        if (state === states.disconnected) {
            setSocketUrl("")
        }
    }, [state])


    function renderThings() {

        let list = []
        for (let [id, t] of things) {
            const thing =
                <Thing key={id} {...t} openPanel={handleOpenThingPanel} sendMessage={handleSendMessage}/>
            list.push(thing)
        }
        return list
    }

    function handleOpenThingPanel(props) {
        setThingPanelShow(true)
        setCurrentThing(props)
    }

    function handle(props) {

    }


    return (
        <>
            <TopBar add={true} show={setAddThingShow} title={t("Home")}/>
            <div className={classes.drawerHeader}/>


            <Grid style={{"justifyContent": !state === states.fetching || things.keys > 0 ? 'flex-start' : "center"}}
                  className={clsx(classes.containerGrid, {
                      [classes.contentShift]: !drawerOpen,
                  })}
                  container spacing={2}>
                {state === states.fetching && <CircularProgress disableShrink/>}
                {state === states.connected && things.length === 0 && <AddCircleIcon/>}
                {state === states.connected && things && renderThings()}
                {state === states.disconnected &&
                <div style={{display: "flex", flexDirection: "column", alignItems: "center"}}><ErrorOutlineIcon/>
                    <h4>{t("disconnect")}</h4></div>}
            </Grid>
            <NewThingsDialog open={addThingShow} show={setAddThingShow}/>
            <ThingPanel open={thingPanelShow} show={setThingPanelShow} {...currentThing}/>
        </>

    )
}


