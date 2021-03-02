import React, {useContext} from 'react';
import {makeStyles, useTheme} from '@material-ui/core/styles';
import IconButton from '@material-ui/core/IconButton';
import MailIcon from '@material-ui/icons/Mail';
import Drawer from "@material-ui/core/Drawer";
import ChevronLeftIcon from "@material-ui/icons/ChevronLeft";
import ChevronRightIcon from "@material-ui/icons/ChevronRight";
import Divider from "@material-ui/core/Divider";
import List from "@material-ui/core/List";
import ListItem from "@material-ui/core/ListItem";
import ListItemIcon from "@material-ui/core/ListItemIcon";
import HomeIcon from "@material-ui/icons/Home";
import AlarmOnIcon from '@material-ui/icons/AlarmOn';
import ListItemText from "@material-ui/core/ListItemText";
import InboxIcon from "@material-ui/icons/MoveToInbox";
import {useTranslation} from 'react-i18next';
import Home from "../page/home.jsx";
import {drawerWidth} from "../js/constant"
import {AppContext} from "../App"
import SettingsIcon from '@material-ui/icons/Settings';
import {useHistory} from "react-router-dom";


const useStyles = makeStyles((theme) => ({
    root: {
        display: 'flex',
    },
    grow: {
        flexGrow: 1,
    },

    drawer: {
        width: drawerWidth,
        flexShrink: 0,
    },
    drawerPaper: {
        width: drawerWidth,
    },
    drawerHeader: {
        display: 'flex',
        alignItems: 'center',
        padding: theme.spacing(0, 1),
        // necessary for content to be below app bar
        ...theme.mixins.toolbar,
        justifyContent: 'flex-end',
    },

}));

export default function SideBar(props) {

    let history = useHistory();
    const classes = useStyles();

    const {drawerOpen,setDrawerOpen, newThingsOpen,setNewThingsOpen} = useContext(AppContext)

    console.log("layout sidebar", drawerOpen)
    const {t, i18n} = useTranslation();
    const theme = useTheme();

    function handleClick(url) {
        setDrawerOpen(false)
        history.push(url)

    }

    return (
        <>
            <Drawer
                className={classes.drawer}
                width="240"
                variant="persistent"
                anchor="left"
                open={drawerOpen}
                classes={{
                    paper: classes.drawerPaper,
                }}
            >
                <div className={classes.drawerHeader}>
                    <IconButton onClick={() => setDrawerOpen(false)}>
                        {theme.direction === 'ltr' ? <ChevronLeftIcon/> : <ChevronRightIcon/>}
                    </IconButton>
                </div>
                <Divider/>

                <List>

                    <ListItem button key={"home"} onClick={() => handleClick("/home")}>
                        <ListItemIcon>
                            <HomeIcon/>
                        </ListItemIcon>
                        <ListItemText primary={t('Home')}/>
                    </ListItem>

                    <ListItem button key={"rules"} onClick={() => handleClick("/rules")}>
                        <ListItemIcon>
                            <AlarmOnIcon/>
                        </ListItemIcon>
                        <ListItemText primary={t('Rules')}/>
                    </ListItem>

                    <ListItem button key={"settings"} onClick={() => handleClick("/settings")}>
                        <ListItemIcon>
                            <SettingsIcon/>
                        </ListItemIcon>
                        <ListItemText primary={t('Settings')}/>
                    </ListItem>


                    {['Inbox', 'Starred', 'Send email', 'Drafts'].map((text, index) => (
                        <ListItem button key={text}>
                            <ListItemIcon>{index % 2 === 0 ? <InboxIcon/> : <MailIcon/>}</ListItemIcon>
                            <ListItemText primary={text}/>
                        </ListItem>
                    ))}
                </List>
                <Divider/>
                <List>
                    {['All mail', 'Trash', 'Spam'].map((text, index) => (
                        <ListItem button key={text}>
                            <ListItemIcon>{index % 2 === 0 ? <InboxIcon/> : <MailIcon/>}</ListItemIcon>
                            <ListItemText primary={text}/>
                        </ListItem>
                    ))}
                </List>
            </Drawer>

        </>
    );
}
