<template>
  <div id="sidebar-item-content">
  <a href="#" class="sidebar-item">
    <img :src="svgFile" id="sidebar-scroll-item-svg" alt=""/>
    <span id="sidebar-scroll-item-title">{{ title }} </span>
  </a>
  </div>
</template>

<script>
import {computed, reactive, toRefs} from 'vue'

export default {

  name: "sidebar-item",

  props: {
    item_style: {
      type: String,
      required: true

    },
    title: {
      type: String,
      required: true
    }
  },


  setup(props) {

    console.log(props)
    const sidebarScrollItemReactiveData = reactive({
      title: props.title,
    })

    const svgFile = computed(
        () => {
          return getSvgFile(props.item_style)
        })

    return {
          svgFile,
      ...toRefs(sidebarScrollItemReactiveData)
    }
  }
}

const styleList = {
  home: {
    order:1,
    iconSrc: "/src/assets/svg/home-24px.svg",
  },
  rule: {
    order:  2,
    iconSrc:  "/src/assets/svg/rule-24px.svg",
  },
  settings: {
    order:  3,
    iconSrc:  "/src/assets/svg/settings-24px.svg",
  },
  exit: {
    order:  4,
    iconSrc: "/src/assets/svg/exit_to_app-24px.svg",
  },
}

function getSvgFile(name) {
  return styleList[name]["iconSrc"]

}

</script>

<style scoped>

#sidebar-item-content{

}

a {
  margin: 4px;
  display: block;
  border: 1px solid red;
}

#sidebar-item-content{
  height: 48px;
  background-color: #ffffff;
}

#sidebar-scroll-item-svg{

}

</style>