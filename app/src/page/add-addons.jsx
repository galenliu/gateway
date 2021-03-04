import React, {useEffect, useState} from "react";
import Dialog from "@material-ui/core/Dialog";
import AppBar from "@material-ui/core/AppBar";
import Toolbar from "@material-ui/core/Toolbar";
import Typography from "@material-ui/core/Typography";
import IconButton from "@material-ui/core/IconButton";
import CloseIcon from "@material-ui/icons/Close";
import Grid from "@material-ui/core/Grid";
import {useTranslation} from "react-i18next";
import Slide from '@material-ui/core/Slide';
import API from "../js/api";
import StoreIcon from '@material-ui/icons/Store';
import Card from "@material-ui/core/Card";
import {Button, CircularProgress} from "@material-ui/core";
import Divider from "@material-ui/core/Divider";
import {Addon, useStyles} from "./addons";


const Transition = React.forwardRef(function Transition(props, ref) {
    return <Slide direction="up" ref={ref} {...props} />;
});




