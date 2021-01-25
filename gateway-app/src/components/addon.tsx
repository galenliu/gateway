import React from "react";
import {Trans} from 'react-i18next';
import "../css/settings.css";

interface addonProps {
    id: string,
    name: string,
    version: string,
    description: string,
    author: string,
    license_url: string,
    homepage_url: string,
    primary_type: string,
    checksum: string,
    url: string,
    onAdd: (addonId: string, addonUrl: string, addonChecksum: string) => {}
}


const Addon: React.FC<addonProps> = props => {

    // @ts-ignore
    return (
        <li className="discovered-addon-item">
            <div className="addon-settings-header adapter">
                <span className="addon-settings-name">{props.name}</span>
                <span className="addon-settings-version">{props.version}</span>
                <span className="addon-settings-description">{props.description}</span>
                <span className="addon-settings-author">by
                    <a href={props.homepage_url} target="_blank" rel="noopener">{props.author}</a>
                </span>
                <span id="addon-license-calendar" className="addon-settings-license"
                      data-license-href={props.license_url}>(license)
                </span>
            </div>
            <div className="addon-settings-controls">
                <button onClick={() => {
                    props.onAdd(props.id, props.url, props.checksum)
                }} id="addon-install-calendar" className="text-button addon-discovery-settings-add"
                        data-l10n-id="addon-discovery-add"><Trans>add</Trans>
                </button>
            </div>
        </li>
    )
}


export default Addon;

