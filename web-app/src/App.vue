<template>
<router-view></router-view>
</template>

<script lang="ts">

import {topMenuStatus} from "./utils/Utils.ts"

import { ref, provide, reactive} from "vue"

import NProgress from 'nprogress'
import 'nprogress/nprogress.css'
import {router} from "./router";
import { getMenuList } from "./utils/Utils.ts"

export default {

  name: 'App',

  setup(){
    const topBarMenuVisible = ref(topMenuStatus())
    const width = document.documentElement.clientWidth
    const sidebarState =  reactive({
      topBarMenuVisible: true
    })
    provide("topBarMenuVisible",topBarMenuVisible)
    provide("sideMenuList",getMenuList())
    provide("sideState",sidebarState)

    router.beforeEach((to, from, next) => {
      NProgress.start()
      next()
    })

    router.afterEach(() => {
      if(width<=900){
        topBarMenuVisible.value = false
      }
      NProgress.done()
      window.scrollTo(0, 0)
    })

  },

  components: {


  }
}
</script>

<style lang="scss">
@import "assets/scss/var.scss";

#nprogress {
  .bar {
    background: $theme !important; //自定义颜色
  }

  .spinner-icon {
    border-color: $theme transparent transparent $theme;
  }
}

</style>