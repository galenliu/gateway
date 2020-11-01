<template>
  <aside class="nav">
    <div class="main-menu">
      <svg-icon icon-name="home"></svg-icon>
    </div>
  </aside>
</template>

<script lang="ts">
import svgIcon from "./svgIcon.vue"
import path from "path"

export default {
  name: "sidebar",

  components: {
    svgIcon,
  }
}
</script>

<style scoped>

</style>