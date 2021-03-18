import React, {useEffect, useState} from "react";
import Grid from "@material-ui/core/Grid";
import Card from "@material-ui/core/Card";
import {makeStyles} from "@material-ui/core/styles";
import Icons, {ActionsIcon} from "./icons";
import Typography from "@material-ui/core/Typography";
import {useTranslation} from "react-i18next";
import Constants, {ThingProperties, ThingType as Things} from "../js/constant";
import App from "../App";

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
    const [description, setDescription] = useState(App.gatewayModel.things.get(props.thingId))
    const [state, setState] = useState()

    useEffect(() => {
        function refreshThing(properties) {

        }

        let thingModel = App.gatewayModel.thingModels.get(props.thingId)
        thingModel.subscribe(Constants.PROPERTY_STATUS, refreshThing)
        setState(states.updating)
        return () => {
            thingModel.unsubscribe(Constants.PROPERTY_STATUS, refreshThing)
        }
    }, [])

    function stateTextStyle() {
        if (state === states.updating) {
            return {color: "primary", variant: "body2"}
        }
        if (state === states.disconnect) {
            return "secondary"
        }
        return {color: "primary", variant: "h6"}
    }


    function handleToggleClick(e) {
        e.stopPropagation()
        if (description.selectedCapability === Things.Light) {
            for (const name in model.propertyDescriptions) {

                let prop = model.propertyDescriptions[name]
                if (prop["@type"] === ThingProperties.OnOffProperty) {
                    let value = prop.value
                    let newValue
                    switch (prop.type) {
                        case 'number':
                            value = parseFloat(value);
                            break;
                        case 'integer':
                            value = parseInt(value);
                            break;
                        case 'boolean':
                            value = Boolean(value);
                            newValue = !value
                            break;
                    }
                    console.log("toggle:", newValue)
                }
            }
        }
    }

    return (
        <>
            <Grid item className={classes.root}>
                <Card elevation={10} className={classes.thingCard} onClick={() => props.open(props.thingId)}>
                    <div className={classes.cardTop}>
                        <Icons state={state} color={"#fb8c00"} type={description.selectedCapability} size={2}/>
                        <ActionsIcon cursor={"pointer"} state={state} type={description.selectedCapability} size={2}
                                     onClick={handleToggleClick}/>
                    </div>
                    <div className={classes.cardBot}>
                        <Typography variant={"body1"}>
                            {description.title}
                        </Typography>
                        <Typography {...stateTextStyle()} t={2}>
                            {t(state)}
                        </Typography>
                    </div>
                </Card>
            </Grid>
        </>)
}
