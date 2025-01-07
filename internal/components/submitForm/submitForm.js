async function submitPost() {
  const form = document.getElementById("fileUploadForm");
  const data = new FormData(form);
  let bodyText = document.getElementById("bodyText").value;
  if (bodyText.length < 5) {
    document.getElementById("errorField").innerHTML = "too short";
  } else if (bodyText.length > 1000) {
    document.getElementById("errorField").innerHTML = "too long";
  } else {
    console.log(data);
    let response = await fetch("/submitRoot", {
      method: "POST",
      body: data,
    });

    let res = await response.json();
    if (res.success == "true") {
      togglePostForm();
      window.location = window.location.origin + "/post/" + res.replyID;
    } else {
      document.getElementById("errorField").innerHTML = res.error;
    }
  }
}
