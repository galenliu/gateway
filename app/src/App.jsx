import './App.css';
import React, {useState} from "react";
import SideBar from "./component/sideBar";
import "./i18n"
import {HashRouter as Router, Route, Switch} from 'react-router-dom'
import Home from "./page/home";
import Settings from "./page/settings";

export const AppContext = React.createContext()


export default function App() {

    const [drawerOpen, setDrawerOpen] = useState(false)
    const [newThingsOpen, setNewThingsOpen] = useState(false)


    return (
        <AppContext.Provider value={{
            drawerOpen: drawerOpen,
            setDrawerOpen: setDrawerOpen,
            newThingsOpen: newThingsOpen,
            setNewThingsOpen: setNewThingsOpen,
        }}>
            <Router>
                <Switch>
                    <Route exact path="/things">
                        <SideBar/>
                        <Home/>
                    </Route>
                    <Route exact path="/settings">
                        <SideBar/>
                        <Settings/>
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

