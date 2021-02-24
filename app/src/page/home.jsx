import React, {useEffect, useReducer,useCallback} from "react";
import Thing from "../component/thing.jsx";
import Grid from "@material-ui/core/Grid";
import mock from "../mock/things.json"

import ThingsReducer,{Actions} from "../js/things-reducer";

export const ThingsContext = React.createContext()

function Home() {

    const [things, setThings] = useReducer(ThingsReducer, [])

    useEffect(() => {
        console.log("things",things)
        mock.forEach((thing) => {
            console.log(thing)
            setThings({messageType: Actions.addThing, thing: thing})
        })
    })

    const handleSendMessage = useCallback((data) =>
           console.log(JSON.stringify(data)),
        []);


    function RenderThingsView() {
        let list = []
        for (let t of things) {
            console.log("item:----",t)
            const thing =
                <Thing key={t.id} thing={t}/>
            list.push(thing)
        }
        return list
    }

    return (

            <ThingsContext.Provider value={{handleSendMessage}}>
                    <Grid container spacing={1}>
                        {RenderThingsView()}
                    </Grid>
            </ThingsContext.Provider>

    );
}

export default Home;
