var app = angular.module("imgturtleWeb", []);

app.controller("SignUpCtrl", function($scope, $http) {
  $scope.submit = function() {
    var uname = document.getElementById("inp_uname").value;
    var pwd = document.getElementById("inp_pwd").value;
    var rpwd = document.getElementById("inp_rpwd").value;
    var email = document.getElementById("inp_email").value;
    var errorSpan = document.getElementById("errortext");

    if (uname == "" && pwd == "" && rpwd == "" && email == "") {
      errorSpan.innerHTML = "You didn't enter any values. Please fill all fields.";
      errorSpan.style.visibility = "visible";
    } else if (uname == "") {
      errorSpan.innerHTML = "You didn't enter a username.";
      errorSpan.style.visibility = "visible";
    } else if (uname.length < 4) {
      errorSpan.innerHTML = "The username you entered is not long enough.";
      errorSpan.style.visibility = "visible";
    } else if (pwd.length < 6) {
      errorSpan.innerHTML = "The password you entered is not long enough.";
      errorSpan.style.visibility = "visible";
    } else if (email == "") {
      errorSpan.innerHTML = "You didn't enter an e-mail address.";
      errorSpan.style.visibility = "visible";
    } else if (email.length < 4 || !$scope.validateEmail(email)) {
      errorSpan.innerHTML = "The e-mail address you entered is no valid address.";
      errorSpan.style.visibility = "visible";
    } else {
      errorSpan.style.visibility = "hidden";

      var request = $http({
        method: "post",
        url: "/signup",
        data: {
          "uname": uname,
          "pwd": pwd,
          "email": email
        }
      });

      request.success(function(data) {
        alert("Success.");
        window.location = "/";
      });

      request.error(function(data) {
        errorSpan.innerHTML = data;
        errorSpan.style.visibility = "visible";
      });
    }
  };

  $scope.validateEmail = function(email) {
    var re = /^([\w-]+(?:\.[\w-]+)*)@((?:[\w-]+\.)*\w[\w-]{0,66})\.([a-z]{2,6}(?:\.[a-z]{2})?)$/i;
    return re.test(email);
  }

  $scope.checkPasswordMatch = function() {
    if (("" + $scope.pwd).length > 0 && ("" + $scope.rpwd).length > 0) {
      if ($scope.pwd == $scope.rpwd) {
        document.getElementById("btn_submit").style.visibility = "visible";
        document.getElementById("inp_rpwd").style.color = "#000";
        document.getElementById("inp_pwd").style.color = "#000";
      } else {
        if (document.getElementById("btn_submit").style.visibility == "visible") {
          document.getElementById("btn_submit").style.visibility = "hidden";
        }
        document.getElementById("inp_pwd").style.color = "rgb(193, 96, 96)";
        document.getElementById("inp_rpwd").style.color = "rgb(193, 96, 96)";
      }
    }
  }
});
