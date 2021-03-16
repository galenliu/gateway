import React, {useContext, useEffect, useState} from "react";
import Grid from "@material-ui/core/Grid";
import {ErrorOutlined} from "@material-ui/icons";
import TopBar from "../component/topBar";
import {useTranslation} from "react-i18next";
import {makeStyles} from "@material-ui/core/styles";
import {AppContext} from "../Router";
import Constants, {drawerWidth} from "../js/constant";
import clsx from "clsx";
import {CircularProgress} from "@material-ui/core";
import App from "../App";
import {ThingPanel} from "../component/thing-panel";
import NewThingsDialog from "./AddThing";
import Thing from "../component/thing";


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
    pending: "pending",
    connected: "connected",
    disconnected: "disconnected",
}


export default function Things(props) {
    const classes = useStyles()
    const [things, setThings] = useState(props.things)
    const {drawerOpen} = useContext(AppContext)
    const [addThingShow, setAddThingShow] = useState(false)
    const [thingPanelShow, setThingPanelShow] = useState(false)
    const [currentThingId, setCurrentThingId] = useState(null)
    //const [things, dispatch] = useReducer(ThingsReducer, new Map())
    const [state, setState] = useState(states.connected)
    const {t, i18n} = useTranslation();

    useEffect(() => {
        const refreshThings = (things) => {
            console.log("=================", things)
            setThings(things)
        }
        App.gatewayModel.subscribe(Constants.REFRESH_THINGS, refreshThings)
        return () => {
            App.gatewayModel.unsubscribe(Constants.REFRESH_THINGS, refreshThings)
        }
    }, [])

    function renderThings() {
        if (things === undefined) {
            return
        }
        let thingList = []
        for (const [key, value] of things) {
            console.log("props:", props)
            let model = App.gatewayModel.thingModels.get(key)
            let thing = <Thing key={key} description={value} model={model} open={openPanel}/>
            thingList.push(thing)
            console.log(thing)
        }
        return thingList
    }

    function openPanel(thingId) {
        setCurrentThingId(thingId)
        setThingPanelShow(true)
    }

    return (
        <>
            <TopBar add={true} show={setAddThingShow} title={t("Things")}/>
            <div className={classes.drawerHeader}/>
            <Grid style={{"justifyContent": !state === states.pending ? 'flex-start' : "center"}}
                  className={clsx(classes.containerGrid, {
                      [classes.contentShift]: !drawerOpen,
                  })}
                  container spacing={2}>
                {state === states.pending && <CircularProgress disableShrink/>}
                {state === states.connected && renderThings()}
                {state === states.disconnected &&
                <div style={{display: "flex", flexDirection: "column", alignItems: "center"}}><ErrorOutlined/>
                    <h4>{t("disconnect")}</h4></div>}
            </Grid>
            <NewThingsDialog open={addThingShow} show={setAddThingShow}/>
            {currentThingId !== null && <ThingPanel open={thingPanelShow} show={setThingPanelShow} thingId={currentThingId}/>}
        </>

    )
}


