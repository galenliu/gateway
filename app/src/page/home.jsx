import React, {useCallback, useContext, useEffect, useReducer, useState} from "react";
import Thing from "../component/thing.jsx";
import Grid from "@material-ui/core/Grid";
import mock from "../mock/things.json"
import API from "../js/api";
import NewThingsDialog from "./newThings";

import ThingsReducer, {Actions} from "../js/things-reducer";
import TopBar from "../component/topBar";

export const HomeContext = React.createContext()

function Home() {

    const [things, thingsDispatch] = useReducer(ThingsReducer, [])
    const [newThingsAdd, setNewThingsAdd] = useState(false)


    function handleNewThingsAddClose() {
        setNewThingsAdd(false)
    }
    function handleNewThingsAddOpen() {
        setNewThingsAdd(true)
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

        <HomeContext.Provider value={{open: newThingsAdd,close: handleNewThingsAddClose,handleSendMessage:handleSendMessage}}>
            <>
                <TopBar/>
                <Grid container spacing={1}>
                    {RenderThingsView()}
                </Grid>
            </>
            <NewThingsDialog />
        </HomeContext.Provider>
    );
}

export default Home;
