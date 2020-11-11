<template>
<div class="side">
<slot name="sidebar"></slot>
</div>
  <div class="work-space">
    <slot name="top-bar" ></slot>
    <slot name="work-space"></slot>
  </div>
</template>

<script>
export default {
name: "Layout"
}
</script>

<style scoped>
.side{

}

.work-space{
  display: flex;
  flex-grow: 1;
  flex-direction: column;
  background-color: #ffffff;
}

</style>