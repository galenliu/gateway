<template>
  <router-view/>
</template>

<script lang="ts">
import {provide, ref} from "vue"

export default {
  name: 'App',
  setup() {
    const menuVisible = ref(false)
    provide("visible", ref())
  }
}
</script>

<style lang="scss">
@use ".@material/drawer";
@use "@material/list";

@include drawer.core-styles;
@include drawer.dismissible-core-styles;
@include drawer.modal-core-styles;
@include list.core-styles;
</style>
