import React from "react";
import {Trans} from "react-i18next";
import "../css/settings.css";

interface adapterInfo {
    name: string,
    version: string,
    description: string,
    homepage_url: string,
    author: string,
    license: string,
    enabled: boolean,
    isUpdate: boolean,

}

const Adapter: React.FC<adapterInfo> = props => {

    return <li className="addon-item">
        <div className="addon-settings-header adapter">
            <span className="addon-settings-name">{props.name}</span>
            <span className="addon-settings-version">{props.version}</span>
            <span className="addon-settings-description">{props.description}</span>
            <span className="addon-settings-author"><Trans>author</Trans>
                <a href={props.homepage_url} target="_blank" rel="noopener">{props.author}</a>
            </span>
            <span className="addon-settings-license"
                  data-license-href={props.license}><Trans>license</Trans>
           </span>
        </div>
        <div className="addon-settings-controls">
            <button className="text-button addon-settings-config"><Trans>configure</Trans>
            </button>
            {props.isUpdate && <button
                className="text-button addon-settings-update"><Trans>update</Trans>
            </button>}
            <span className="addon-settings-spacer"></span>
            <button className="text-button addon-settings-remove"
            ><Trans>remove</Trans>
            </button>
            <button
                className={props.enabled ? "text-button addon-settings-disable" : "text-button addon-settings-enable"}>
                {props.enabled ? <Trans>disable</Trans> : <Trans>enable</Trans>}
            </button>
        </div>
    </li>
}


export default Adapter;

