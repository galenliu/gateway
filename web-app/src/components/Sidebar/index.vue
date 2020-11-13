<template>
  <aside class="sidebar">
    <main-menu title="刘桂林的家庭世界"></main-menu>
    <scroll-box>
      <list-item v-for="(value,key,i) in items" :key=i :ref="el => {if(el) navItems[i] = el}" :name="key"
                 :class="{'selected': selectedKey === i }"
                 :title=value.title :icon="value.icon" @ItemSelected="onSelected(i,key,$event)">
      </list-item>
    </scroll-box>
  </aside>
</template>

<script lang="ts">
import MainMenu from "./MainMenu.vue"
import ScrollBox from "./ScrollBox.vue"
import ListItem from "./ListItem.vue"
import {provide, ref, inject, reactive, toRefs} from 'vue'
import {store} from "../../utils/Store.ts"
import {getMenuList} from "../../utils/Utils.ts"
import {router} from "../../router";


export default {
  name: "Sidebar",

  components: {
    ScrollBox,
    ListItem,
    MainMenu,
  },

  props: {
    activeKey: {
      type: Number,
      default: 0
    },
    items: {
      type: Object,
      require: true,
    }
  },

  setup(props,context) {

    const selectedKey= ref(props.activeKey)
    const navItems = ref<ListItem[]>([])

    //侧边栏有选项被选中时，执行方法，key表示第几项
    const onSelected = (key:number,name:String,e) => {
      selectedKey.value = key
      context.emit("goto",name)

    }

    return {
      store,
      navItems,
      selectedKey,
      onSelected,
      ...toRefs(store),
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