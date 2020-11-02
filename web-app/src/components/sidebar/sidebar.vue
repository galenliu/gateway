<template>
  <aside class="sidebar">
    <main-menu title="刘桂林的家庭自动化"></main-menu>
    <sidebarScroll></sidebarScroll>
  </aside>
</template>

<script lang="ts">
import mainMenu from "../sidebar/main-menu.vue"
import sidebarScroll from '../sidebar-scroll.vue'

export default {
  name: "sidebar",

  components: {
    mainMenu,
    sidebarScroll,
  },
  setup() {

    return {}
  }
}
</script>

<style scoped>

aside {
  position: fixed;
  display: flex;
  flex-direction: column;
  background-color: lavender;
}

* {
  margin: 0;
  padding: 0;
}
</style>