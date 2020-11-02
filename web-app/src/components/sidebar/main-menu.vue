<template>
  <div id="main-menu">
    <svgBtn btn-style="menu"></svgBtn>
    <span id="menu-btn-title"> {{ title }}</span>
  </div>
</template>

<script>
import svgBtn from "../svg-btn.vue"
import {reactive, toRefs} from "vue";

export default {
  name: "main-menu",

  components: {
    svgBtn: svgBtn
  },

  props: {
    title: {
      type: String,
      default: 'WebThings'
    }
  },

  setup(props) {
    const menuReactiveData = reactive({
      title: props.title
    })

    return {
      ...toRefs(menuReactiveData)
    }
  }
}
</script>

<style scoped>
#main-menu {
  display: inline-flex;
  align-items: center;
  background-color: #fafafa;
  border-bottom-color: red;
  overflow: hidden;
  padding-left: 8px;

  border-bottom: 1px solid rgba(0, 0, 0, 0.12);
  border-right: 1px solid rgba(0, 0, 0, 0.12);
  z-index: 10;
}

#menu-btn-title {
  margin-right: 20px;
  margin-left: 10px;
}

* {
  margin: 0;
  padding: 0;
}
</style>