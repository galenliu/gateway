import React, {useCallback, useContext} from "react";
import Light from "./capability/light";
import Grid from "@material-ui/core/Grid";
import Card from "@material-ui/core/Card";
import {ThingContext} from "../views/home"


export default function Thing(props) {

    const {handleSendMessage} = useContext(ThingContext)

    const setPropertyValue = useCallback((property) => {
            console.log("set property", property)
            const data = {}
            console.log(data)
            let message = {}
            message.messageType = "setProperty"
            message.id = props.thing.id
            message.data = {}
            message.data[property.name] = property.value
            console.log("set property message:",message)
            handleSendMessage(message)
        }, []
    )

    const propertyValueChanged = (thingId, value) => {

    }


    function RenderThing() {
        const data = props.thing
        console.log("thing render data:",data)
        let thing
        if (data.hasOwnProperty("selectedCapability")) {
            let selectedCapability = data["selectedCapability"]
            if (selectedCapability === "Light") {
                thing = <Light key={data.id} thing={data}/>
            }
        }
        return thing
    }

    return <><ThingContext.Provider value={{setPropertyValue}}>
        <Grid item xs={6} sm={3} id="thing-card">
            <Card elevation={5} >
                {RenderThing()}
            </Card>
        </Grid>
    </ThingContext.Provider>
    </>

}
