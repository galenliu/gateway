import React, {useEffect, useState} from "react";
import Grid from "@material-ui/core/Grid";
import Card from "@material-ui/core/Card";
import {makeStyles} from "@material-ui/core/styles";
import ThingIcon, {ActionsIcon} from "./thing-icon";
import Typography from "@material-ui/core/Typography";
import {useTranslation} from "react-i18next";
import {ThingPanel} from "./thing-panel";

// function rand() {
//     return Math.round(Math.random() * 20) - 10;
// }
//
// function getModalStyle() {
//     const top = 50 + rand();
//     const left = 50 + rand();
//     return {
//         top: `${top}%`,
//         left: `${left}%`,
//         transform: `translate(-${top}%, -${left}%)`,
//     };
// }

const useStyles = makeStyles((theme) => ({
    thingCard: {
        borderRadius: 12,
        display: 'flex',
        minWidth: 160,
        maxWidth: 160,
        justifyContent: "space-between",
        flexDirection: "column",
    },
    modal: {
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'center',
    },
    cardTop: {
        display: "flex",
        flexDirection: "row",
        marginTop: 5,
        marginRight: 5,
        alignItems: "flex-start",
        justifyContent: 'space-between',
    },
    cardBot: {
        display: "flex",
        marginLeft: 4,
        flexDirection: "column",
        alignItems: "flex-start",
        justifyContent: 'space-between',
    },
    root: {
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'center',
    },
}))

const states = {
    on: "on",
    off: "off",
    disconnect: "disconnect",
    updating: "updating"
}

export default function Thing(props) {

    const classes = useStyles()
    const {t, i18n} = useTranslation();
    const [open, setOpen] = useState(false)

    const [thing, setThing] = useState(props)
    const [state, setState] = useState()

    useEffect(() => {
        setState(states.updating)

    }, [])


    // const {handleSendMessage} = useContext(HomeContext)
    //
    //
    // const setPropertyValue = useCallback((property) => {
    //         console.log("set property", property)
    //         const data = {}
    //         console.log(data)
    //         let message = {}
    //         message.messageType = "setProperty"
    //         message.id = props.thing.id
    //         message.data = {}
    //         message.data[property.name] = property.value
    //         console.log("set property message:", message)
    //         handleSendMessage(message)
    //     }, []
    // )
    //
    // const propertyValueChanged = (thingId, value) => {
    //     console.log(thingId, value)
    // }
    //
    // const onClick = (e) => {
    //     console.log("on Click")
    // }
    function stateTextStyle() {
        if (state === states.updating) {
            return {color: "primary", variant: "body2"}
        }
        if (state === states.disconnect) {
            return "secondary"
        }
        return {color: "primary", variant: "h6"}
    }

    return (
        <>

            <Grid item className={classes.root} onDoubleClick={() => {
                console.log("doubleClick")
                setOpen(true)
            }}>

                <Card elevation={10} className={classes.thingCard}>
                    <div className={classes.cardTop}>
                        <ThingIcon state={state} color={"#fb8c00"} type={thing.selectedCapability} size={2}/>
                        <ActionsIcon state={state} type={thing.selectedCapability} size={2}/>
                    </div>
                    <div className={classes.cardBot}>
                        <Typography variant={"body1"}>
                            {thing.title}
                        </Typography>
                        <Typography {...stateTextStyle()} t={2}>
                            {t(state)}
                        </Typography>
                    </div>
                </Card>
            </Grid>
        </>)
}
