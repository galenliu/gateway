import React, {useEffect, useState} from "react";
import Dialog from "@material-ui/core/Dialog";
import AppBar from "@material-ui/core/AppBar";
import Toolbar from "@material-ui/core/Toolbar";
import Typography from "@material-ui/core/Typography";
import IconButton from "@material-ui/core/IconButton";
import CloseIcon from "@material-ui/icons/Close";
import Grid from "@material-ui/core/Grid";
import {makeStyles} from "@material-ui/core/styles";
import {useTranslation} from "react-i18next";
import Slide from '@material-ui/core/Slide';
import API from "../js/api";
import ExtensionIcon from "@material-ui/icons/Extension";
import AddCircleIcon from '@material-ui/icons/AddCircle';
import Card from "@material-ui/core/Card";
import Icons from "../component/icons";
import {Button, CircularProgress, FormHelperText, Link} from "@material-ui/core";
import Divider from "@material-ui/core/Divider";
import StoreIcon from "@material-ui/icons/Store";
import {versionStringCompare} from "../js/util";
import AddIcon from "@material-ui/icons/Add";


export const useStyles = makeStyles((theme) => ({
    root: {},
    addonCard: {
        padding: 3,
        display: "flex",
        maxWidth: 400,
        minWidth: 240,
        margin: 5,
        justifyContent: "space-between",
        alignItems: "center",
    },
    cardContent: {
        display: "flex",
        minWidth: 200,
        maxWidth: 200,
        flexDirection: "column",
        justifyContent: "flex-start",
        alignItems: "flex-start",
    },
    sideContent: {
        display: "flex",
        flexDirection: "column",
        direction: "column",
        minWidth: 80,
        justifyContent: "flex-start",
        alignItems: "stretch",
    },
    content: {

        justifyContent: "flex-start", alignItems: "center", direction: "column",

    },
    drawerHeader: {
        display: 'flex',
        alignItems: 'center',
        padding: theme.spacing(0, 1),
        // necessary for content to be below app bar
        ...theme.mixins.toolbar,
        justifyContent: 'flex-end',
    },
    appBar: {},
    title: {
        marginLeft: theme.spacing(2),
        flex: 1,
    },
}));
const Transition = React.forwardRef(function Transition(props, ref) {
    return <Slide direction="up" ref={ref} {...props} />;
});

const states = {
    Loading: "loading",
    Empty: "empty",
    Completed: "completed",
}

