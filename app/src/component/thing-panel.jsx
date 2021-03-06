import React from 'react';
import {makeStyles, withStyles} from '@material-ui/core/styles';
import Button from '@material-ui/core/Button';
import Dialog from '@material-ui/core/Dialog';
import MuiDialogTitle from '@material-ui/core/DialogTitle';
import MuiDialogContent from '@material-ui/core/DialogContent';
import MuiDialogActions from '@material-ui/core/DialogActions';
import IconButton from '@material-ui/core/IconButton';
import CloseIcon from '@material-ui/icons/Close';
import Typography from '@material-ui/core/Typography';
import {useTranslation} from "react-i18next";
import ListItem from '@material-ui/core/ListItem';
import List from "@material-ui/core/List";
import {ListItemSecondaryAction, ListSubheader, Switch, useTheme} from "@material-ui/core";
import ListItemText from '@material-ui/core/ListItemText';
import ListItemIcon from '@material-ui/core/ListItemIcon';

import ThingIcon from "./thing-icon";
import Grid from "@material-ui/core/Grid";

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
        display:"flex",
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
    },
}));

const DialogTitle = withStyles(styles)((props) => {
    const {children, classes, onClose, ...other} = props;
    return (
        <MuiDialogTitle disableTypography className={classes.root} {...other}>
            <Typography variant="h6">{children}</Typography>
            {onClose ? (
                <IconButton aria-label="close" className={classes.closeButton} onClick={onClose}>
                    <CloseIcon/>
                </IconButton>
            ) : null}
        </MuiDialogTitle>
    );
});

const DialogContent = withStyles((theme) => ({
    root: {
        padding: theme.spacing(2),
    },
}))(MuiDialogContent);

const DialogActions = withStyles((theme) => ({
    root: {
        margin: 0,
        padding: theme.spacing(1),
    },
}))(MuiDialogActions);


export function ThingPanel(props) {

    const {t, i18n} = useTranslation();
    const classes = useStyles();
    return (
        <div>

            <Dialog onClose={() => props.show(false)} aria-labelledby="customized-dialog-title" open={props.open}>
                <DialogTitle id="customized-dialog-title" onClose={() => props.show(false)}>
                    {props.title}
                </DialogTitle>
                <Grid className={classes.content}>
                    <DetailsPanel  {...props}/>
                </Grid>

                <DialogActions>
                    <Button autoFocus onClick={() => props.show(false)} color="primary">
                        {t("Close")}
                    </Button>
                </DialogActions>
            </Dialog>
        </div>
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
                <ListItem>
                    <ListItemIcon>
                        <ThingIcon color={"red"} type={"Light"} {...props}/>
                    </ListItemIcon>
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

