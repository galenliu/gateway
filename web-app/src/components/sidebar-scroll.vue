<template>
  <div class="sidebar-scroll">
    <sidebarItem v-for="(value, name) in getListItem" :item_style=name :title=value></sidebarItem>
  </div>
</template>

<script lang="ts">
import sidebarItem from './sidebar-item.vue'
import {computed, reactive, toRefs} from 'vue'

export default {
  name: "sidebar-scroll",

  components: {
    sidebarItem: sidebarItem,
  },

  setup() {
    const scrollReactiveData = reactive({})

    const getListItem = computed(getItem)

    return {
      getListItem,
      ...toRefs(scrollReactiveData)
    }
  }
}

function getItem(): object {
  return {
    home: "家庭",
    rule: "自动化",
    settings: "设置",
    exit: "退出",
  }
}
</script>

<style scoped>
.sidebar-scroll {
  display: flex;
  flex-direction: column;
}
</style>