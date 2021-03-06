import React from 'react';
import {makeStyles} from '@material-ui/core/styles';
import Dialog from '@material-ui/core/Dialog';
import IconButton from '@material-ui/core/IconButton';
import CloseIcon from '@material-ui/icons/Close';
import Typography from '@material-ui/core/Typography';
import {useTranslation} from "react-i18next";
import ListItem from '@material-ui/core/ListItem';
import List from "@material-ui/core/List";
import {ListItemSecondaryAction, ListSubheader, Switch, useTheme} from "@material-ui/core";
import ListItemText from '@material-ui/core/ListItemText';
import ListItemIcon from '@material-ui/core/ListItemIcon';
import DialogContent from '@material-ui/core/DialogContent';

import ThingIcon from "./thing-icon";
import Grid from "@material-ui/core/Grid";
import AppBar from "@material-ui/core/AppBar";
import Toolbar from "@material-ui/core/Toolbar";
import StoreIcon from "@material-ui/icons/Store";
import {drawerWidth} from "../js/constant";
import {Transition} from "mdi-material-ui";
import Divider from "@material-ui/core/Divider";
import DomainIcon from "@material-ui/icons/Domain";
import NavigateNextIcon from "@material-ui/icons/NavigateNext";
import ExtensionIcon from "@material-ui/icons/Extension";

const styles = (theme) => ({
    root: {
        margin: 0,
        padding: theme.spacing(2),
    },

    closeButton: {
        position: 'absolute',
        right: theme.spacing(1),
        top: theme.spacing(1),
        color: theme.palette.grey[500],
    },
});

const useStyles = makeStyles((theme) => ({
    margin: {
        margin: theme.spacing(1),
    },
    list: {
        display: "flex",
        width: '100%',
        maxWidth: 360,
        flexDirection: "row",
        backgroundColor: theme.palette.background.paper,

    },
    content: {
        flexDirection: "column",
        display: "flex",
        alignItems: "flex-start",
        justifyContent: 'center',
        flexGrow: 1,
        padding: theme.spacing(3),
        transition: theme.transitions.create('margin', {
            easing: theme.transitions.easing.sharp,
            duration: theme.transitions.duration.leavingScreen,
        }),
        marginLeft: drawerWidth,
    },

    modal: {
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'center',
    },

    contentShift: {
        transition: theme.transitions.create('margin', {
            easing: theme.transitions.easing.easeOut,
            duration: theme.transitions.duration.enteringScreen,
        }),
        marginLeft: 0,
    },
    drawerHeader: {
        display: 'flex',
        alignItems: 'center',
        padding: theme.spacing(0, 1),
        // necessary for content to be below app bar
        ...theme.mixins.toolbar,
        justifyContent: 'flex-end',
    },
    appBar: {},
    title: {
        marginLeft: theme.spacing(2),
        flex: 1,
    },

    listItem: {
        maxWidth: 400,
        minWidth: 300,
    },
    root: {
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'center',
    },
}));


export function ThingPanel(props) {

    const {t, i18n} = useTranslation();
    const classes = useStyles();
    return (
        <DialogContent>
            <Dialog fullScreen className={classes.root} open={props.open} onClose={() => props.show(false)}
                    TransitionComponent={Transition}>
                <AppBar className={classes.appBar}>
                    <Toolbar>
                        <StoreIcon/>
                        <Typography variant="h6" className={classes.title}>
                            {props.title}
                        </Typography>
                        <IconButton autoFocus color="inherit" onClick={() => props.show(false)} aria-label="close">
                            <CloseIcon/>
                        </IconButton>
                    </Toolbar>
                </AppBar>
                <div className={classes.drawerHeader}/>
                <Grid className={classes.content} container justify="flex-start" alignItems="center" direction="column">
                    <DetailsPanel {...props}/>
                </Grid>
            </Dialog></DialogContent>

    );
}


export function ControlPanel(props) {

    const {t, i18n} = useTranslation();
    const theme = useTheme();


    return (
        <></>
    )
}

export function DetailsPanel(props) {

    const classes = useStyles();
    const {t, i18n} = useTranslation();
    const theme = useTheme();

    return (
        <>
            <List subheader={<ListSubheader>Settings</ListSubheader>} className={classes.list}>

                <Divider/>
                <ListItem button
                          className={classes.listItem} variant="contained" elevation={111}>
                    <ListItemIcon>
                        <ThingIcon  {...props} color={"#fb8c00"} type={props.selectedCapability}
                                    size={1}/>
                    </ListItemIcon>
                    <ListItemText primary={t("Domain")}/>
                    <NavigateNextIcon/>
                </ListItem>
                <Divider/>
                <ListItem button
                          className={classes.listItem} variant="contained"
                >
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
                    <ListItemText primary={t("Domain")}/>
                    <NavigateNextIcon/>
                </ListItem>
                <Divider/>

                <ListItem>

                    <ListItemText id="switch-list-label-wifi" primary="Wi-Fi"/>
                    <ListItemSecondaryAction>
                        <Switch
                            edge="end"
                            inputProps={{'aria-labelledby': 'switch-list-label-wifi'}}
                        />
                    </ListItemSecondaryAction>
                </ListItem>
            </List>

        </>
    )
}

