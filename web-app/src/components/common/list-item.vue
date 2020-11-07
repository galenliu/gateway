<template>
  <div v-if="visible" class="scroll-item" :class="['scroll-item', {'selected':selected}]">
    <a @click="OnClick">
      <component class="item-icon" v-bind:is="icon"/>
      <span v-if="unfold" id="item-title">{{ title }} </span>
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

export default {

  name: "list-item",

  props: {

    visible: {
      type: Boolean,
      default: true
    },

    selected: {
      type: Boolean,
      default: false,
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

    console.log(props)
    const ScrollItemReactiveData = reactive({

      icon: props.icon,
      title: props.title,
      visible: props.visible,
      selected: props.selected,
    })
    const unfold = inject("sidebar-unfold")
    const OnClick = (e) => {
      context.emit("OnClick:ItemSelected",e)
    }
    watchEffect(() => {

      console.log(ScrollItemReactiveData.selected)

    })


     const doSelected = function() {
        ScrollItemReactiveData.selected= !props.selected

    };


    return {
      unfold,
      OnClick,
      ...toRefs(ScrollItemReactiveData)
    }
  }
}
</script>

<style scoped>
.scroll-item > a {
  display: flex;
  align-items: center;
  margin: 4px;
}

.selected {
  border-radius: 12px;
  background-color: lightblue;
}

a {
  cursor: pointer;
}

.item-icon {
  margin-left: 4px;
  margin-right: 8px;
  padding: 8px 12px;
}
</style>
