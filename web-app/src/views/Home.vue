<template>
  <h1>Home</h1>
  <router-link class:="link-a" to="/overview">qwerqwe
    <button>123123412123</button>
  </router-link>
  <div class="router-view-wrapper">
    <router-view></router-view>
  </div>
</template>

<script lang="ts">

import TopNav from '../components/TopNav.vue'
import SideNav from "../components/SideNav.vue";

import {ref, inject} from 'vue'

export default {
  name: "Home",


  components: {
    TopNav,
    SideNav,
  },

  setup() {

    const topBarMenuVisible = inject<ref<boolean>>('topBarMenuVisible')

    return {

      topBarMenuVisible
    }

  }


}
</script>

<style scoped lang="scss">
div {
  width: 100px;
  height: 100px;
  background: linear-gradient(145deg, rgb(239, 228, 255) 0%, rgb(206, 187, 232) 100%);
}
</style>