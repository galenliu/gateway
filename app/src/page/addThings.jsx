import React, {useCallback, useContext, useEffect, useRef, useState} from 'react';
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
    const {t, i18n} = useTranslation();

    const ws = useRef(null)

    const {open, setNewThingsClose, setNewThingsOpen} = useContext(HomeContext)
    const [availableThings, setAvailableThings] = useState({})
    const [actionUrl, setActionUrl] = useState()
    const [message, setMessage] = useState('');
    const [readyState, setReadyState] = useState('正在链接中');

    const webSocketInit = useCallback(() => {
        const stateArr = [
            '正在链接中',
            '已经链接并且可以通讯',
            '连接正在关闭',
            '连接已关闭或者没有链接成功',
        ];

        console.log("start requestPairing......")
        let proto = 'ws://';
        if (window.location.protocol === 'https:') {
            proto = 'wss://';
        }
        let host = window.location.host
        const path = proto + host + "/new_things"

        ws.current = new WebSocket(path);
        if (!ws.current || ws.current.readyState === 3) {

            ws.current.onopen = _e =>
                setReadyState(stateArr[ws.current?.readyState ?? 0]);
            ws.current.onclose = _e =>
                setReadyState(stateArr[ws.current?.readyState ?? 0]);
            ws.current.onerror = e =>
                setReadyState(stateArr[ws.current?.readyState ?? 0]);
            ws.current.onmessage = e => {
                setMessage(e.data);
            };
        }
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
                console.log("message:", message)
                // if (lastMessage !== null) {
                //     const things = availableThings
                //     if (!availableThings.hasOwnProperty(lastJsonMessage.id)) {
                //         things[lastJsonMessage.id] = lastJsonMessage
                //         setAvailableThings({...things})
                //         console.log("availableThings:", availableThings)
                //     }
                // }

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
            if (open) {
                console.log("....................")
                requestPairing()
            }
            if (!open) {
                cancelPairing()
            }
        }, [open]
    )

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
                // const newThing = <NewThing key={thing.id} thing={thing}
                //                            onSave={saveRequest}
                // />
                list.push(thing)
            }
        }
        return list
    }


    return (
        <div>{open &&
        <Dialog fullScreen open={open} onClose={setNewThingsOpen} TransitionComponent={Transition}>
            <AppBar className={classes.appBar}>
                <Toolbar>
                    <Typography variant="h6" className={classes.title}>
                        {t("AddNewThings")}......{readyState}
                    </Typography>
                    <IconButton autoFocus color="inherit" onClick={() => {
                        {
                            setNewThingsClose()
                            cancelPairing()
                        }
                    }} aria-label="close">
                        <CloseIcon/>
                    </IconButton>
                </Toolbar>
            </AppBar>


            {open && RenderAvailableThings()}
            {/*<ul>*/}
            {/*    {messageHistory.current*/}
            {/*        .map((message, idx) => <span key={idx}>{message.data}</span>)}*/}
            {/*</ul>*/}
            {/*{RenderAvailableThings()}*/}

        </Dialog>}
        </div>
    );
}