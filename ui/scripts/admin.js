document.addEventListener("DOMContentLoaded", newUserFormClose);

function newUserFormOpen() {
    document.getElementById("new-user").style.display = "block";
}

function newUserFormClose() {
    document.getElementById("new-user").style.display = "none";
}