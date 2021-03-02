import React, {useState} from "react";
import Card from "@material-ui/core/Card";
import {makeStyles} from "@material-ui/core/styles";
import ThingIcon from "./thing-icon";
import {Button, FormControl, FormHelperText, MenuItem, Select, TextField} from "@material-ui/core";
import {useTranslation} from "react-i18next";
import Divider from "@material-ui/core/Divider";


const useStyles = makeStyles((theme) => ({
    newThingCard: {
        display: "flex",
        maxWidth: 345,
        margin: 10,
        alignItems: "center",
    },

    formControl: {
        margin: theme.spacing(1),
        minWidth: 120,

    },
    button: {
        margin: theme.spacing(1),
        minHigh: 50,
    },

    icon: {},
    content: {
        padding: 12,
        display: "flex",
        flexDirection: "column"
    }
}))


export default function NewThing(props) {

    const {t, i18n} = useTranslation();

    const classes = useStyles();
    const [thing, setThing] = useState({...props["thing"], selectedCapability: props.thing["@type"][0]})

    console.log("selectedCapability:", thing.selectedCapability)

    // function setTitle(title) {
    //     setThings({...thing, title: title})
    // }

    // const [selectedCapability,setSelectedCapability] = useState(thing.selectedCapability)
    function capabilityItem() {
        let list = []
        for (const item of props.thing["@type"]) {
            list.push(<MenuItem value={item} key={item}>{t(item)}</MenuItem>)
        }
        return list
    }

    const handleChange = (event) => {
        let t = thing
        t.selectedCapability = event.target.value
        setThing({...t});
    };

    const handleTitleChange = (event) => {
        let t = thing
        t.title = event.target.value
        console.log("event.target.value", event.target.value)
        setThing({...t});
    };

    return <>
        <Card className={classes.newThingCard} elevation={5}>
            <ThingIcon className={classes.icon} style={{fontSize: 80}} color={"primary"}
                       type={thing.selectedCapability}/>
            <div className={classes.content}>
                <TextField id="standard-basic" label={t("Title")} defaultValue={thing.title}
                           onChange={handleTitleChange}/>
                <FormControl className={classes.formControl}>
                    <FormHelperText>{t("SelectedCapability")}</FormHelperText>
                    <Select
                        onChange={handleChange}
                        labelId="demo-simple-select-label"
                        id="demo-simple-select"
                        value={thing.selectedCapability}
                    >
                        {capabilityItem()}
                    </Select>
                </FormControl>

            </div>
            <Button
                variant="contained"
                color="primary"
                className={classes.button}
            >
                {t("Save")}
            </Button>

        </Card>
        <Divider/>

    </>
}
