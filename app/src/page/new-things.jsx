import React, {useCallback, useContext, useEffect, useRef, useState} from 'react';
import {makeStyles} from '@material-ui/core/styles';
import Dialog from '@material-ui/core/Dialog';
import AppBar from '@material-ui/core/AppBar';
import Toolbar from '@material-ui/core/Toolbar';
import IconButton from '@material-ui/core/IconButton';
import Typography from '@material-ui/core/Typography';
import CloseIcon from '@material-ui/icons/Close';
import Slide from '@material-ui/core/Slide';
import API from "../js/api";
import {useTranslation} from "react-i18next";
import NewThing from "../component/new-thing";
import Grid from "@material-ui/core/Grid";
import {AppContext} from "../App";

const useStyles = makeStyles((theme) => ({
    appBar: {

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
    const {t, i18n} = useTranslation();

    const ws = useRef(null)

    const {drawerOpen, setDrawerOpen, newThingsOpen, setNewThingsOpen} = useContext(AppContext)


    const [availableThings, setAvailableThings] = useState({})
    const [actionUrl, setActionUrl] = useState()
    const [message, setMessage] = useState('');
    const [readyState, setReadyState] = useState('正在链接中');

    const webSocketInit = useCallback(() => {
        const stateArr = [
            '正在链接中',
            '正在扫描新设备',
            '连接正在关闭',
            '扫描完成',
        ];

        console.log("start requestPairing......")
        let proto = 'ws://';
        if (window.location.protocol === 'https:') {
            proto = 'wss://';
        }
        let host = window.location.host
        const path = proto + host + "/new_things"

        ws.current = new WebSocket(path);
        ws.current.onopen = _e =>
            setReadyState(stateArr[ws.current?.readyState ?? 0]);
        ws.current.onclose = _e =>
            setReadyState(stateArr[ws.current?.readyState ?? 0]);
        ws.current.onerror = e =>
            setReadyState(stateArr[ws.current?.readyState ?? 0]);
        ws.current.onmessage = e => {
            console.log(e)
            setMessage(e.data);
        };

    }, [ws]);


    console.log("addThings websocket starting ...,open:", open)

    function requestPairing() {
        webSocketInit();
        API.startPairing(5000).then((action) => {
            setActionUrl(action.href)
            setTimeout(() => {
                cancelPairing()
            }, 5000)
        }).catch((err) => {
            console.log("startPairing err:", err)
        })
    }

    useEffect(
        () => {
            try {
                if (message !== undefined) {
                    let newThing = JSON.parse(message)
                    const things = availableThings
                    if (!availableThings.hasOwnProperty(newThing.id)) {
                        things[newThing.id] = newThing
                        setAvailableThings({...things})
                        console.log("AvailableThings update :", availableThings)
                    }
                }
            } catch (e) {
                cancelPairing()
                console.log("message err:", e)
            }

        }, [message]
    )


    function cancelPairing() {
        console.log("cancel pairing....")
        ws.current?.close();
        if (actionUrl !== undefined) {
            API.cancelPairing(actionUrl).catch((err) => {
                console.log("cancelParing err:", err)
            })
        }
    }


    useEffect(
        () => {
            if (newThingsOpen) {
                console.log("....................")
                setAvailableThings([])
                requestPairing()
            }
            if (!newThingsOpen) {
                cancelPairing()
            }
        }, [newThingsOpen]
    )




    function RenderAvailableThings() {
        let list = []
        for (let thingId in availableThings) {
            let thing = availableThings[thingId]
            console.log("render thing :", thing)
            const newThing = <NewThing key={thing.id} thing={thing}

            />
            list.push(newThing)
        }
        return list
    }


    return (
        <div>
            <Dialog fullScreen open={newThingsOpen} onClose={() => setNewThingsOpen(true)}
                    TransitionComponent={Transition}>
                <AppBar className={classes.appBar}>
                    <Toolbar>
                        <Typography variant="h6" className={classes.title}>
                            {t(readyState)}......
                        </Typography>
                        <IconButton autoFocus color="inherit" onClick={() => {
                            {
                                setNewThingsOpen(false)
                                cancelPairing()
                            }
                        }} aria-label="close">
                            <CloseIcon/>
                        </IconButton>
                    </Toolbar>
                </AppBar>
                <Grid container justify="flex-start" alignItems="center" direction="column">
                    {RenderAvailableThings()}
                </Grid>
            </Dialog>
        </div>
    );
}