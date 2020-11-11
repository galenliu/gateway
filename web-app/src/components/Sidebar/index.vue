<template>
  <aside class="sidebar">
    <main-menu  title="刘桂林的家庭世界"></main-menu>
    <scroll-box>
      <list-item v-for="(value,key,i) in items" :key=i :ref="key" :name= "key"
                 :selected="key == selectedItemName ? true: false" :title=value.title :icon="value.icon" @ItemSelected="selectedItem">
      </list-item>
    </scroll-box>
  </aside>
</template>

<script lang="ts">
import MainMenu from "./MainMenu.vue"
import ScrollBox from "./ScrollBox.vue"
import  ListItem  from "./ListItem.vue"
import {provide, ref, inject, reactive, toRefs }from 'vue'
import { store } from "../../utils/Store.ts"
import { getMenuList } from "../../utils/Utils.ts"



export default {
  name: "Sidebar",

  components: {
    ScrollBox,
    ListItem,
    MainMenu,
  },
  setup(props,context) {

    const selectedItemName = ref("name")

    const items = getMenuList()

    return {
      store,
      selectedItemName,
      ...toRefs(store),
      items,
    }
  }
}
</script>

<style scoped>

* {
  margin: 0;
  padding: 0;
}

.sidebar {
  height: 100%;
  display: flex;
  flex-direction: column;
  background-color: #ffffff;
  border-right: solid 1px rgba(0, 0, 0, 0.12);
}


</style>