import React, {useContext, useState} from "react";
import {makeStyles} from "@material-ui/core/styles";
import ListItemIcon from '@material-ui/core/ListItemIcon';
import ListItemText from '@material-ui/core/ListItemText';
import ListItem from '@material-ui/core/ListItem';
import List from "@material-ui/core/List";
import {useTranslation} from "react-i18next";
import DomainIcon from '@material-ui/icons/Domain';
import ExtensionIcon from '@material-ui/icons/Extension';
import TopBar from "../component/topBar";
import {CssBaseline, useTheme} from "@material-ui/core";
import clsx from "clsx";
import {AppContext} from "../App";
import Grid from "@material-ui/core/Grid";
import AddonsDialog from "./Addons";
import Divider from "@material-ui/core/Divider";
import NavigateNextIcon from '@material-ui/icons/NavigateNext';
import {drawerWidth} from "../js/constant";


const useStyles = makeStyles((theme) => ({
    root: {
        display: "flex",

    },
    drawerHeader: {
        display: 'flex',
        alignItems: 'center',
        padding: theme.spacing(0, 1),
        // necessary for content to be below app bar
        ...theme.mixins.toolbar,
        justifyContent: 'flex-end',
    },
    content: {
        flexGrow: 1,
        padding: theme.spacing(3),
        transition: theme.transitions.create('margin', {
            easing: theme.transitions.easing.sharp,
            duration: theme.transitions.duration.leavingScreen,
        }),
        marginLeft: drawerWidth,
    },
    contentShift: {
        transition: theme.transitions.create('margin', {
            easing: theme.transitions.easing.easeOut,
            duration: theme.transitions.duration.enteringScreen,
        }),
        marginLeft: 0,
    },

    listItem: {
        maxWidth: 400,
        minWidth: 300,
    }

}))


export default function Settings(props) {

    const classes = useStyles(props);
    const theme = useTheme();
    const {t, i18n} = useTranslation();
    const {drawerOpen, setNewThingsClose, setNewThingsOpen, newThingsOpen} = useContext(AppContext)
    const [addonsDialogShow, setAddonsDialogShow] = useState(false)

    return (
        <>
            <CssBaseline/>
            <TopBar add={false} title={t("Settings")}/>
            <main
                className={clsx(classes.content, {
                    [classes.contentShift]: !drawerOpen,
                })}
            >
                <Grid container justify="flex-start" alignItems="center" direction="column">
                    <div className={classes.drawerHeader}/>

                    <List component="nav" aria-label="main mailbox folders" className={classes.list}>
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
                        <ListItem button
                                  className={classes.listItem} variant="contained"
                                  onClick={() => setAddonsDialogShow(true)}>
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
                        <ListItem button
                                  className={classes.listItem} variant="contained" elevation={111}>
                            <ListItemIcon>
                                <DomainIcon/>
                            </ListItemIcon>
                            <ListItemText primary={t("Domain")}/>
                            <NavigateNextIcon/>
                        </ListItem>
                        <Divider/>
                    </List>
                </Grid>

            </main>
            <AddonsDialog open={addonsDialogShow} show={setAddonsDialogShow}/>
        </>
    )
}