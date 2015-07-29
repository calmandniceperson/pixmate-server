/*function upload() {
  var fromData = $("#fileBox").serialize(); // your form's data
  $.ajax({
      type: "POST",
      url: "/upload",
      data: fromData //sends the data to the new page.
    })
    .done(function(msg) {
      window.location.href = '/me' // redirects the page when finished.
    });
}
*/
var app = angular.module("imgturtleWeb", []);

app.controller("FileUploadCtrl", function($scope, $http) {
  $scope.upload = function() {
    //var fromData = $("#fileBox").serialize(); // your form's data
    var imgInp = document.getElementById("fileBox");
    var img = imgInp.files[0];
    var fullPath = imgInp.value;
    if (fullPath) {
      var startIndex = (fullPath.indexOf('\\') >= 0 ? fullPath.lastIndexOf('\\') : fullPath.lastIndexOf('/'));
      var filename = fullPath.substring(startIndex);
      if (filename.indexOf('\\') === 0 || filename.indexOf('/') === 0) {
        filename = filename.substring(1);
      }
    }

    var request = $http({
      method: "post",
      url: "/upload",
      data: {
        "img": img,
        name: filename,
        size: img.size.toString()
      }
    });

    request.success(function(data) {
      window.location = "/me";
    });

    request.error(function(data) {
      alert(data);
    });
  }
});
