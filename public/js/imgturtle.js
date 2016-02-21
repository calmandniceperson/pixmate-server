updateAdditionalOptions();

function updateTTLDef() {
  if (document.getElementById("optionsRadiosDefault").checked == true) {
    document.getElementById("uploadTTLTime").disabled = true;
    document.getElementById("uploadTTLTime").value = "";
    document.getElementById("uploadTTLViews").disabled = true;
    document.getElementById("uploadTTLViews").value = "";
  }
}

function updateTTLTime() {
  if (document.getElementById("optionsRadiosTime").checked == true) {
    document.getElementById("uploadTTLTime").disabled = false;
    document.getElementById("uploadTTLTime").value = 1;
    document.getElementById("uploadTTLViews").disabled = true;
    document.getElementById("uploadTTLViews").value = "";
  } else {
    document.getElementById("uploadTTLTime").disabled = true;
    document.getElementById("uploadTTLTime").value = "";
  }
}

function updateTTLViews() {
    if(document.getElementById("optionsRadiosViews").checked == true) {
        document.getElementById("uploadTTLViews").disabled = false;
        document.getElementById("uploadTTLViews").value = 1;
        document.getElementById("uploadTTLTime").disabled = true;
        document.getElementById("uploadTTLTime").value = "";
    } else {
        document.getElementById("uploadTTLViews").disabled = true;
        document.getElementById("uploadTTLViews").value = "";
    }
}

function updateSubmit() {
  var file = document.getElementById("uploadFile").files[0];
  if (file) {
    document.getElementById("submitButton").disabled = false;
  }
}

function updateAdditionalOptions() {
  var radios = document.getElementsByClassName("radio");
  if (document.getElementById("addOptCheckbox").checked == false) {
    for(var i = 0; i < radios.length; ++i) {
      var ele = radios[i];
      ele.style.display = "none";
    }
  } else {
    for(var i = 0; i < radios.length; ++i) {
      var ele = radios[i];
      ele.style.display = "block";
    }
  }
}
