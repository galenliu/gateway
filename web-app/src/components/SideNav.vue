<template>
  <aside class="sidebar">
    <main-menu @menuClick= toggleSidebar title="刘桂林的家庭世界"></main-menu>
    <scroll-listbox></scroll-listbox>
  </aside>
</template>

<script lang="ts">
import MainMenu from "./sidebar/main-menu.vue"
import ScrollListbox from "./sidebar/scroll-listbox.vue"
import { provide , ref } from 'vue'
export default {
  name: "side-nav",

  components: {
    ScrollListbox,
    MainMenu,
  },
  setup() {

    const unfold = ref(true)
    provide("sidebar-unfold", unfold)

    const toggleSidebar = function (){
      console.log("menu click")
      unfold.value = !unfold.value
      console.log(unfold.value)
    }

    return {
      toggleSidebar,
      unfold,
    }
  }
}
</script>

<style scoped>

* {
  margin: 0;
  padding: 0;
}

aside {
  position: fixed;
  height: 100%;
  display: flex;
  flex-direction: column;
  background-color: #ffffff;
  border-right: solid 1px rgba(0,0,0,0.12);
}


</style>