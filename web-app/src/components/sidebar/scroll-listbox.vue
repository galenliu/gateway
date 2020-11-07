<template>
<div class="sidebar-scroll">
    <list-item v-for="(item, i) in items" :itemid="i" :rel="item.name" :title=item.title :icon="item.icon" @OnClick:ItemSelected="SelectedItem">
    </list-item>
</div>
</template>

<script lang="ts">
import {
    reactive,
    toRefs
} from 'vue'
import ListItem from '../common/list-item.vue'

import {
    CheckCorrect,
    Home,
    SettingConfig
} from "@icon-park/vue-next"

const items = [{
        name: "home",
        title: "家庭",
        icon: Home,
        selected: true
    },
    {
        name: "rule",
        title: "规则",
        icon: CheckCorrect,
        selected: false
    },
    {
        name: "config",
        title: "配置",
        icon: SettingConfig,
        selected: false
    }
]

export default {
    name: "ScrollListbox",

    components: {
        Home,
        SettingConfig,
        CheckCorrect,
        ListItem,
    },
    props: {
        items: {
            type: Object,
            default: items,
        }
    },

    setup(props, context) {
        const scrollReactiveData = reactive({
            items: props.items,

        })

      const SelectedItem =(e)=> {
        console.log(e.currentTarget )
      }

        return {
            SelectedItem,
            Home,
            ...toRefs(scrollReactiveData)
        }
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
