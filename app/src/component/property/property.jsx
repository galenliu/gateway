import {useContext} from "react";
import {ThingsContext} from "../../views/Things.jsx";



export default function property(Component, props) {

    return <Component props={props}/>

}
