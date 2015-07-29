var app = angular.module("imgturtleWeb", []);

app.controller("SignInCtrl", function($scope, $http) {
  $scope.submit = function() {
    var ue = document.getElementById("inp_ue").value;
    var pwd = document.getElementById("inp_pwd").value;
    var errorSpan = document.getElementById("errortext");

    if (ue == "" && pw == "") {
      errorSpan.innerHTML = "You didn't enter any values. Please fill all fields.";
      errorSpan.style.visibility = "visible";
    } else if (ue == "") {
      errorSpan.innerHTML = "You didn't enter a username or an e-mail address.";
      errorSpan.style.visibility = "visible";
    } else if (ue.length < 4) {
      if ($scope.validateEmail(ue)) {
        errorSpan.innerHTML = "The e-mail address you entered is not long enough.";
      } else {
        errorSpan.innerHTML = "The username you entered is not long enough.";
      }
      errorSpan.style.visibility = "visible";
    } else if (pwd == "") {
      errorSpan.innerHTML = "You didn't enter a password.";
      errorSpan.style.visibility = "visible";
    } else {
      errorSpan.style.visibility = "hidden";

      var request = $http({
        method: "post",
        url: "/signin",
        data: {
          "ue": ue,
          "pwd": pwd
        }
      });

      request.success(function(data) {
        alert("Success.");
        window.location = "/me";
      });

      request.error(function(data) {
        errorSpan.innerHTML = data;
        errorSpan.style.visibility = "visible";
      });
    }
  }

  $scope.validateEmail = function(email) {
    var re = /^([\w-]+(?:\.[\w-]+)*)@((?:[\w-]+\.)*\w[\w-]{0,66})\.([a-z]{2,6}(?:\.[a-z]{2})?)$/i;
    return re.test(email);
  }

  $scope.checkInput = function() {
    if (document.getElementById("inp_ue").value.length >= 4 &&
      document.getElementById("inp_pwd").value.length >= 6) {
      document.getElementById("btn_submit").style.visibility = "visible";
    } else {
      document.getElementById("btn_submit").style.visibility = "hidden";
    }
  }
});
