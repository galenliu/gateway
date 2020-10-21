
let Menu;

init: function() {
  this.menuButton = document.getElementById('menu-button');
  this.menuButton.addEventListener('click', Menu.toggle.bind(Menu));
}