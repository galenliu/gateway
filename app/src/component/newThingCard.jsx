import React, {useCallback, useContext, useEffect, useRef, useState} from "react";
import Grid from "@material-ui/core/Grid";
import Modal from '@material-ui/core/Modal';
import Card from "@material-ui/core/Card";
import {ThingsContext} from "../page/home"
import { Camera, Settings } from 'mdi-material-ui/light'
import {makeStyles} from "@material-ui/core/styles";



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


export default function NewThingCard(props) {
    const classes = useStyles();

    const [selectedCapability, setSelectedCapability] = useState(props.type[0])
    const [title, setTitle] = useState(props.title)




    return <>
    <Card classes={classes.newThingCard}>
  <div >{}</div>
    </Card>

    </>
}
