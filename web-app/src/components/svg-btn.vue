<template>
  <div id="svg-btn-container" class="svg-btn" :style="{ height : btnH , width : btnW}">
    <img :src='iconFile' alt="">
    <div></div>
  </div>
</template>

<script lang="ts">


import path from "path"

import {computed, reactive, toRefs,} from 'vue'

const resDir = "/src/assets/svg"

export default {
  name: "svgBtn",
  props: {
    iconSize: {
      type: Number,
      default: 24
    },

    btnSize: {
      type: Number,
      default: 48
    },

    btnStyle: {
      type: String,
    }
  },

  setup(props) {


    const iconData = reactive(
        {
          btnStyle: props.btnStyle,
          iconSize: props.iconSize
        }
    );

    const iconFile = computed(() => {
      const file = path.join(resDir, iconData.btnStyle + "-" + iconData.iconSize + "px.svg")
      console.log(file)
      return file
    })

    const btnH = computed(() => {
      return props.btnSize + "px"

    })

    const btnW = computed(() => {
      return props.btnSize + "px"

    })

    return {
      btnH,
      btnW,
      iconFile,
      ...toRefs(iconData),
    }

  }
}

</script>


<style scoped>
* {
  margin: 0;
  padding: 0;
}

#svg-btn-container {
  display: flex;
  width: 48px;
  height: 48px;
  justify-content: center;

}

#svg-btn-container:hover {
  border-radius: 50%;
  background-color: rgba(0, 0, 0, .1);
}

img {

}

</style>
