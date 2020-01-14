document.addEventListener("DOMContentLoaded", newUserFormClose);

function newUserFormOpen() {
    document.getElementById("new-user").style.display = "block";
}

function newUserFormClose() {
    document.getElementById("new-user").style.display = "none";
}

function deleteUser(login) {
    console.log(login);
    if(login !== 'admin') {
        let xmlHttpRequest = new XMLHttpRequest();
        let url = "/admin/delete_user";
        xmlHttpRequest.open("POST", url, true);
        xmlHttpRequest.setRequestHeader("Content-Type", "application/json");
        xmlHttpRequest.onreadystatechange = function () {
            if (xmlHttpRequest.readyState === 4 && xmlHttpRequest.status === 200) {
            }
        }
        let data = JSON.stringify(login);
        xmlHttpRequest.send(data);
    } else {
        alert("Фиг ты удалишь центрального админа!");
    }
}

function deleteReport(orderNumber) {
        let xmlHttpRequest = new XMLHttpRequest();
        let url = "/admin/delete_report";
        xmlHttpRequest.open("POST", url, true);
        xmlHttpRequest.setRequestHeader("Content-Type", "application/json");
        xmlHttpRequest.onreadystatechange = function () {
            if (xmlHttpRequest.readyState === 4 && xmlHttpRequest.status === 200) {
            }
        }
        let data = JSON.stringify(orderNumber);
        xmlHttpRequest.send(data);
}

function openReport(orderNumber) {
    let url = "/admin/open_report/?order_number=" + orderNumber;
    fetch(url).then((response) => {
        return response.json();
    }).then((dataJson) => {
        fillReportTable(dataJson);
    });
}

function fillReportTable(dataJson){
    console.log(dataJson);
}