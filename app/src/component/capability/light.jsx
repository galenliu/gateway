import React, {useEffect} from "react";
import {SvgIcon} from "@material-ui/core";
import {yellow} from "@material-ui/core/colors";

const Light = (props) => {
    function LightOnIcon() {
        return (
            <SvgIcon style={{fontSize: 40}}>
                <path fill="currentColor"
                      d="M12,2A7,7 0 0,0 5,9C5,11.38 6.19,13.47 8,14.74V17A1,1 0 0,0 9,18H15A1,1 0 0,0 16,17V14.74C17.81,13.47 19,11.38 19,9A7,7 0 0,0 12,2M9,21A1,1 0 0,0 10,22H14A1,1 0 0,0 15,21V20H9V21Z"/>
            </SvgIcon>
        );
    }

    function LightOffIcon() {
        return (
            <SvgIcon style={{fontSize: 40, color: yellow[900]}}>
                <path fill="currentColor"
                      d="M12,2A7,7 0 0,0 5,9C5,11.38 6.19,13.47 8,14.74V17A1,1 0 0,0 9,18H15A1,1 0 0,0 16,17V14.74C17.81,13.47 19,11.38 19,9A7,7 0 0,0 12,2M9,21A1,1 0 0,0 10,22H14A1,1 0 0,0 15,21V20H9V21Z"/>
            </SvgIcon>
        );
    }

    useEffect(() => {
        update(props.thing)
        props.updataState(<h1>打开</h1>)
        props.updataIcon(LightOnIcon())
    }, [])

    function update(thing) {
        props.updataState(<h1>打开</h1>)
        if (thing) {
            if (thing.value) {
                if (thing.value == true) {
                    props.updataIcon(LightOnIcon())
                }
                if (thing.value == false) {
                    props.updataIcon(LightOffIcon())
                }
            }
        }
    }

    return (
        <>
            <h1> 亮度
            </h1>
            <switch></switch>
        </>
    )

    // renderProperty() {
    //     const list = []
    //     for (const prop in props.thing.properties) {
    //         if (props.thing.properties.hasOwnProperty(prop)) {
    //             const p = props.thing.properties[prop]
    //             if (p.hasOwnProperty("@type")) {
    //                 let type = p["@type"]
    //                 if (type === "OnOffProperty") {
    //                     const propComponent =
    //                         <OnOffProperty data={p} key={p.name} >
    //                         </OnOffProperty>
    //
    //                     list.push(propComponent)
    //                 }
    //                 if (type === "BrightnessProperty") {
    //                     const propComponent =
    //                         <BrightnessProperty data={p} key={p.name} >
    //                         </BrightnessProperty>
    //
    //                     list.push(propComponent)
    //                 }
    //             }
    //         }
    //     }
    //     return list
    // },

}

export default Light;



