<template>
  <layout>
    <template v-slot:sidebar>
      <sidebar >
      </sidebar>
    </template>
    <template v-slot:top-bar>
      <top-bar></top-bar>
    </template>
      <router-view v-slot: work-space>
      </router-view>
  </layout>

</template>

<script lang="ts">

import TopBar from '../components/Topbar.vue'
import Sidebar from "../components/Sidebar/index.vue"



import Layout from "./Layout.vue";

export default {
  name: "index",



  components: {
    Layout,
    TopBar,
    Sidebar,
  },

  setup() {
    const itemSelected=function (){
        console.log("1111111111111")
    }

  }


}
</script>

<style scoped lang="scss">

</style>