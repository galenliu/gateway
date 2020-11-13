<template>
  <layout>
    <template v-slot:sidebar>
      <sidebar :items="items" @goto="goto">
      </sidebar>
    </template>
    <template v-slot:top-bar>
      <top-bar></top-bar>
    </template>
    <template v-slot:work-space>
      <router-view>
      </router-view>
    </template>
  </layout>

</template>

<script lang="ts">

import TopBar from '../components/Topbar.vue'
import Sidebar from "../components/Sidebar/index.vue"


import Layout from "./Layout.vue";
import { reactive } from "vue";
import {getMenuList} from "../utils/Utils";

import "vue-router";
import {router} from "../router";


export default {
  name: "index",


  components: {

    Layout,
    TopBar,
    Sidebar,
  },

  setup() {

    const items = getMenuList()

    const goto= (name:String) =>{
      router.push(name)
    }
    return{
      goto,
      items,
    }
  }

}
</script>

<style scoped lang="scss">

</style>