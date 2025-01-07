function togglePostForm() {
    let pf = document.getElementById("section-submitForm").style.display;
    if (pf != "block") {
        document.getElementById("section-submitForm").style.display = "block";
        document.getElementById("newPostButt").innerHTML = "-";
        document.getElementById("newPostButt").style.background = "#8d561f";
        document.getElementById("newPostButt").style.border = "1px solid #6b3000";
    } else {
        document.getElementById("section-submitForm").style.display = "none";
        document.getElementById("newPostButt").innerHTML = "+";
        document.getElementById("newPostButt").style.background = "#9cd4d5";
        document.getElementById("newPostButt").style.border = "1px solid #538d99";
    }

}
function getStream(sortOrder) {
    window.location = window.location.origin + "/" + sortOrder;
}
