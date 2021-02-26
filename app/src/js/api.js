let debug = true;

export function getUrl(url) {
    if (debug) {
        return "http://localhost:9090" + url
    } else {
        return url
    }
}

export function getWsUrl(url) {
    if (debug) {
        return "ws://localhost:9090" + url
    } else {
        return url
    }
}

const API = {

    jwt: localStorage.getItem('jwt'),

    headers(contentType) {
        const headers = {
            Accept: 'application/json',
        };

        if (this.jwt) {
            headers.Authorization = `Bearer ${this.jwt}`;
        }

        if (contentType) {
            headers['Content-Type'] = contentType;
        }

        return headers;
    },

    getJson(url) {
        const opts = {
            method: 'GET',
            headers: this.headers(),
        };
        console.log("url:",url)
        return fetch(url, opts).then((res) => {
            console.log("fetch getJson api:", url)
            if (!res.ok) {
                throw new Error(res.status);
            }
            return res.json();
        });
    },

    putJson(url, data) {

        const opts = {
            method: 'PUT',
            headers: this.headers('application/json'),
            body: JSON.stringify(data),
        };

        return fetch(url, opts).then((res) => {
            if (!res.ok) {
                throw new Error(res.status);
            }
            return res.json();
        });
    },

    postJson(url, data) {

        const opts = {
            method: 'POST',
            headers: this.headers('application/json'),
            body: JSON.stringify(data),
        };
        return fetch(url, opts).then((res) => {
            if (!res.ok) {
                throw new Error(res.status);
            }
            if (res.status !== 204) {
                return res.json();
            }
        });
    },

    delete(url) {

        const opts = {
            method: 'DELETE',
            headers: this.headers(),
        };

        return fetch(url, opts).then((res) => {
            if (!res.ok) {
                throw new Error(res.status);
            }
        });
    },


    installAddon(addonId, addonUrl, addonChecksum) {
        return this.postJson("/addons", {
            id: addonId,
            url: addonUrl,
            checksum: addonChecksum,
        })
    },

    getAddonsInfo() {
        return this.getJson("/settings/addons_info")
    },

    getInstalledAddons() {
        return this.getJson("/addons")
    },

    startPairing(timeout) {
        return this.postJson('/actions', {
            pair: {
                input: {
                    timeout,
                },
            },
        });
    },
    cancelPairing(actionUrl) {
        return this.delete(actionUrl)
    },

    addThing(thing) {
        return this.postJson("/things", thing)
    },

    getThings() {
        return this.getJson("/things")
    },

    setThingPropertyValue(url, data) {
        return this.putJson(url, data)
    },


}

export default API
