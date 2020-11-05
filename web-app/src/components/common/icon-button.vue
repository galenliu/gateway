<template>
  <div class="wt-icon-button" :class="[`wt-icon-button--${type)}`]">
    {{svgFile}}
    <svg :src='svgFile'></svg>
  </div>

</template>

<script>
import {computed, reactive, toRefs} from "../../../../web-app1/node_modules/vue"
import path from '../../../../web-app1/node_modules/path'

const resDir = "/src/assets/svg"

export default {
  name: "icon-button",
  props: {
    type: {
      type: String,
      required: true,
    },
    svgSize: {
      type: Number,

    },
    iconSize:{
      type: Number,
          default: 48,
}
  },

  setup(props) {

    const reactiveDate = reactive({
      svgFile,
      svgSize: props.svgSize,

    })

    const svgFile = computed(() => {
      const file = path.join(resDir, iconData.btnStyle + "-" + reactiveDate.svgSize + "px.svg")
      console.log(file)
      return file
    })
    return {
      svgFile,
      ...toRefs(reactiveDate)
    }
  }
}


</script>


<style scoped>

</style>