export default function AddonsDialog(props) {

    const classes = useStyles();
    const {t, i18n} = useTranslation();
    const [fetchAddonsShow, setFetchAddonsShow] = useState(false)
    const [installedAddons, setInstalledAddons] = useState(new Map())
    const [availableAddons, setAvailableAddons] = useState(new Map())
    const [state, setState] = useState(states.Loading)


    function fetchAvailableAddonList() {
        return new Promise(function (resolve, reject) {
            API.getAddonsInfo().then((data) => {
                console.log(data)
                if (!data || !data.urls || !data.architecture || !data.version ||
                    !data.nodeVersion) {
                    return;
                }

                const params = new URLSearchParams();
                params.set('arch', data.architecture);
                params.set('version', data.version);
                params.set('node', data.nodeVersion);

                console.log(data)

                if (data.pythonVersions && data.pythonVersions.length > 0) {
                    params.set('python', data.pythonVersions.join(','));
                }

                if (data.testAddons) {
                    params.set('test', '1');
                }
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
                return Promise.all(promises)
            }).then(
                (responses) => {
                    const promises = []
                    console.log(responses, "responses.....")
                    for (const resp of responses) {
                        promises.push(resp.json())
                    }
                    return Promise.all(promises)
                }).then((bodies) => {
                    let availableAddons = new Map()
                    for (const body of bodies) {
                        if (!body) {
                            continue
                        }
                        for (const addon of body) {
                            const entry = {
                                id: addon.id,
                                name: addon.name,
                                description: addon.description,
                                author: addon.author,
                                homepage_url: addon.homepage_url,
                                license_url: addon.license_url,
                                version: addon.version,
                                url: addon.url,
                                checksum: addon.checksum,
                                primary_type: addon.primary_type,
                                installed: installedAddons.has(addon.id),
                            }

                            if (installedAddons.has(addon.id)) {
                                continue;
                            }
                            availableAddons.set(addon.id, entry)

                        }
                    }
                    resolve(availableAddons);
                }
            ).catch((e) => {
                reject(e)
            })

        })


    }

    function fetchInstalledAddonsList() {
        return new Promise(function (resolve, reject) {
            API.getInstalledAddons().then(body => {
                if (body === null) {
                    console.log(body)
                    return reject(new Error("installed empty"))
                }
                console.log("fetch the installed addon body :", body)
                let newMap = new Map()
                try {
                    for (const s of body) {
                        s.isUpdate = false
                        newMap.set(s.id, s)
                    }

                } catch (err) {
                    return reject(err);
                }
                return resolve(newMap)
            })
        })
    }

    function renderInstalledAddonsList() {
        console.log("installedAddons:", installedAddons)
        const list = []
        for (const [id, a] of installedAddons) {
            console.log("installedAddons:", id, a)
            const addon = <InstalledAddon key={a.id}
                                          id={a.id}
                                          name={a.name}
                                          short_name={a.short_name}
                                          author={a.author}
                                          description={a.description}
                                          enabled={a.enabled}
                                          exec={a.exec}
                                          license={a.license_url}
                                          version={a.version}
                                          primary_type={a.primary_type}
                                          schema={a.schema}
                                          isUpdate={a.isUpdate}

            />
            list.push(addon)
        }
        return list
    }

    useEffect(() => {
        if (!props.open) {
            setFetchAddonsShow(false)
            setState(states.Empty)
            setInstalledAddons(new Map())
            setAvailableAddons(new Map())
        } else {
            setState(states.Loading)

            let installed = new Map()
            fetchInstalledAddonsList().then((installedMap) => {
                    installed = installedMap
                    return fetchAvailableAddonList()
                }, (e) => {
                    setState(states.Empty)
                    return fetchAvailableAddonList()
                }
            ).then((fetchAddons) => {
                for (const [id, addon] of fetchAddons) {
                    if (installed !== null) {
                        if (installed.has(id)) {
                            if (!addon.installed) {
                                fetchAddons.get(id).installed = true

                            }
                            if (versionStringCompare(addon.version, installed.get(id).version) > 0) {
                                console.log(versionStringCompare(addon.version, installed.get(id).version))
                                installed.get(id).isUpdate = true
                                installed.get(id).url = addon.url

                            }

                        }
                    }
                }
                if (fetchAddons) {
                    console.log("update available:", fetchAddons)
                    setAvailableAddons(fetchAddons)
                }
                if (installed !== null) {
                    console.log("update installed:", installed)
                    setInstalledAddons(installed)
                    setState(states.Completed)
                }
            }, (err) => {
                console.log(err)
                setState(states.Completed)
                setInstalledAddons(installed)
            }).catch((e) =>
                console.log(e)
            )
        }

    }, [props.open])

    useEffect(() => {
        console.info("update availableAddons:")
    }, [availableAddons])

    useEffect(() => {
        console.info("update installedAddons:", installedAddons)
        if (installedAddons.size === 0) {
            setState(states.Empty)
        } else {
            setState(states.Completed)
        }
    }, [installedAddons])

    return (
        <div>
            <Dialog fullScreen open={props.open} onClose={() => props.show(false)}
                    TransitionComponent={Transition}>
                <AppBar className={classes.appBar}>
                    <Toolbar>
                        <ExtensionIcon/>
                        <Typography variant="h6" className={classes.title}>
                            {t("Addons")}
                        </Typography>
                        <IconButton autoFocus color="inherit" onClick={() => props.show(false)}
                                    aria-label="close">
                            <CloseIcon/>
                        </IconButton>
                        <IconButton autoFocus color="inherit" aria-label="close" onClick={() => {
                            setFetchAddonsShow(true)
                        }}>
                            <AddCircleIcon/>
                        </IconButton>
                    </Toolbar>
                </AppBar>

                <Grid className={classes.content} container justify="flex-start" alignItems="center" direction="column">
                    <div className={classes.drawerHeader}/>
                    {state === states.Loading && <CircularProgress disableShrink/>}
                    {state === states.Completed && renderInstalledAddonsList()}
                    {state === states.Empty &&
                    <div style={{display: "flex", flexDirection: "column", alignItems: "center"}}><AddIcon cursor={"pointer"}
                        onClick={() => setFetchAddonsShow(true)} style={{fontSize: 50}}/><Typography variant="h6"
                                                                                                     className={classes.title}>
                        {t("Click Added")}
                    </Typography></div>}

                </Grid>
            </Dialog>
            <AddAddonsDialog installedAddons={installedAddons} availableAddons={availableAddons}
                             open={fetchAddonsShow} show={setFetchAddonsShow}/>
        </div>
    )

}

