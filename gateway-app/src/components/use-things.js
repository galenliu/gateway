import React from "react";


export default function thingsReducer(context, action) {
    console.log("action:", action)
    switch (action.messageType) {
        case "propertyStatus":

            return context
        case "addThing":
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



