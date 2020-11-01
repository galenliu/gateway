<template>
  <div class="svg-icon container">
    {{ iconFile }}
    <div></div>
  </div>
</template>

<script lang="ts">

//import path from "path";
import path from "path"


import {computed, reactive, toRefs,} from 'vue'


export default {
  name: "svgIcon",
  props: {
    iconSize: {
      type: Number,
      default: 24
    },

    iconName: {
      type: String,

    }
  },

  setup(props) {


    const iconData = reactive(
        {
          iconName: props.iconName,
          iconSize: props.iconSize
        }
    );

    const iconFile = computed(() => {

      let file = path.join( __dirname, iconData.iconName + "-" + iconData.iconSize + "px.svg")

      console.log(file)
      return file
    })

    return {
      iconFile,
      ...toRefs(iconData),
    }

  }
}


</script>


<style scoped>
div {
  width: 100px;
  height: 100px;
  background-color: red;
}
</style>
