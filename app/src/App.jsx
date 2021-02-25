import './App.css';
import React, {useState} from "react";
import SideBar from "./component/sideBar";
import "./i18n"
import {HashRouter as Router, Route, Switch} from 'react-router-dom'
import Home from "./page/home";

export const AppContext = React.createContext()


function App() {

    const [open, setOpen] = useState(false)

    function handleDrawerClose(){
        setOpen(false)
    }
    function handleDrawerOpen(){
        setOpen(true)
    }

    return (
        <AppContext.Provider value={{drawerOpenState: open,drawerClose: handleDrawerClose}}>
            <Router>
                <Switch>
                    <Route exact path="/things">
                        <SideBar/>
                        <Home/>
                    </Route>
                    <Route path="/">
                        <SideBar/>
                        <Home/>
                    </Route>
                </Switch>
            </Router>
        </AppContext.Provider>
    );
}

export default App;
