<template>
  <div class="addon-settings-controls">
    <button id="addon-config-timer-adapter" class="text-button addon-settings-config " data-l10n-id="addon-configure">Configure</button>
    <button id="addon-update-timer-adapter" class="text-button addon-settings-update hidden" data-l10n-id="addon-update">Update</button>
    <span class="addon-settings-spacer"></span>
    <button id="addon-remove-timer-adapter" class="text-button addon-settings-remove" data-l10n-id="addon-remove">Remove</button>
    <button id="addon-toggle-timer-adapter" class="text-button addon-settings-disable">
      Disable
    </button>
  </div>
</template>

<script>
export default {
  name: "discovered-addon-item"
}
</script>

<style scoped>
div{
  display: block;
}

.addon-settings-config, .addon-settings-controls{
  padding: 0 1rem 1rem 1rem;
  text-align: right;
}

</style>