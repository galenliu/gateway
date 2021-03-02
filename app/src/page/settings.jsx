import React, {useContext} from "react";
import {makeStyles} from "@material-ui/core/styles";
import ListItemIcon from '@material-ui/core/ListItemIcon';
import ListItemText from '@material-ui/core/ListItemText';
import ListItem from '@material-ui/core/ListItem';
import List from "@material-ui/core/List";
import {useTranslation} from "react-i18next";
import DomainIcon from '@material-ui/icons/Domain';
import TopBar from "../component/topBar";
import {CssBaseline} from "@material-ui/core";
import Typography from "@material-ui/core/Typography";
import clsx from "clsx";
import {AppContext} from "../App";


const useStyles = makeStyles((theme) => ({
    root: {
        display: "flex",

    },
    content: {
        flexGrow: 1,
        padding: theme.spacing(3),
        transition: theme.transitions.create('margin', {
            easing: theme.transitions.easing.sharp,
            duration: theme.transitions.duration.leavingScreen,
        }),
        marginLeft: -240,
    },
    contentShift: {
        transition: theme.transitions.create('margin', {
            easing: theme.transitions.easing.easeOut,
            duration: theme.transitions.duration.enteringScreen,
        }),
        marginLeft: 0,
    },

    LItem: {},

}))


export default function Settings(props) {

    const classes = useStyles();
    const {t, i18n} = useTranslation();
    const {drawerOpen, setNewThingsClose, setNewThingsOpen, newThingsOpen} = useContext(AppContext)


    return (
        <>
            <CssBaseline/>
            <TopBar/>
            <div>
                <List component="nav" aria-label="main mailbox folders">
                    <ListItem>
                        <ListItemIcon>
                            <DomainIcon/>
                        </ListItemIcon>
                        <ListItemText primary="Inbox"/>
                    </ListItem>
                </List>
            </div>
            <main
                className={clsx(classes.content, {
                    [classes.contentShift]: !drawerOpen,
                })}
            >
                <div className={classes.drawerHeader}/>
                <Typography paragraph>
                    Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt
                    ut labore et dolore magna aliqua. Rhoncus dolor purus non enim praesent elementum
                    facilisis leo vel. Risus at ultrices mi tempus imperdiet. Semper risus in hendrerit
                    gravida rutrum quisque non tellus. Convallis convallis tellus id interdum velit laoreet id
                    donec ultrices. Odio morbi quis commodo odio aenean sed adipiscing. Amet nisl suscipit
                    adipiscing bibendum est ultricies integer quis. Cursus euismod quis viverra nibh cras.
                    Metus vulputate eu scelerisque felis imperdiet proin fermentum leo. Mauris commodo quis
                    imperdiet massa tincidunt. Cras tincidunt lobortis feugiat vivamus at augue. At augue eget
                    arcu dictum varius duis at consectetur lorem. Velit sed ullamcorper morbi tincidunt. Lorem
                    donec massa sapien faucibus et molestie ac.
                </Typography>
                <Typography paragraph>
                    Consequat mauris nunc congue nisi vitae suscipit. Fringilla est ullamcorper eget nulla
                    facilisi etiam dignissim diam. Pulvinar elementum integer enim neque volutpat ac
                    tincidunt. Ornare suspendisse sed nisi lacus sed viverra tellus. Purus sit amet volutpat
                    consequat mauris. Elementum eu facilisis sed odio morbi. Euismod lacinia at quis risus sed
                    vulputate odio. Morbi tincidunt ornare massa eget egestas purus viverra accumsan in. In
                    hendrerit gravida rutrum quisque non tellus orci ac. Pellentesque nec nam aliquam sem et
                    tortor. Habitant morbi tristique senectus et. Adipiscing elit duis tristique sollicitudin
                    nibh sit. Ornare aenean euismod elementum nisi quis eleifend. Commodo viverra maecenas
                    accumsan lacus vel facilisis. Nulla posuere sollicitudin aliquam ultrices sagittis orci a.
                </Typography>
            </main>
        </>
    )
}