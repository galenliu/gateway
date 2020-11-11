import {CheckCorrect, Home, SettingConfig} from "@icon-park/vue-next";

const topMenuStatus = function (): Boolean {
    return document.body.clientWidth > 900
}

const getMenuList = function (): Object {

    const items = {
        "home": {
            title: "家庭",
            icon: Home,
            selected: true
        },
        "rule": {

            title: "规则",
            icon: CheckCorrect,
            selected: false,
        },
        "profile": {
            title: "配置",
            icon: SettingConfig,
            selected: false
        },
    }

    return items
}

export {topMenuStatus, getMenuList}