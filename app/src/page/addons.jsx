import React, {useCallback, useEffect, useState} from "react";
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
import ThingIcon from "../component/thing-icon";
import {Button, CircularProgress, FormHelperText, Link} from "@material-ui/core";
import Divider from "@material-ui/core/Divider";
import StoreIcon from "@material-ui/icons/Store";
import {versionStringCompare} from "../js/util";


export const useStyles = makeStyles((theme) => ({
    root: {},
    addonCard: {
        padding: 3,
        display: "flex",
        maxWidth: 400,
        minWidth: 300,
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
        marginTop: 80,
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


export default function AddonsDialog(props) {

    const classes = useStyles();
    const {t, i18n} = useTranslation();
    const [fetchAddonsShow, setFetchAddonsShow] = useState(false)
    const [installedAddons, setInstalledAddons] = useState(new Map())
    const [availableAddons, setAvailableAddons] = useState(new Map())

    const [open, setOpen] = useState(props.open)

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
        return API.getInstalledAddons().then(body => {
            if (!body) {
                return;
            }

            console.log("fetch the installed addon body :", body)
            let list = new Map()
            for (const s of body) {
                console.log("fetch addon  :", s)
                try {
                    s.isUpdate = false
                    list.set(s.id, s)
                } catch (err) {
                    console.error(`Failed to parse add-on settings: ${err}`);
                }
            }
            setInstalledAddons(list)
        })
    }


    function renderInstalledAddonsList() {
        const list = []
        for (const [id, a] of installedAddons) {
            const addon = <InstalledAddon key={id}
                                          id={id}
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
                                          isUpdate = {a.isUpdate}

            />
            list.push(addon)
        }
        return list
    }

    useEffect(() => {
            if (!open) {
                return
            }
            fetchInstalledAddonsList().then(() =>
                fetchAvailableAddonList()
            ).then((fetchAddons) => {
                console.log("fetchAddons:", fetchAddons)
                setAvailableAddons(fetchAddons)
            }).catch((e) =>
                console.error(e)
            )

        }, [open])

    useEffect(() => {
        let availableUpdate = false
        let installedUpdate = false
        let available = availableAddons
        let installed = installedAddons

        console.log("available:", available)
        for (const [id, addon] of availableAddons) {
            if (installedAddons.has(id)) {
                if (!addon.installed) {
                    available.get(id).installed = true
                    availableUpdate = true
                }
                if (versionStringCompare(addon.version, installed.get(id).version) === 0) {
                    console.log(versionStringCompare(addon.version, installed.get(id).version))
                    installed.get(id).isUpdate = true
                    installed.get(id).url = addon.url
                    installedUpdate = true
                }

            }
            if (availableUpdate) {
                setAvailableAddons(available)
            }
            if (installedUpdate) {
                setInstalledAddons(installed)
            }
            console.log("available:", availableAddons)
            console.log("installed:", installedAddons)
        }
    }, [availableAddons])


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
                <Grid container justify="flex-start" alignItems="center" direction="column">
                    {renderInstalledAddonsList()}
                </Grid>
            </Dialog>
            {fetchAddonsShow && <AddAddonsDialog installedAddons={installedAddons} availableAddons={availableAddons}
                                                 open={fetchAddonsShow} show={setFetchAddonsShow}/>}
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
        console.log("addon:", availableAddons)
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

    useEffect(() => {
        if (props.open) {

        }
    }, [props.open])

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
            <Card className={classes.addonCard} elevation={5}>
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
        <ThingIcon style={{fontSize: 80}} color={"primary"} type={addon.primary_type}/>
        <div className={classes.cardContent}>
            <FormHelperText>{t("name")}: </FormHelperText>
            <Typography variant="h5">
                {addon.name}
            </Typography>
            <FormHelperText>{t("id")}: </FormHelperText>
            <Typography variant="subtitle1">
                {addon.id}
            </Typography>
            <div style={{display: "flex", "flex-direction": "row"}}>
                <FormHelperText>{t("version")}: </FormHelperText>
                <Typography variant="h9">
                    {addon.version}
                </Typography>
            </div>
            <div style={{display: "flex", "flex-direction": "row"}}>
                <FormHelperText>{t("author")}: </FormHelperText>
                <Typography variant="h9">
                    <Link href={addon.homepage_url}>
                        {addon.author}
                    </Link>
                </Typography>
            </div>
            <div style={{display: "flex", "flex-direction": "row"}}>
                <Typography variant="h9">
                    <Link href={addon.license_url}>
                        {t("license")}
                    </Link>
                </Typography>
            </div>

        </div>
    </>

}
