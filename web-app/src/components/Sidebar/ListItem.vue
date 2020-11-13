<template>
  <div v-if="visible" class="scroll-item" >
    <a @click.prevent="OnClick" >
      <component class="item-icon" size="24px" v-bind:is="icon"/>
      <span v-if="store.state.sideUnfold" id="item-title">{{ title }} </span>
    </a>
  </div>
</template>

<script>
import {
  inject,
  reactive,
  toRefs,
  watchEffect,
} from 'vue'

import {
  Tag
} from "@icon-park/vue-next"
import { store } from "../../utils/Store"

export default {

  name: "ListItem",

  props: {
    name:{
      type: String,
    },

    visible: {
      type: Boolean,
      default: true
    },

    selected: {
      type: Boolean,
      default: true,
    },
    icon: {
      type: Object,
      default: Tag,
    },

    title: {
      type: String,
      required: true
    }
  },

  components: {},

  setup(props, context) {



    const OnClick = () => {
      context.emit("ItemSelected")
    }


    return {
      store,
      OnClick,

    }
  }
}
</script>

<style scoped lang="scss">
.scroll-item {

  display: flex;
  margin-bottom: 4px;
  margin-top: 4px;
  margin-right: 4px;
}

#item-title{
  display: inline-block;
}

.selected {
  cursor: pointer;
  border-radius: 12px;
  background-color: lightblue;
}

.scroll-item > a {
  display: flex;
  align-items: center;
  cursor: pointer;
}

.item-icon {
  display: flex;
  margin-left: 4px;
  margin-right: 8px;
  padding: 8px 12px;
}

</style>
