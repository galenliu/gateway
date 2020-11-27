<template>
  <ul>
    <li class="settings-item"><a id="domain-settings-link" href="/settings/domain"
                                 data-l10n-id="settings-domain">域名</a></li>
    <li class="settings-item"><a id="network-settings-link" href="/settings/network" data-l10n-id="settings-network">网络</a>
    </li>
    <li class="settings-item"><a id="user-settings-link" href="/settings/users" data-l10n-id="settings-users">用户</a>
    </li>
    <li class="settings-item"><a id="addon-settings-link" href="/settings/addons" data-l10n-id="settings-add-ons">附加组件</a>
    </li>
    <li class="settings-item"><a id="localization-settings-link" href="/settings/localization"
                                 data-l10n-id="settings-localization">本地化</a></li>
    <li class="settings-item"><a id="update-settings-link" href="/settings/updates" data-l10n-id="settings-updates">更新</a>
    </li>
    <li class="settings-item"><a id="authorization-settings-link" href="/settings/authorizations"
                                 data-l10n-id="settings-authorizations">授权 </a></li>
    <li class="settings-item"><a id="experiment-settings-link" href="/settings/experiments"
                                 data-l10n-id="settings-experiments">实验</a></li>
    <li class="settings-item"><a id="developer-settings-link" href="/settings/developer"
                                 data-l10n-id="settings-developer">开发者</a></li>
  </ul>
</template>

<script>
export default {
name: "Profile"
}
</script>
<style scoped>
.settings-item{
  margin: 0;
  padding: 0;
  list-style-type: none;

  background-position: right 2.8rem center;
  background-size: auto 2rem;
}
li{
  display: list-item;
  text-align: center;
}

.settings-item a {
  display: block;
  margin: 1rem 0;
  padding: 2rem 2rem 2rem 5.5rem;
  text-decoration: none;
  font-size: 1.8rem;
  background-position: 1rem;
  background-size: 3.2rem;
  background-repeat: no-repeat;
  outline: 0;
}
</style>