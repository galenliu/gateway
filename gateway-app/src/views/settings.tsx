import React from 'react'
import {Trans, useTranslation} from "react-i18next";
import {Link} from "react-router-dom";
import "../css/settings.css";


const Settings: React.FC = () => {

    return (
        <>
            <nav id="settings-menu" className="settings-section">

                <ul>
                    <li className="settings-item"><Link id="domain-settings-link"
                                                        to="/settings/domain"><Trans>Domain</Trans></Link></li>
                    <li className="settings-item"><Link id="network-settings-link"
                                                        to="/settings/network"><Trans>Network</Trans></Link></li>
                    <li className="settings-item"><Link id="user-settings-link"
                                                        to="/settings/users"><Trans>Users</Trans></Link>
                    </li>
                    <li className="settings-item"><Link id="addon-settings-link"
                                                        to="/settings/addons"><Trans>Addons</Trans></Link></li>

                    <li className="settings-item"><Link id="localization-settings-link"
                                                        to="/settings/localization"><Trans>Localization</Trans></Link>
                    </li>
                    <li className="settings-item"><Link id="update-settings-link"
                                                        to="/settings/updates"><Trans>Updates</Trans></Link></li>
                    <li className="settings-item"><Link id="authorization-settings-link"
                                                        to="/settings/authorizations"><Trans>Authorizations</Trans></Link>
                    </li>
                    <li className="settings-item"><Link id="experiment-settings-link"
                                                        to="/settings/experiments"><Trans>Experiments</Trans></Link>
                    </li>
                    <li className="settings-item"><Link id="developer-settings-link"
                                                        to="/settings/developer"><Trans>Developer</Trans></Link></li>
                </ul>
            </nav>
        </>
    );
}

export default Settings
