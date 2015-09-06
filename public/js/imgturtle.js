document.getElementById("me-link").addEventListener("click", function(){
  document.getElementById("feed-link").classList.remove("selected");
  document.getElementById("me-link").classList.add("selected");
  document.querySelector(".flipdiv").classList.add("flipped");
});

document.getElementById("feed-link").addEventListener("click", function(){
  document.getElementById("me-link").classList.remove("selected");
  document.getElementById("feed-link").classList.add("selected");
  document.querySelector(".flipdiv").classList.remove("flipped");
});
