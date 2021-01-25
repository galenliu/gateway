import React, {useContext} from "react";
import Slider from '@material-ui/core/Slider';
import {ThingContext} from "../../views/home";


const BrightnessProperty = (props) => {

    const [state, setState] = React.useState(props.data)
    const {setPropertyValue} = useContext(ThingContext)

    const handleChange = (event, newValue) => {
        if (state.value !== newValue){
            setState({...state, value: newValue})
            setPropertyValue(state)
        }
    };

    return (
        <>
            <div id={"property-item"} className={"property-brightness"}>
                <Slider
                    value={state.value}
                    min={props.data.minimum}
                    step={1}
                    max={props.data.maximum}
                    onChange={handleChange}
                    valueLabelDisplay="auto"
                    aria-labelledby="non-linear-slider"
                />
            </div>
        </>
    )
}

export default BrightnessProperty
