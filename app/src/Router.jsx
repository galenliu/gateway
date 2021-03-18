import React, {useState} from "react";
import SideBar from "./component/sideBar";
import "./i18n"
import {HashRouter as AppRouter, Route, Switch} from 'react-router-dom'
import Things from "./views/Things";
import Settings from "./views/Settings";
import {createMuiTheme, ThemeProvider} from '@material-ui/core/styles';
import App from "./App";


export const AppContext = React.createContext()
const theme = createMuiTheme({
    status: {},
});

App.init()

export default function Router() {
    const [drawerOpen, setDrawerOpen] = useState(false)
    const [newThingsOpen, setNewThingsOpen] = useState(false)

    return (
        <ThemeProvider theme={theme}>
            <AppContext.Provider value={{
                drawerOpen: drawerOpen,
                setDrawerOpen: setDrawerOpen,
                newThingsOpen: newThingsOpen,
                setNewThingsOpen: setNewThingsOpen,
                ts: App.gatewayModel.things,
            }}>
                <AppRouter>
                    <Switch>
                        <Route exact path="/things">
                            <SideBar/>
                            <Things things={App.gatewayModel.things}/>
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
                </AppRouter>
            </AppContext.Provider>
        </ThemeProvider>
    );
}

