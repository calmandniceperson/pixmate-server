document.getElementById("feed-link").focus();
document.getElementById("random-link").addEventListener("click", function(){
  document.getElementById("feed-link").classList.remove("selected");
  document.getElementById("random-link").classList.add("selected");
  document.querySelector(".flipdiv").classList.add("flipped");
});

document.getElementById("feed-link").addEventListener("click", function(){
  document.getElementById("random-link").classList.remove("selected");
  document.getElementById("feed-link").classList.add("selected");
  document.querySelector(".flipdiv").classList.remove("flipped");
});
