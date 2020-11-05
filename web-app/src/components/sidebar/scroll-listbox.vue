<template>
  <div class="sidebar-scroll">

    <list-item v-for="(item,i) in items" v-bind:title="item.title">
      <template v-if="item.name === 'home'"  v-slot:icon>
        <Home size="24px" class="item-icon" />
      </template>

      <template v-if="item.name === 'rule'" v-slot:icon>
        <CheckCorrect size="24px" class="item-icon"/>
      </template>

      <template v-if="item.name === 'config'" v-slot:icon>
        <SettingConfig size="24px" class="item-icon"/>
      </template>


    </list-item>

  </div>
</template>

<script lang="ts">

import {reactive, toRefs, onMounted} from 'vue'
import ListItem from '../common/list-item.vue'

import {CheckCorrect, Home, SettingConfig} from "@icon-park/vue-next"

const items = [
  {name: "home",  title: "家庭"},
  {name: "rule", title: "规则"},
  {name: "config", title: "配置"}
]

export default {
  name: "ScrollListbox",

  components: {
    Home,
    SettingConfig,
    CheckCorrect,

    ListItem,
  },

  setup() {
    const scrollReactiveData = reactive({
      items,
    })


    return {
      ...toRefs(scrollReactiveData)
    }
  }
}


</script>


<style lang="scss" scoped>
.sidebar-scroll {
  display: flex;
  flex-direction: column;
  box-sizing: border-box;
}

.item-icon {
  margin-left: 4px;
  margin-right: 8px;
  padding: 8px 12px;
}
</style>