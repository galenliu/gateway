import React, {useState} from "react";
import "../css/menu.css"
import {Trans} from 'react-i18next';
import {NavLink} from 'react-router-dom'


interface IStataController {
    menuBtnShow: boolean,
    menuListShow: boolean,
    menuScrimShow: boolean,
}

const Menu: React.FC = props => {


    function initState(): IStataController {
        return {
            menuBtnShow: true,
            menuListShow: false,
            menuScrimShow: false,
        }
    }

    const [state, setState] = useState<IStataController>(initState());


    function hidden() {
        setState({
            menuBtnShow: state.menuBtnShow,
            menuListShow: false,
            menuScrimShow: false,
        })

    }

    function toggleMenu() {
        setState({
            menuBtnShow: state.menuBtnShow,
            menuListShow: !state.menuListShow,
            menuScrimShow: !state.menuScrimShow,
        })
    }

    const MenuHidden = () => {
        setState({
            menuBtnShow: false,
            menuListShow: state.menuListShow,
            menuScrimShow: state.menuScrimShow,
        })
    }

    return (
        <React.Fragment>
            <div id="menu-scrim" onClick={hidden} className={state.menuScrimShow ? "" : "hidden"}/>
            <button onClick={toggleMenu} id="menu-button"
                    className={state.menuBtnShow ? "icon-button" : "icon-button hidden"}/>
            <nav id="main-menu" className={state.menuScrimShow ? "" : "hidden"}>
                <ul>
                    <li>
                        <NavLink onClick={() => {
                            hidden()
                        }} id="things-menu-item" to="/things" activeClassName="selected">
                            <Trans>Things</Trans></NavLink>
                    </li>
                    <li>
                        <NavLink onClick={() => {
                            hidden()
                        }} id="rules-menu-item" to="/rules" activeClassName="selected"><Trans>Rules</Trans></NavLink>
                    </li>
                    <li>
                        <NavLink onClick={() => {
                            hidden()
                        }} id="settings-menu-item" to="/settings"
                                 activeClassName="selected"><Trans>Settings</Trans></NavLink>
                    </li>
                </ul>
            </nav>
        </React.Fragment>
    )
};


export default Menu;
