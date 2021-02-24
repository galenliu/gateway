import i18n from "i18next";
import LanguageDetector from 'i18next-browser-languagedetector';
import {initReactI18next} from "react-i18next";
import enTrans from "./i18n/en-us.json"
import zhTrans from "./i18n/zh-cn.json"


// the translations
// (tip move them in a JSON file and import them)
const resources = {
    en: {
        translation: enTrans,
    },
    zh: {
        translation: zhTrans,
    }
};

i18n.use(LanguageDetector) //嗅探当前浏览器语言
    .use(initReactI18next) // passes i18n down to react-i18next
    .init({
        resources,
        lng: "zh",

        keySeparator: false, // we do not use keys in form messages.welcome

        interpolation: {
            escapeValue: false // react already safes from xss
        }
    }).catch((e) => {
    console.log(e)
});

export default i18n;
