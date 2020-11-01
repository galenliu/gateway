'use strict';

const Menu = {

    init: function () {

        this.hidden = true;

        //菜单按钮
        this.menuBtn = document.getElementById("menu-button")

        //找到menu-main和menu-scrim，绑定事件
        this.scrim = document.getElementById("menu-scrim")
        this.mainMenu = document.getElementById("main-menu")

        //菜单的每一项
        this.items = {}
        this.items.home = document.getElementById("home-menu-item")
        this.items.room = document.getElementById("rooms-menu-ite")
        this.items.rules = document.getElementById("rules-menu-item")
        this.items.log = document.getElementById("logs-menu-item")
        this.items.settings = document.getElementById("settings-menu-item")
        this.currtItem = "home"
    };


    toggle: function () {
        console.log("menu button is cliek")
        if this.hidden{
            this.show();
        } else {
            this.hiden();
        }
    },

    hiden: function () {
        this.scrim.classList.add("hidden")
        this.mainMenu.classList.add("hidden")
        this.menuBtn.classList.remove("menu-show")
        this.hidden = true
    }

    show: function () {
        this.scrim.classList.remove("hidden")
        this.mainMenu.classList.remove("hidden")
        this.menuBtn.classList.add("menu-show")
        this.hidden = false
    }


};

module.exports = Menu;