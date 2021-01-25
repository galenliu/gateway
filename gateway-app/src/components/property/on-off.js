import React, {useContext} from "react";
import Switch from '@material-ui/core/Switch';
import {ThingContext} from "../../views/home";

const OnOffProperty = (props) => {

    const {name, value} = props

    const {setPropertyValue} = useContext(ThingContext)

    const [state, setState] = React.useState(props.data);

    const handleChange = (event) => {
        setState({...state, value: !event.target.checked})
        setPropertyValue(state)
    };


    return (
        <>
            <div id={"property-item"} className={"property-on-off-switch"}>
                <Switch
                    checked={!state.value}
                    onChange={handleChange}
                    color="primary"
                    name="on-off-switch"
                    inputProps={{'aria-label': 'primary checkbox'}}
                />
            </div>
        </>
    )
}

export default OnOffProperty
