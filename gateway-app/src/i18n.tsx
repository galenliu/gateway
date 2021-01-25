import i18n from "i18next";
import LanguageDetector from 'i18next-browser-languagedetector';
import enUsTrans from "./i18n/en-us.json";
import zhCnTrans from "./i18n/zh-cn.json";
import {initReactI18next} from 'react-i18next';

i18n.use(LanguageDetector) //获取当前浏览器语言
    .use(initReactI18next)
    .init({
        //资源文件
        resources: {
            enUs: {
                translation: enUsTrans,
            },
            zhCn: {
                translation: zhCnTrans,
            },
        },
        fallbackLng: "zhCn",
        debug: false,
        interpolation: {
            escapeValue: false,
        },
    })

export default i18n;
