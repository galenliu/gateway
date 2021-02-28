import React, {useContext, useEffect, useState} from 'react';
import {makeStyles} from '@material-ui/core/styles';
import Dialog from '@material-ui/core/Dialog';
import AppBar from '@material-ui/core/AppBar';
import Toolbar from '@material-ui/core/Toolbar';
import IconButton from '@material-ui/core/IconButton';
import Typography from '@material-ui/core/Typography';
import CloseIcon from '@material-ui/icons/Close';
import Slide from '@material-ui/core/Slide';
import {HomeContext} from "./home"
import API from "../js/api";
import {useTranslation} from "react-i18next";
import Grid from "@material-ui/core/Grid";
import NewThing from "../component/new-thing";
import * as url from "url";

const useStyles = makeStyles((theme) => ({
    appBar: {
        position: 'relative',
    },
    title: {
        marginLeft: theme.spacing(2),
        flex: 1,
    },
}));

const Transition = React.forwardRef(function Transition(props, ref) {
    return <Slide direction="up" ref={ref} {...props} />;
});

export default function NewThingsDialog(props) {

    let url

    const {t, i18n} = useTranslation();

    const classes = useStyles();
    const {open, setNewThingsClose, setNewThingsOpen} = useContext(HomeContext)
    const [availableThings, setAvailableThings] = useState({})


    useEffect(
        () => {
            if (open) {
                requestPairing()
            } else {
                cancelPairing()
            }
            return () => cancelPairing(url)

        }, [open]
    )


    function requestPairing() {
        console.log("start pairing")
        let proto = 'ws://';
        if (window.location.protocol === 'https:') {
            proto = 'wss://';
        }
        let host = window.location.host
        const path = proto + host + "/new_things"
        let socket = new WebSocket(path)

        socket.onmessage = (e) => {
            console.log("e.data:",e.data)
            try {
                const thingObj = JSON.parse(e.data)
                if (thingObj !== null) {
                    const things = availableThings
                    if (!availableThings.hasOwnProperty(thingObj.id)) {
                        things[thingObj.id] = thingObj
                        setAvailableThings({...things})
                        console.log("availableThings:", availableThings)
                    }
                }

            } catch (e) {
                console.log(e)
            }
        }
        socket.onerror = (err) => {
            console.log("websocket err:", err)
            cancelPairing()
        }

        API.startPairing(5000).then((action) => {
            url = action.href
            setTimeout(() => {
                cancelPairing(url)
                socket.close()
            }, 5000)
        }).catch((err) => {
            console.log(err)
            cancelPairing()
        })
    }

    function cancelPairing(url) {
        API.cancelPairing(url).catch((err) => {
            console.log(err)
        })
    }

    function close(){
        cancelPairing()
    }

    function saveRequest(id, option) {
        if (availableThings.hasOwnProperty(id)) {
            let saveThing = {...availableThings[id], ...option}
            console.log(saveThing)
            API.addThing(saveThing).catch(e => {
                console.log(e)
            })

        }
    }

    function RenderAvailableThings() {
        let list = []
        for (let thingId in availableThings) {
            if (availableThings.hasOwnProperty(thingId)) {
                let thing = availableThings[thingId]
                const newThing = <NewThing key={thing.id} thing={thing}
                                           onSave={saveRequest}
                />
                list.push(newThing)
            }
        }
        return list
    }


    return (
        <div>
            <Dialog fullScreen open={open} onClose={setNewThingsOpen} TransitionComponent={Transition}>
                <AppBar className={classes.appBar}>
                    <Toolbar>
                        <Typography variant="h6" className={classes.title}>
                            {t("AddNewThings")}
                        </Typography>
                        <IconButton autoFocus color="inherit" onClick={()=>{setNewThingsClose()}} aria-label="close">
                            <CloseIcon/>
                        </IconButton>
                    </Toolbar>
                </AppBar>
                <Grid color="red">
                    {RenderAvailableThings()}
                </Grid>
            </Dialog>
        </div>
    );
}