<template>
  <div id="main-menu">
    <menu-unfold v-if="!state.sideUnfold" class="menu-toggle-button" theme="outline"  size="24px"  fill="#333" @click=menuClick  />
    <menu-fold v-if="state.sideUnfold" class="menu-toggle-button" theme="outline"  size="24px"  fill="#333" @click=menuClick  />
    <span v-if="state.sideUnfold" id="menu-btn-title"> {{ title }}</span>
  </div>
</template>

<script>



import { MenuFold } from '@icon-park/vue-next'
import { MenuUnfold } from '@icon-park/vue-next'
import {reactive, toRefs ,inject } from "vue";

import { store } from "../../utils/Store"

export default {
  name: "MainMenu",

  components: {
    MenuFold,
    MenuUnfold,
  },

  props: {
    title: {
      type: String,
      default: 'WebThings'
    },
  },

  setup(props,context) {
    const menuReactiveData = reactive({
      title: props.title,

    })

    const menuClick =() => {
      store.setSideUnfoldAction(!store.state.sideUnfold)
    }

    return {
      ...toRefs(store),
      menuClick,
      ...toRefs(menuReactiveData)
    }
  }
}



</script>

<style lang="scss" scoped>

$menu-color: #fafafa;


#main-menu {
  border-bottom: solid 1px rgb(0,0,0,.12);
  background-color: $menu-color;

}

.menu-toggle-button{
  background-color: transparent;
  margin: 8px;
  padding: 12px;
  &:hover{
    background-color: rgb(0,0,0,.1);
    border-radius: 50%;
  }
}

#menu-btn-title{

  margin-right: 1em;
}



</style>