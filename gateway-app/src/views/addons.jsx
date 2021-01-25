import React, {useEffect, useState} from "react";
import "../css/settings.css"
import Adapter from "../components/adapter";
import {Link} from "react-router-dom";
import Addon from "../components/addon";
import API from "../api"


// Adapter like this
//------------
// author: "Tim Hellhake"
// description: "Runs for the specified number of seconds"
// enabled: true
// exec: "{nodeLoader} {path}"
// homepage_url: "https://github.com/tim-hellhake/timer-adapter#readme"
// id: "timer-adapter"
// name: "Timer"
// primary_type: "adapter"
// version: "1.6.1"
//----------------


// Addon like this
//------------
// "id": "yo-notifier",
// "name": "Yo",
// "description": "Send a Yo to your phone (gateway >= 0.9.0 only)",
// "author": "Tim Hellhake",
// "homepage_url": "https://github.com/tim-hellhake/yo-notifier#readme",
// "license_url": "https://raw.githubusercontent.com/tim-hellhake/yo-notifier/master/LICENSE",
// "version": "0.1.4",
//"url": "https://github.com/tim-hellhake/yo-notifier/releases/download/0.1.4/yo-notifier-0.1.4.tgz",
// "checksum": "83fb9753956b6e786137d294d801b4589971a53b5c6f04618bd3855ddeea5862",
// "primary_type": "notifier"
//------------


const Addons = () => {

    const [installedAddons, setInstalledAddons] = useState(new Map())
    const [availableAddons, setAvailableAddons] = useState(new Map())


    //获取到的gateway已安装的addons
    const [currentAddonsList, setCurrentAddonList] = useState()

    function fetchInstalledAddonsList() {
        console.log("begin fetch the installed addon list")
        API.getInstalledAddons().then(body => {
            if (!body) {
                return;
            }
            let installedList = new Map()
            for (const s of body) {
                try {
                    installedList.set(s.id, s)
                } catch (err) {
                    console.error(`Failed to parse add-on settings: ${err}`);
                }
            }
            setInstalledAddons(installedList)
        }).catch((err) => {
            console.log("fetch the installed addon list err:", err)
        })
    }

    function fetchAvailableAddonList() {
        API.getAddonsInfo().then((data) => {
            if (!data || !data.urls || !data.architecture || !data.version) {
                console.log(data.urls)
                return;
            }

            const params = new URLSearchParams();
            params.set('arch', data.architecture);
            params.set('version', data.version);
            params.set('node', data.node_version);
            params.set('python', data.python_version);
            const promises = [];
            for (const url of data.urls) {
                console.log(url)
                promises.push(fetch(`${url}?${params.toString()}`, {
                    method: 'GET',
                    cache: 'reload',
                    headers: {
                        Accept: 'application/json',
                    },
                }));
            }
            console.log(promises, "promises.....")
            return Promise.all(promises)
        }).then(
            (responses) => {
                const promises = []
                for (const resp of responses) {
                    promises.push(resp.json())
                }
                return Promise.all(promises)
            }).then((bodies) => {
                let availableList = availableAddons
                for (const body of bodies) {
                    for (const addon of body) {
                        const entry = {
                            id: addon.id,
                            name: addon.name,
                            description: addon.description,
                            author: addon.author,
                            homepage_url: addon.homepage,
                            license_url: addon.license_url,
                            version: addon.version,
                            url: addon.url,
                            checksum: addon.checksum,
                            primary_type: addon.primary_type,
                            installed: installedAddons.has(addon.id)
                        }
                        if (availableList.has(addon.id) &&
                            this.compareSemver(availableList.get(addon.id).version,
                                entry.version) >= 0) {
                            continue;
                        }
                        availableList.set(addon.id, entry)

                    }
                }
                setAvailableAddons(availableList);
                console.log("availableList:", availableList)
            }
        ).catch((err) => {
            console.error(err)
        })
    }

    function installAddon(id, url, checksum) {
        API.installAddon(id, url, checksum).catch((err) => {
            console.error("install addon fail err", err)
        })
    }

    function renderInstalledAddonsList(adapterList) {
        const list = []
        for (const [id, adapter] of adapterList) {
            const addon = <Adapter key={id}
                                   name={adapter.name}
                                   author={adapter.author}
                                   description={adapter.description}
                                   homepage_url={adapter.homepage_url}
                                   license={adapter.license_url}
                                   version={adapter.version}
                                   primary_type={adapter.primary_type}
                                   exec={adapter.exec}
                                   enabled={adapter.enabled}
                                   updata={false}
                                   schema={adapter.schema}
                                   isUpdate={adapter.isUpdate}
            />
            list.push(addon)
        }
        return list
    }

    function renderAvailableAddonsList(addonsList) {

        let addonList = []
        for (const [id, entry] of addonsList) {
            const addon =
                <Addon key={id}
                       id={id}
                       name={entry.name}
                       author={entry.author}
                       description={entry.description}
                       homepage_url={entry.homepage_url}
                       license_url={entry.license_url}
                       version={entry.version}
                       primary_type={entry.primary_type}
                       checksum={entry.checksum}
                       url={entry.url}
                       onAdd={installAddon}/>
            addonList.push(addon)
        }
        return addonList
    }


    useEffect(() => {
        fetchInstalledAddonsList()
        fetchAvailableAddonList()
    }, [])


    const [isDiscover, setDiscover] = useState(false);

    return (
        <>
            <section id="addon-settings" className="settings-section">

                <section id="addon-main-settings"
                         className={isDiscover ? "settings-subsection hidden" : "settings-subsection"}>
                    <button id="discover-addons-button" onClick={() => setDiscover(true)}
                            className="icon-button right-top-button"/>
                    <Link to="/settings">
                        <button id="back-button" className="icon-button"/>
                    </Link>
                    <ul id="installed-addons-list">
                        {
                            installedAddons && renderInstalledAddonsList(installedAddons)
                        }
                    </ul>
                </section>

                {/*<section id="addon-config-settings" className="hidden settings-subsection"/>*/}


                <section id="addon-discovery-settings"
                         className={isDiscover ? "settings-subsection" : "settings-subsection hidden"}>

                    <button onClick={() => setDiscover(false)} id="back-button" className="icon-button"/>

                    <div id="addon-discovery-container" className=''>
                        <input id="discovered-addons-search" className="hidden"/>
                        <ul id="discovered-addons-list">
                            {availableAddons && renderAvailableAddonsList(availableAddons)};
                        </ul>
                    </div>
                </section>

            </section>
        </>
    );
}


export default Addons;
