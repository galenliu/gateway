import './App.css';
import React, {useState} from "react";
import SideBar from "./component/sideBar";
import "./i18n"
import {HashRouter as Router, Route, Switch} from 'react-router-dom'
import Home from "./page/home";

export const AppContext = React.createContext()


function App() {

    const [drawerOpen, setDrawerOpen] = useState(false)

    function handleDrawerClose() {
        setDrawerOpen(false)
    }

    function handleDrawerOpen() {
        setDrawerOpen(true)
    }

    return (
        <AppContext.Provider value={{open: drawerOpen, setOpen: handleDrawerOpen, setClose: handleDrawerClose}}>
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
