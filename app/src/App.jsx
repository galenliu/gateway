import './App.css';
import React, {useEffect, useState} from "react";
import SideBar from "./component/sideBar";
import "./i18n"
import {HashRouter as Router, Route, Switch} from 'react-router-dom'
import Things from "./views/Things";
import Settings from "./views/Settings";
import {createMuiTheme, ThemeProvider} from '@material-ui/core/styles';
import Core from "./core";


export const AppContext = React.createContext()
const theme = createMuiTheme({
    status: {},
});


export default function App() {

    useEffect(() => {
        Core.init()
    }, [])

    const [drawerOpen, setDrawerOpen] = useState(false)
    const [newThingsOpen, setNewThingsOpen] = useState(false)

    return (
        <ThemeProvider theme={theme}>
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
                            <Things/>
                        </Route>
                        <Route exact path="/settings">
                            <SideBar/>
                            <Settings/>
                        </Route>
                        <Route path="/">
                            <SideBar/>
                            <Things/>
                        </Route>

                    </Switch>
                </Router>
            </AppContext.Provider>
        </ThemeProvider>
    );
}

