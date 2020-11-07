<template>
  <router-view/>
</template>

<script lang="ts">


import {topNavMenuStatus} from "./utils/utils.ts"

import { ref, provide } from "vue"

import NProgress from 'nprogress'
import 'nprogress/nprogress.css'
import {router} from "./router";

export default {

  name: 'App',

  setup(){
    const topBarMenuVisible = ref(topNavMenuStatus())
    const width = document.documentElement.clientWidth
    provide("topNavMenuStatus",topNavMenuStatus)

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