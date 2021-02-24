import React, {useEffect, useRef, useState} from "react";
import {useTranslation} from 'react-i18next'
import "../css/things.css"
import Utils from "../util"
import Card from '@material-ui/core/Card';
import {Button} from "@material-ui/core";
import {classes} from "istanbul-lib-coverage";
import {makeStyles} from "@material-ui/core/styles";

const useStyles = makeStyles((theme) => ({
    root: {
        flexGrow: 1,
    },
    paper: {
        padding: theme.spacing(2),
        textAlign: 'center',
        color: theme.palette.text.secondary,
    },
}));


export default function NewThing(props) {

    const classes = useStyles();

    let {t, i18n} = useTranslation()

    const fileInput = useRef()
    const [iconFile, setIconFile] = useState()

    const [selectedCapability, setSelectedCapability] = useState(props.thing["@type"][0])

    const [title, setTitle] = useState(props.thing.title)

    function typeOption(typeList) {
        let options = []
        for (let i of typeList) {
            const option = <option value={i} key={i} selected={options.length === 0}>
                {t(i)}
            </option>
            options.push(option)
        }
        return options
    }

    useEffect(()=>{

       // setSelectedCapability(props["@type"][0])
        console.log(selectedCapability)

    })

    return <Card className={classes.root} id={"new-thing " + Utils.getClassFromCapability(selectedCapability)}>
        <div className=" new-thing-metadata">
            <input type=" text" name="thing-title" readOnly={false} onChange={(e) => setTitle(e.target.value)}
                   className=" new-thing-title" value={props.thing.title}/>

            <span className=" new-thing-spacer"/>

            <select className=" new-thing-type" onChange={(e) => {
                setSelectedCapability(e.target.value)
            }}>
                {typeOption(Utils.sortCapabilities(props.thing["@type"]).concat("Custom"))}
            </select>

            <span className=" new-thing-spacer"/>

            {selectedCapability === "Custom" &&
            <input type="file" className="new-thing-custom-icon-input" ref={fileInput}
                   accept="image/jpeg,image/png,image/svg+xml" onChange={(e) => {
                setIconFile(e.target.value)
            }}/>}

            {selectedCapability === "Custom" && <label className=" new-thing-custom-icon-label text-button"
                                                       onClick={() => fileInput.current.click()}>{t('chooseIcon')}</label>}
            {selectedCapability === "Custom" && <span className=" new-thing-label">{iconFile}</span>}

        </div>
        <Button color="primary" variant="contained" onClick={() => {
                props.addNewThingRequest({...props.thing, title: title, "selectedCapability": selectedCapability, iconFile: iconFile})
        }}>  {t('Save')}
        </Button>
    </Card>

}

