import {useContext} from "react";
import {ThingContext} from "../../views/home";


const {setPropertyValue} = useContext(ThingContext)

export default function property(Component, props) {

    return <Component props={props}/>

}
