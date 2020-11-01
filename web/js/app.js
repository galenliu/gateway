'use strict';

// eslint-disable-next-line prefer-const
let Menu;

const App = {

  init: function() {
    this.menuButton = document.getElementById('menu-button');
    this.menuButton.addEventListener('click', Menu.toggle.bind(Menu));

    Menu.init()
  }

};
module.exports = App;

Menu = require('./menu');


window.addEventListener("load",function app_onLoad(){
  window.removeEventListener("load",app_onLoad);
  console.log("widows load")
  App.init()
})