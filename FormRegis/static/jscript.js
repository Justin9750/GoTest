//Logout confirmation
var coll = document.getElementsByClassName("colapsible")
for (var i = 0; i < coll.length; i++) {
  coll[i].addEventListener("click", function() {
    this.classList.toggle("on")
    var content = this.nextElementSibling
    if (content.style.display === "block") {
      content.style.display = "none"
    } else {
      content.style.display = "block"
    }
  })
}

//Back button
document.getElementById("home").addEventListener("click", function(){
  window.history.back();
})

/*
//Validation
function validate() {
  alert("Email atau password tidak terdaftar!\n " +
      "Periksa kembali penulisan kredensial anda");
}
*/