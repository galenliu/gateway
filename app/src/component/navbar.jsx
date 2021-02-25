import React from 'react';
import {makeStyles, useTheme} from '@material-ui/core/styles';
import {useTranslation} from 'react-i18next';
import clsx from "clsx";
import CssBaseline from "@material-ui/core/CssBaseline";
import SideBar from "./sideBar";
import top from "./topBar";
import {drawerWidth} from "../js/constant"


const useStyles = makeStyles((theme) => ({
    root: {
        display: 'flex',
    },

    drawerHeader: {
        display: 'flex',
        alignItems: 'center',
        padding: theme.spacing(0, 1),
        // necessary for content to be below app bar
        ...theme.mixins.toolbar,
        justifyContent: 'flex-end',
    },

    title: {
        display: 'none',
        [theme.breakpoints.up('sm')]: {
            display: 'block',
        },
    },
    content: {
        flexGrow: 1,
        padding: theme.spacing(3),
        transition: theme.transitions.create('margin', {
            easing: theme.transitions.easing.sharp,
            duration: theme.transitions.duration.leavingScreen,
        }),
        marginLeft: -drawerWidth,
    },
    contentShift: {
        transition: theme.transitions.create('margin', {
            easing: theme.transitions.easing.easeOut,
            duration: theme.transitions.duration.enteringScreen,
        }),
        marginLeft: 0,
    },

    inputRoot: {
        color: 'inherit',
    },
    inputInput: {
        padding: theme.spacing(1, 1, 1, 0),
        // vertical padding + font size from searchIcon
        paddingLeft: `calc(1em + ${theme.spacing(4)}px)`,
        transition: theme.transitions.create('width'),
        width: '100%',
        [theme.breakpoints.up('md')]: {
            width: '20ch',
        },
    },

}));

export default function Navbar(props) {
    const classes = useStyles();
    const {t, i18n} = useTranslation();


    //侧边栏打开、关闭状态
    const [open, setOpen] = React.useState(false);

    const handleDrawerOpen = () => {
        setOpen(true);
    };
    const handleDrawerClose = () => {
        setOpen(false);
    };

    const handleAddNewThings = () => {
        setOpen(false);
    };

    return (
        <div className={classes.root}>
            <CssBaseline/>
            <Topbar handleOpen={handleDrawerOpen} open={open} addNewThings={handleAddNewThings}/>
            <SideBar handleClose={handleDrawerClose} open={open}/>


            <main
                className={clsx(classes.content, {
                    [classes.contentShift]: open,
                })}
            >
                <div className={classes.drawerHeader}/>
            </main>


        </div>
    );
}
