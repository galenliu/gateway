import React from "react";

export const Actions = {
    addThing: "addonThing",
    propertyStatus: "propertyStatus"
}

export default function ThingsReducer(context, action) {
    console.log("action:", action)
    switch (action.messageType) {
        case Actions.propertyStatus:
            return context
        case Actions.addThing:
            console.log("thingsReducer addThing:", action)
            console.log("state", context)
            let newThing = action.thing
            for (let thing of context) {
                if (thing.id === newThing.id) {
                    return context
                }
            }
            let newThings = context
            newThings.push(newThing)
            console.log("newThings:", newThings)
            return [...newThings]
        default:
            return context
    }
}



