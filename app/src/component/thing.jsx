import React, {useCallback, useContext, useEffect, useRef, useState} from "react";
import Light from "./capability/light";
import Grid from "@material-ui/core/Grid";
import Modal from '@material-ui/core/Modal';
import Card from "@material-ui/core/Card";
import {ThingsContext} from "../page/home"
import {makeStyles} from "@material-ui/core/styles";

function rand() {
    return Math.round(Math.random() * 20) - 10;
}

function getModalStyle() {
    const top = 50 + rand();
    const left = 50 + rand();
    return {
        top: `${top}%`,
        left: `${left}%`,
        transform: `translate(-${top}%, -${left}%)`,
    };
}

const useStyles = makeStyles((theme) => ({
    thingCard: {
        // background: 'linear-gradient(45deg, #FE6B8B 30%, #FF8E53 90%)',
        padding: 10,
        border: 2,
        borderRadius: 25,
        boxShadow: '0 3px 5px 2px rgba(0, 0, 0, .3)',
        color: 'black',
        animationName: "click",
    },
    modal: {
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'center',
    },
}))


export default function Thing(props) {

    const classes = useStyles()
    const [modalStyle] = React.useState(getModalStyle);

    const [icon, setIcon] = useState()
    const [state, setState] = useState()
    const [open, setOpen] = useState(false)

    const [thing,setThing] = useState()

    const thingRef = useRef()

    const updateIcon=(icon)=> {
        setIcon(icon)
    }

    const updateState=(state)=> {
        setState(state)
    }

    function handleOpen() {
        setOpen(true)
    }

    function handleClose() {
        setOpen(false)
    }


    function InitThing() {
        const data = props.thing
        let thing
        if (data.hasOwnProperty("selectedCapability")) {
            let selectedCapability = data["selectedCapability"]
            if (selectedCapability === "Light") {
                thing = <Light updataIcon={updateIcon} ref={thingRef}  updataState ={updateState}/>
            }
        }
        return thing
    }

    useEffect(() => {
        console.log("11111",InitThing())
        setThing(InitThing())
        console.log("thing---------:",thing)

    },[])




    const {handleSendMessage} = useContext(ThingsContext)


    const setPropertyValue = useCallback((property) => {
            console.log("set property", property)
            const data = {}
            console.log(data)
            let message = {}
            message.messageType = "setProperty"
            message.id = props.thing.id
            message.data = {}
            message.data[property.name] = property.value
            console.log("set property message:", message)
            handleSendMessage(message)
        }, []
    )

    const propertyValueChanged = (thingId, value) => {
        console.log(thingId, value)
    }

    const onClick = (e) => {
        console.log("on Click")
    }

    return <><ThingsContext.Provider value={{setPropertyValue}}>
        <Grid item xs={4} sm={2} id="thing-card">
            <Card elevation={10} className={classes.thingCard} onClick={onClick} onDoubleClick={handleOpen}>
                {icon}
                {state}
                <Modal open={open} onClose={handleClose}
                       className={classes.modal}
                       aria-labelledby="simple-modal-title"
                       aria-describedby="simple-modal-description">
                  <h1>sdfasdf</h1>
                </Modal>
            </Card>
        </Grid>
    </ThingsContext.Provider>
    </>

}
