import { createApp } from 'vue'
import './index.css'
import App from './App.vue'

import {createWebHashHistory, createRouter} from './node_modules/vue-router'

const history = createWebHashHistory()
const router = createRouter({
    history:history,
    routes: [
       { path:"/",component: App}
    ]
})

const app = createApp(App)
app.use(router)
app.mount('#app')



