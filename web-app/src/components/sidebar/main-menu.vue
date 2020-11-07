<template>
  <div id="main-menu">
    <menu-unfold v-if="!unfold" class="menu-toggle-button" theme="outline"  size="24px"  fill="#333" @click=menuClick  />
    <menu-fold v-if="unfold" class="menu-toggle-button" theme="outline"  size="24px"  fill="#333" @click=menuClick  />
    <span v-if="unfold" id="menu-btn-title"> {{ title }}</span>
  </div>
</template>

<script>



import { MenuFold } from '@icon-park/vue-next'
import { MenuUnfold } from '@icon-park/vue-next'
import {reactive, toRefs ,inject } from "vue";

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
       context.emit("menuClick")
    }

    const unfold =inject("sidebar-unfold")
    return {
      menuClick,
      unfold,
      ...toRefs(menuReactiveData)
    }
  }
}



</script>

<style lang="scss" scoped>

$menu-color: #fafafa;

* {
  margin: 0;
  padding: 0;
}

#main-menu {
  display: inline-flex;
  box-sizing: border-box;
  align-items: center;
  border-bottom: solid 1px rgb(0,0,0,.12);
  background-color: $menu-color;
  overflow: hidden;
  z-index: 10;
}

.menu-toggle-button{
  background-color: transparent;
  margin: 8px;
  padding: 12px;
  &:hover{
    background-color: rgb(0,0,0,.1);
    border-radius: 50%;
  }

  &+span{
    font-size: 1em;
    margin-right: 1em;
  }
}



</style>