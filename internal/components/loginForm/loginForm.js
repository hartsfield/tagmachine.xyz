{{ define "loginForm.js" }}
function toggleLoginForm() {
    let pf = document.getElementById("section-submitForm").style.display;
    if (pf != "block") {
        document.getElementById("section-loginForm").style.display = "block";
        document.getElementById("showLoginButt").innerHTML = "-";
        document.getElementById("showLoginButt").style.background = "#8d561f";
        document.getElementById("showLoginButt").style.border = "1px solid #6b3000";
    } else {
        document.getElementById("section-loginForm").style.display = "none";
        document.getElementById("showLoginButt").innerHTML = "Login";
        document.getElementById("showLoginButt").style.background = "darkred";
        document.getElementById("showLoginButt").style.border = "1px solid red";
    }

}
// auth is used for signing up and signing in/out. path could be:
// /api/signup
// /api/signin
// /api/logout
function auth(path) {
    var xhr = new XMLHttpRequest();

    xhr.open("POST", "/" + path);
    xhr.setRequestHeader("Content-Type", "application/json");
    xhr.onload = function() {
        if (xhr.status === 200) {
            var res = JSON.parse(xhr.responseText);
            if (res.success == "false") {
                // If we aren't successful we display an error.
                document.getElementById("errorField").innerHTML = res.error;
            } else {
                // Reload the page now that the user is signed in.
                window.location.reload();
            }
        }
    };

    // For now, all we're sending is a username and password, but we may start
    // asking for email or mobile number at some point.
    xhr.send(JSON.stringify({
        password: document.getElementById("password").value,
        username: document.getElementById("username").value,
    }));
}


{{end}}
