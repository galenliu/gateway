import React, {useCallback, useContext, useEffect, useReducer, useState} from "react";
import Thing from "../component/thing.jsx";
import Grid from "@material-ui/core/Grid";
import API from "../js/api";
import NewThingsDialog from "./new-things";

import ThingsReducer from "../js/things-reducer";
import TopBar from "../component/topBar";
import {AppContext} from "../App";
import {useTranslation} from "react-i18next";

export const HomeContext = React.createContext()

function Home() {

    const [things, thingsDispatch] = useReducer(ThingsReducer, [])
    const {t, i18n} = useTranslation();

    useEffect(() => {
        API.getThings().then((data) => {
            console.log("data", data)
        }).catch((e) => {
            console.log("getThings err:", e)
        })
        console.log("things", things)

    }, [])

    const handleSendMessage = useCallback((data) =>
            console.log(JSON.stringify(data)),
        []);


    function RenderThingsView() {
        let list = []
        for (let t of things) {
            console.log("item:----", t)
            const thing =
                <Thing key={t.id} thing={t}/>
            list.push(thing)
        }
        return list
    }


    return (
        <>
            <TopBar add={true} title={t("Home")}/>
            <Grid container spacing={1}>
                {RenderThingsView()}
            </Grid>
            <NewThingsDialog/>  </>

    );
}

export default Home;
