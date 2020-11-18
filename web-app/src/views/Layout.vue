<template>
<div class="side">
<slot name="sidebar"></slot>
</div>
  <div class="right-space">
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

.right-space{
  display: flex;
  width: 80%;
  flex-direction: column;
  background-color: #ffffff;
}

</style>