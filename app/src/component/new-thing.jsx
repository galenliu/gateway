import React, {useState} from "react";
import Card from "@material-ui/core/Card";
import {makeStyles} from "@material-ui/core/styles";
import ThingIcon from "./thingIcon";
import Typography from "@material-ui/core/Typography";
import Toolbar from "@material-ui/core/Toolbar";


const useStyles = makeStyles((theme) => ({
    newThingCard: {
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


export default function NewThing(props) {
    const classes = useStyles();

    const [thing, setThings] = useState({...props.thing, selectedCapability: props.thing.type[0]})



    function setTitle(title) {
        setThings({...thing, title: title})
    }

    const [selectedCapability,setSelectedCapability] = useState(thing.selectedCapability)



    return <>
        <Card classes={classes.newThingCard}>
            <ThingIcon type={thing.selectedCapability}/>
            <Typography variant="h6">

            </Typography>
        </Card>

    </>
}
