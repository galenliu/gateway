class API {

    public jwt: string;

    constructor() {
        this.jwt = localStorage.getItem('jwt')
    }

    //region  json handler

    headers(contentType?: string): object {
        const headers = {
            Accept: 'application/json',
            Authorization: null
        };
        if (this.jwt) {
            headers.Authorization = 'Bearer ${ this.jwt }';
        }
        if (contentType) {
            headers['Content-Type'] = contentType
        }
        return headers
    }


    async getJson(url: string): Promise<string | null> {
        await fetch(url).then(res => {
            if (!res.ok) {
                throw new Error(res.status.toString())
            }
            return res.json()
        })
        return
    }

    async postJson(url, data: object): Promise<string | null> {
        const opts = {
            method: "POST",
            header: this.headers("application/json"),
            body: JSON.stringify(data),
        };
        await fetch(url, opts).then(res => {
            if (!res.ok) {
                throw new Error(res.status.toString())
            }
            return res.json()
        })
        return
    }


    async putJson(url, data: object): Promise<string | null> {
        const opts = {
            method: "PUT",
            header: this.headers("application/json"),
            body: JSON.stringify(data),
        };
        await fetch(url, opts).then(res => {
            if (!res.ok) {
                throw new Error(res.status.toString())
            }
            return res.json()
        })
        return
    }

    async patchJson(url, data: Promise<string | null>) {
        const opts = {
            method: "PATCH",
            header: this.headers("application/json"),
            body: JSON.stringify(data),
        };
        await fetch(url, opts).then(res => {
            if (!res.ok) {
                throw new Error(res.status.toString())
            }
            return res.json()
        })
        return
    }

    async delete(url: string) {

        const opts = {
            method: "DELETE",
            header: this.headers(),
        };
        await fetch(url, opts).then(res => {
            if (!res.ok) {
                throw new Error(res.status.toString())
            }
        })
    }

    //endregion

    //region  user handler

    getUser(id: string): string {
        return ""
    }

    addUser(name, email, password: string): string {
        return ""
    }

    editUser(id, name, email, password, newPassword: string): string {
        return ""
    }

    deleteUser(id: string): string {
        return ""
    }

    login(email, password, totp: string): string {
        return ""
    }

    logout(): string {
        return ""
    }

    //endregion

    //region addons handle

     getInstalledAddons(): Promise<string> {
        return  this.getJson('/addons')
    }

     getAddonConfig(addonId: string): Promise<string> {
        return  this.getJson('/addons/' + addonId + '/config')
    }

    setAddonSetting(addonId: string, enabled: boolean): Promise<string> {
        return  this.putJson('/addons/${ addonId }/config', {enabled})
    }

    installAddon(addonId, addonUrl, addonChecksum: string): Promise<string> {
        return this.postJson('/addons', {
            id: addonId,
            url: addonUrl,
            checksum: addonChecksum
        })
    }

    uninstallAddon(addonId: string) {
        return this.delete('/addons/${ addonId }')
    }

    getAddonsInfo(): Promise<string | null> {
        return this.getJson("/settings/addonsInfo");
    }

    //endregion

    //region things handler
    getThings(): string {
        return ""
    }

    getThing(thingId: string): string {
        return ""
    }

    setThingCredentials(thingId: string, data: string): string {
        return ""
    }

    removeThing(thingId: string) {

    }

    //endregion

}