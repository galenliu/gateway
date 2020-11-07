// @ts-ignore
import {createRouter, createWebHistory , createWebHashHistory} from "vue-router";

const history = createWebHistory()
const hashRouter = createWebHashHistory()

import Home from "./views/Home.vue"
import Overview from './views/Overview.vue'
import Rule from "./views/Rule.vue"
import Profile from "./views/Profile.vue"
import Login from "./views/Login.vue"

const router = createRouter( {
    history: hashRouter,
    linkActiveClass:'link-active',
    routes: [
        {
            path: '/',
            component: Home,
            children: [
                {path: 'home', component: Home},
                {path: 'overview', component: Overview},
                {path: 'Rule', component: Rule},
                {path: 'Profile', component: Profile},
                {path: 'Login', component: Login},
            ]
        }
    ]
})

export {
    router
}