import React from "react";

export const Actions = {
    addThing: "addonThing",
    initial: "initial",
    propertyStatus: "propertyStatus"
}

export default function ThingsReducer(context, action) {
    switch (action.type) {
        case Actions.propertyStatus:
            return context
        case Actions.addThing:
            console.log("Actions.addThing:", context)
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

        case Actions.initial:
            console.log("Actions.initial:", action.initialState)
            return action.initialState
        default:
            return context
    }
}



