// @ts-ignore
import {createRouter, createWebHistory , createWebHashHistory} from "vue-router";

const history = createWebHistory()
const hashRouter = createWebHashHistory()

import index from "./views/index.vue"
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
            component: index,
            children: [
                {path: 'home', component: index},
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