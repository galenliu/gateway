import { createApp } from 'vue'
import App from './App.vue'
import './assets/scss/index.scss'
import '@icon-park/vue-next/styles/index.css';

import { router } from "./router"

createApp(App)
    .use(router)
    .mount('#app')
