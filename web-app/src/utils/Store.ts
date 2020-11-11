import { reactive } from "vue"
import {getMenuList} from  "./Utils.ts"

const debug =true;

const store = reactive({
    state: reactive({
        sideUnfold: true,
    }),
    setSideUnfoldAction(newVale){
        if(debug){
            console.log("setSideUnfoldAction triggered with:  ",newVale)
        }
        this.state.sideUnfold=newVale
    }

})

const SideItemReactiveData =reactive({
    items: reactive({
       stata:getMenuList()
    }),
    setSelected(newValue){
        console.log("setSelected triggered with:--------------  ",newVale)

    }
})

export   { SideItemReactiveData, store }