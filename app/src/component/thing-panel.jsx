import React, {useCallback, useEffect, useState} from 'react';
import {makeStyles} from '@material-ui/core/styles';
import Dialog from '@material-ui/core/Dialog';
import DialogActions from '@material-ui/core/DialogActions';
import IconButton from '@material-ui/core/IconButton';
import DialogTitle from '@material-ui/core/DialogTitle';
import CloseIcon from '@material-ui/icons/Close';
import Typography from '@material-ui/core/Typography';
import {useTranslation} from "react-i18next";
import ListItem from '@material-ui/core/ListItem';
import List from "@material-ui/core/List";
import {
    Button,
    DialogContentText,
    FormControl,
    InputLabel,
    ListItemSecondaryAction,
    ListSubheader,
    MenuItem,
    Select,
    Switch,
    TextField,
} from "@material-ui/core";
import ListItemText from '@material-ui/core/ListItemText';
import ListItemIcon from '@material-ui/core/ListItemIcon';
import DialogContent from '@material-ui/core/DialogContent';
import Icons from "./icons";
import Grid from "@material-ui/core/Grid";
import AppBar from "@material-ui/core/AppBar";
import Toolbar from "@material-ui/core/Toolbar";
import Divider from "@material-ui/core/Divider";
import DomainIcon from "@material-ui/icons/Domain";
import NavigateNextIcon from "@material-ui/icons/NavigateNext";
import ExtensionIcon from "@material-ui/icons/Extension";
import Slide from '@material-ui/core/Slide';
import CheckCircleIcon from '@material-ui/icons/CheckCircle';
import App from "../App";
import Constants from "../js/constant";

const useStyles = makeStyles((theme) => ({
    margin: {
        margin: theme.spacing(1),
    },
    list: {
        display: "flex",
        width: '100%',
        flexDirection: "column",
        backgroundColor: theme.palette.background.paper,
    },
    content: {
        flexDirection: "column",
        display: "flex",
        alignItems: "center",
        justifyContent: 'center',
        flexGrow: 1,
        padding: theme.spacing(3),
        transition: theme.transitions.create('margin', {
            easing: theme.transitions.easing.sharp,
            duration: theme.transitions.duration.leavingScreen,
        }),
    },
    drawerHeader: {
        display: 'flex',
        alignItems: 'center',
        padding: theme.spacing(0, 1),
        // necessary for content to be below app bar
        ...theme.mixins.toolbar,
        justifyContent: 'flex-end',
    },
    appBar: {
        position: 'relative',
    },
    title: {
        marginLeft: theme.spacing(2),
        flex: 1,
    },

    listItem: {
        marginTop: 5,
        marginBottom: 5,
        maxWidth: 400,
        minWidth: 300,
    },
    root: {
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'center',
    },
}));

const Transition = React.forwardRef(function Transition(props, ref) {
    return <Slide direction="up" ref={ref} {...props} />;
});

