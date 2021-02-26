import React, {useContext, useEffect, useState} from 'react';
import {makeStyles} from '@material-ui/core/styles';
import Button from '@material-ui/core/Button';
import Dialog from '@material-ui/core/Dialog';
import ListItemText from '@material-ui/core/ListItemText';
import ListItem from '@material-ui/core/ListItem';
import List from '@material-ui/core/List';
import Divider from '@material-ui/core/Divider';
import AppBar from '@material-ui/core/AppBar';
import Toolbar from '@material-ui/core/Toolbar';
import IconButton from '@material-ui/core/IconButton';
import Typography from '@material-ui/core/Typography';
import CloseIcon from '@material-ui/icons/Close';
import Slide from '@material-ui/core/Slide';
import {HomeContext} from "./home"
import API, {getWsUrl} from "../js/api";

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

        }, [open]
    )


    function requestPairing() {
        console.log("start pairing")
        let proto = 'ws://';
        if (window.location.protocol === 'https:') {
            proto = 'wss://';
        }
        let host = window.location.host
        const path =proto+host+"/new_things"
        console.log("ws new thing path:",path)
        let socket = new WebSocket(path)
        socket.onmessage = (e) => {
            console.log(e.data)
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
        }

        API.startPairing(5000).then((json) => {
            console.log("action:", json)
        }).catch((err) => {
            console.log(err)
        })
    }

    function cancelPairing() {
        // API.cancelPairing().then((json) => {
        //     console.log("action:", json)
        // }).catch((err) => {
        //     console.log(err)
        // })
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
                const newThing = <NewThing key={thing.id}
                                           id={thing.id}
                                           type={thing["@type"]}
                                           title={thing.title}
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
                        <IconButton edge="start" color="inherit" onClick={setNewThingsClose} aria-label="close">
                            <CloseIcon/>
                        </IconButton>
                        <Typography variant="h6" className={classes.title}>
                            Sound
                        </Typography>
                        <Button autoFocus color="inherit" onClick={setNewThingsClose}>
                            save
                        </Button>
                    </Toolbar>
                </AppBar>
                <List>
                    <ListItem button>
                        <ListItemText primary="Phone ringtone" secondary="Titania"/>
                    </ListItem>
                    <Divider/>
                    <ListItem button>
                        <ListItemText primary="Default notification ringtone" secondary="Tethys"/>
                    </ListItem>
                </List>
            </Dialog>
        </div>
    );
}