export function AddAddonsDialog(props) {

    const classes = useStyles();
    const {t, i18n} = useTranslation();

    const availableAddons = props.availableAddons

    function renderAvailableAddons() {
        if (!availableAddons) {
            return
        }
        const list = []
        for (const [id, addon] of availableAddons) {
            if (availableAddons.has(id)) {
                const a = <NewAddon key={addon.id} id={addon.id}
                                    name={addon.name}
                                    author={addon.author}
                                    description={addon.description}
                                    homepage_url={addon.homepage_url}
                                    license_url={addon.license_url}
                                    url={addon.url}
                                    version={addon.version}
                                    primary_type={addon.primary_type}
                                    checksum={addon.checksum}
                                    installed={addon.installed}

                />
                list.push(a)
            }

        }
        return list
    }


    return (
        <Dialog fullScreen className={classes.root} open={props.open} onClose={() => props.show(false)}
                TransitionComponent={Transition}>
            <AppBar className={classes.appBar}>
                <Toolbar>
                    <StoreIcon/>
                    <Typography variant="h6" className={classes.title}>
                        {t("AddonsMarket")}
                    </Typography>
                    <IconButton autoFocus color="inherit" onClick={() => props.show(false)} aria-label="close">
                        <CloseIcon/>
                    </IconButton>
                </Toolbar>
            </AppBar>
            <div className={classes.drawerHeader}/>
            <Grid className={classes.content} container justify="flex-start" alignItems="center" direction="column">
                {renderAvailableAddons()}
            </Grid>
        </Dialog>
    )

}

function InstalledAddon(props) {

    const classes = useStyles();
    const {t, i18n} = useTranslation();

    const [addon, setAddon] = useState(props)

    return (
        <>
            <Card className={classes.addonCard} elevation={5}>

                <Addon {...props}/>
                <div className={classes.sideContent}>
                    {!addon.enabled && <Button style={{margin: 3}}
                                               variant="contained"
                                               color="primary">
                        {t("enable")}
                    </Button>}

                    {addon.enabled && <Button style={{margin: 3}}
                                              variant="contained"
                                              color="primary">
                        {t("disable")}
                    </Button>}

                    {addon.isUpdate && <Button style={{margin: 3}}
                                               variant="contained"
                                               color="primary">
                        {t("update")}
                    </Button>}

                    <Button style={{margin: 3}}
                            variant="contained"
                            color="primary">
                        {t("configure")}
                    </Button>
                </div>
            </Card>
            <Divider/>
        </>
    )

}

function NewAddon(props) {

    const classes = useStyles();
    const {t, i18n} = useTranslation();


    const states = {
        install: "install",
        installed: "installed",
        pending: "pending",
        failed: "failed",
    }
    const [state, setState] = useState()

    useEffect(() => {
        if (props.installed) {
            setState(states.installed)
        } else {
            setState(states.install)
        }
    }, [])


    const handleInstallAddon = () => {
        setState(states.pending)
        console.log("install args:", props.id, props.url, props.checksum)
        API.installAddon(props.id, props.url, props.checksum).then((req) => {
            console.log("install addon OK", req)
            setState(states.installed)
        }).catch((e) => {
            setState(states.failed)
            console.error(e)
        })
    }


    return (
        <>
            <Card className={classes.addonCard} elevation={10}>
                <Addon {...props} />

                <div className={classes.sideContent}>
                    {
                        state === states.install && <Button
                            onClick={() => {
                                handleInstallAddon()
                            }}
                            variant="contained"
                            color="primary">
                            {t("install")}
                        </Button>
                    }
                    {
                        state === states.installed && <Button
                            variant="contained"
                            disabled={true}
                            color="primary"
                        >
                            {t("installed")}
                        </Button>
                    }
                    {
                        state === states.pending && <CircularProgress/>
                    }

                    {
                        state === states.failed && <Button
                            color="secondary"
                            disabled={true}
                            className={classes.button}
                        >
                            <Typography gutterBottom color="secondary">
                                {t(states.failed)}
                            </Typography></Button>
                    }
                </div>
            </Card>
            <Divider/>
        </>
    )

}

export function Addon(props) {

    const classes = useStyles();
    const {t, i18n} = useTranslation();
    const addon = props
    return <>
        <Icons style={{fontSize: 80}} color={"primary"} type={addon.primary_type}/>
        <div className={classes.cardContent}>
            <FormHelperText>{t("name")}: </FormHelperText>
            <Typography variant="h5">
                {addon.name}
            </Typography>
            <FormHelperText>{t("id")}: </FormHelperText>
            <Typography variant="subtitle1">
                {addon.id}
            </Typography>
            <div style={{
                display: "flex",
                "flexDirection": "row",
                "justifyContent": "flex-start",
                "alignItems": "center"
            }}>
                <FormHelperText>{t("version")}: </FormHelperText>
                <Typography variant="body1">
                    {addon.version}
                </Typography>
            </div>
            <div style={{
                display: "flex",
                "flexDirection": "row",
                "justifyContent": "flex-start",
                "alignItems": "center"
            }}>
                <FormHelperText>{t("author")}: </FormHelperText>

                <Link href={addon.homepage_url}>
                    {addon.author}
                </Link>

            </div>
            <div style={{display: "flex", "flexDirection": "row"}}>

                <Link href={addon.license_url}>
                    {t("license")}
                </Link>

            </div>

        </div>
    </>

}
