import React from 'react';
import './App.css';
import Home from './views/home'
import Rules from './views/rules'
import Settings from './views/settings'
import Menu from "./components/menu";
import {HashRouter as Router, Route, Switch} from 'react-router-dom'
import Addons from './views/addons'
import ThingsAdd from "./views/thing-add";
import './i18n'


function App() {
    //const [state, dispatch] = useReducer(goto, {path: "/things"})

    return (
        <Router>
            <Switch>
                <Route exact path="/things">
                    <Menu/>
                    <Home/>
                </Route>
                <Route path="/things/add">
                    <ThingsAdd/>
                </Route>
                <Route path="/rules">
                    <Menu/>
                    <Rules/>
                </Route>
                <Route exact path="/settings">
                    <Menu/>
                    <Settings/>
                </Route>
                <Route path="/settings/addons">
                    <Addons/>
                </Route>
                <Route path="/">
                    <Menu/>
                    <Home/>
                </Route>
            </Switch>
        </Router>
    );
}


export default App;
