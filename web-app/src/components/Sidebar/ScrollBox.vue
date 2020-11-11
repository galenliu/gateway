<template>
<div class="sidebar-scroll">
  <slot>
  </slot>
</div>
</template>

<script lang="ts">


export default {
    name: "ScrollBox",

    components: {

    },

    setup() {

    }
}

</script>

<style lang="scss" scoped>
.sidebar-scroll {
    display: flex;
    flex-direction: column;
    box-sizing: border-box;
}
</style>
