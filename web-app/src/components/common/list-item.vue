<template>
  <div v-if="visible" class="scroll-item" :class="{'selected':selected}">
    <a>
      <slot name="icon"></slot>
      <span v-if="unfold" id="item-title">{{ title }} </span>
    </a>
  </div>
</template>

<script>
import {inject, reactive, toRefs} from 'vue'

export default {

  name: "list-item",

  props: {
    visible: {
      type: Boolean,
      default: true
    },

    selected: {
      type: Boolean,
      default: true,
    },

    title: {
      type: String,
      required: true
    }
  },

  setup(props) {

    console.log(props)
    const ScrollItemReactiveData = reactive({
      title: props.title,
      visible: props.visible,
      selected: props.selected,
    })
    const unfold = inject("sidebar-unfold")

    return {
      unfold,
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
.selected a{
  border-radius: 12px;
  background-color: lightblue;
}

</style>