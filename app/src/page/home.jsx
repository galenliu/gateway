import React, {useCallback, useEffect, useReducer, useState} from "react";
import Thing from "../component/thing.jsx";
import Grid from "@material-ui/core/Grid";
import mock from "../mock/things.json"
import API from "../js/api";
import NewThingsDialog from "./addThings";

import ThingsReducer, {Actions} from "../js/things-reducer";
import TopBar from "../component/topBar";

export const HomeContext = React.createContext()

function Home() {

    const [things, thingsDispatch] = useReducer(ThingsReducer, [])
    const [newThingsShow, setNewThingsShow] = useState(false)


    function handleNewThingsAddClose() {
        setNewThingsShow(false)
    }

    function handleNewThingsAddOpen() {
        setNewThingsShow(true)
    }

    useEffect(() => {
        API.getThings().then((data) => {
            console.log("data", data)
        }).catch((e) => {
            console.log(e)
            return
        })
        console.log("things", things)
        mock.forEach((thing) => {
            console.log(thing)
            thingsDispatch({messaegType: Actions.addThing, thing: thing})
        })
    })

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

        <HomeContext.Provider value={{
            open: newThingsShow,
            setNewThingsClose: handleNewThingsAddClose,
            setNewThingsOpen: handleNewThingsAddOpen
        }}>
            <>
                <TopBar/>
                <Grid container spacing={1}>
                    {RenderThingsView()}
                </Grid>
            </>
            <NewThingsDialog/>
        </HomeContext.Provider>
    );
}

export default Home;