export function ThingPanel(props) {
    const {t, i18n} = useTranslation();
    const [thing, setThing] = useState(App.gatewayModel.things.get(props.thingId))
    const [properties, setProperties] = useState(thing.properties)
    const [removeDialogOpen, setRemoveDialog] = useState(false)
    const classes = useStyles();

    const removeThing = useCallback(() => {

            props.show(false)
        }
    )


    useEffect(() => {
            console.log("properties:", properties)
            const onRemoved = (things) => {
                if (!things.get(props.thingId))
                    setRemoveDialog(true)
            }
            const refreshThing = () => {
                setThing(App.gatewayModel.things.get(props.thingId))
            }
            App.gatewayModel.thingModels.get(props.thingId).subscribe(Constants.DELETE_THING, onRemoved)
            App.gatewayModel.thingModels.get(props.thingId).subscribe(Constants.REFRESH_THINGS, refreshThing)

            return () => {
                App.gatewayModel.thingModels.get(props.thingId).unsubscribe(Constants.DELETE_THING, onRemoved)
                App.gatewayModel.thingModels.get(props.thingId).unsubscribe(Constants.REFRESH_THINGS, refreshThing)
            }
        })

    useEffect(() => {
            setProperties(thing.properties)
        }, [thing])

    return (
        <>
            <Dialog
                open={removeDialogOpen}
                onClose={() => setRemoveDialog(false)}
                aria-labelledby="alert-dialog-title"
                aria-describedby="alert-dialog-description">
                <DialogTitle id="alert-dialog-title">{t("accessories removed")}</DialogTitle>
                <DialogContent>
                    <DialogContentText id="alert-dialog-description">
                        {t("accessories removed")}
                    </DialogContentText>
                </DialogContent>
                <DialogActions>
                    <Button onClick={() => {
                        setRemoveDialog(false)
                    }} color="primary" autoFocus>
                        {t("Ok")}
                    </Button>
                </DialogActions>
            </Dialog>
            {!removeDialogOpen &&
            <Dialog fullScreen className={classes.root} open={props.open} onClose={() => props.show(false)}
                    TransitionComponent={Transition}>

                <AppBar className={classes.appBar}>
                    <Toolbar>
                        <Icons prop color={"#fb8c00"} type={thing.selectedCapability}
                               size={2}/>
                        <Typography variant="h6" className={classes.title}>
                            {thing.title}
                        </Typography>
                        <IconButton autoFocus color="inherit" onClick={() => props.show(false)} aria-label="close">
                            <CloseIcon/>
                        </IconButton>
                    </Toolbar>
                </AppBar>
                <DialogContent>
                    <div className={classes.drawerHeader}/>
                    <Grid className={classes.content} container>
                        <ControlPanel thing={thing}/>
                        <DetailsPanel thing={thing} remove={removeThing}/>
                    </Grid>
                </DialogContent>
            </Dialog>}
        </>

    );
}

export function DetailsPanel(props) {

    const classes = useStyles();
    const {t, i18n} = useTranslation();
    const [title, setTitle] = useState(props.thing.title)

    function update() {

    }

    function renderProperty() {

    }

    return (
        <>
            <List subheader={<ListSubheader>Settings</ListSubheader>} className={classes.list}>
                <Divider/>
                <ListItem
                    className={classes.listItem} variant="contained" elevation={90}>
                    {renderProperty()}
                </ListItem>
                <ListItem
                    className={classes.listItem} variant="contained" elevation={90}>
                    <ListItemIcon>
                        <Icons prop color={"#fb8c00"} type={props.thing.selectedCapability}
                               size={1}/>
                    </ListItemIcon>
                    <TextField defaultValue={props.thing.title}
                               onChange={(e) => setTitle(e.target.value)}/>
                    {title !== props.thing.title &&
                    <CheckCircleIcon cursor={"pointer"}/>}
                </ListItem>

                <Divider/>
                <ListItem button className={classes.listItem} variant="contained">
                    <ListItemIcon>
                        <ExtensionIcon/>
                    </ListItemIcon>
                    <ListItemText primary={t("Addons")}/>
                    <NavigateNextIcon/>
                </ListItem>
                <Divider/>
                <ListItem button
                          className={classes.listItem} variant="contained" elevation={111}>
                    <ListItemIcon>
                        <DomainIcon/>
                    </ListItemIcon>
                    <FormControl className={classes.formControl} style={{width: "100%"}}>
                        <InputLabel id="demo-simple-select-label">{t("Room")}</InputLabel>
                        <Select
                            labelId="demo-simple-select-label"
                            id="demo-simple-select">
                            <MenuItem value={10}>客厅</MenuItem>
                            <MenuItem value={20}>卧室</MenuItem>
                            <MenuItem value={30}>厨房</MenuItem>
                        </Select>
                    </FormControl>

                </ListItem>
                <Divider/>
                <ListItem>
                    <ListItemText id="switch-list-label-wifi" primary={t("On")}/>
                    <ListItemSecondaryAction>
                        <Switch
                            edge="end"
                            inputProps={{'aria-labelledby': 'switch-list-label-wifi'}}
                        />
                    </ListItemSecondaryAction>
                </ListItem>
                <Divider/>
                <ListItem color={"red"} button onClick={() => props.remove()}>
                    <Button variant="contained" color="secondary" style={{width: "100%"}}>
                        {t("remove the accessories")}
                    </Button>
                </ListItem>
                <Divider/>
            </List>

        </>
    )
}

export function ControlPanel(props) {
    const {t, i18n} = useTranslation();

    return (
        <></>
    )
}

