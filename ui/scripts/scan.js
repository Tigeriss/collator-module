let barcode = "";
let firstBarcode = "";
let orderNumber = "";
let layers = 0;
let isStarted = false;
let currentRow = 0;
let currentColumn = 0;
let hasNoError = true; // if there is a fault number in this row - set to true

document.addEventListener("DOMContentLoaded", startScanFormClose);
document.addEventListener('keydown', addNumber);

window.onbeforeunload = function(){
    return "Все данные сканирования будут потеряны!";
}

function addNumber(event) {
    if (isStarted) {
        if (event.key !== 'Enter') {
            barcode = barcode.concat(event.key);
        } else {
            if (currentColumn < layers) {
                let id = (currentRow - 1).toString().concat("_", currentColumn.toString());
                document.getElementById(id).innerText = barcode;
                if (currentColumn !== 0 && firstBarcode !== barcode) {
                    hasNoError = false;
                    document.getElementById(id).style.color = "red";
                } else if (currentColumn === 0) {
                    firstBarcode = barcode;
                }
                barcode = "";
                if (currentColumn === (layers - 1)) {
                    setStatus();
                    currentColumn = 0;
                    createRow();
                } else {
                    currentColumn++;
                }
            }
        }
    }
}

function startScanFormOpen() {
    document.getElementById("new-scan").style.display = "block";
}
function startScanFormClose() {
    document.getElementById("new-scan").style.display = "none";
}

/*
 set all the things for start scanning.
 */
function startScan() {
    // set values
    orderNumber = document.getElementById("order_number").value;
    layers = document.getElementById("layers").value;

    document.getElementById("order").innerText = "Заказ: ".concat(orderNumber);
    document.getElementById("date").innerText = currentDate();
    document.getElementById("new-scan").style.display = "none";
    if (layers > 0) {
        isStarted = true;
        setTable();
        document.getElementById("new_scan").disabled = true;
    }
}

function finishScan() {
    isStarted = false;
    document.getElementById("order").innerText = "Заказ: ";
    if (currentRow !== 0) {
        document.getElementById("scan_table").deleteRow(currentRow);
    }
    createReport();

}

function setTable() {
    createHeader();
    createRow();
    document.getElementById("scan_table_footer").colSpan = layers + 2;
}

function currentDate() {
    let date = new Date();
    let day = String(date.getDate()).padStart(2, '0');
    let month = String(date.getMonth() + 1).padStart(2, '0');
    let year = date.getFullYear();
    return day.concat(".", month, ".", year.toString())
}

function createHeader() {
    let mainWidth = (100 / (layers + 2)) / 100 * 97;
    let headerRow = document.createElement("tr");
    let firstColumn = document.createElement("th");
    firstColumn.innerText = "№ п/п";
    firstColumn.style.width = mainWidth.toString().concat("%");
    headerRow.appendChild(firstColumn);

    for (let i = 0; i < layers; i++) {
        let thScanNum = document.createElement("th");
        thScanNum.style.width = mainWidth.toString().concat("%");
        thScanNum.innerText = (i + 1).toString();
        headerRow.appendChild(thScanNum);
    }
    let lastColumn = document.createElement("th");
    lastColumn.style.width = (100 - (layers + 1) * mainWidth).toString().concat("%");
    lastColumn.innerText = "Статус";
    headerRow.appendChild(lastColumn);
    document.getElementById("scan_table_head").appendChild(headerRow);
}

function createRow() {
    let currRow = document.createElement("tr");
    currRow.setAttribute("id", currentRow.toString());

    // row number cell. need not id cause inner text set on creation
    let firstColumn = document.createElement("td");
    firstColumn.innerText = (currentRow + 1).toString().concat(".");
    currRow.appendChild(firstColumn);

    // cells for scan numbers
    for (let i = 0; i < layers; i++) {
        let cell = document.createElement("td");
        cell.setAttribute("id", (currentRow.toString().concat("_", i.toString())));
        cell.style.fontSize = "75%";
        currRow.appendChild(cell);
    }

    // status cell
    let lastColumn = document.createElement("td");
    lastColumn.setAttribute("id", (currentRow.toString().concat("_status")));
    currRow.appendChild(lastColumn);
    document.getElementById("scan_table_body").appendChild(currRow);
    currentRow++;
    hasNoError = true;
}

function setStatus() {
    let statusHTML;
    if (hasNoError) {
        statusHTML = "<img src=\"/static/imgs/true.png\" alt='ОК'>";
    } else {
        statusHTML = "<img src=\"/static/imgs/false.png\" alt='Ошибка'>";
    }
    document.getElementById((currentRow - 1).toString().concat("_status")).innerHTML = statusHTML;
}

function clearScanTable() {
    orderNumber = "";
    layers = 0;
    hasNoError = true;
    let allRows = document.getElementById("scan_table").querySelectorAll("tr");
    for (let i = 0; i < allRows.length - 1; i++) {
        document.getElementById("scan_table").deleteRow(0);
    }
    currentRow = 0;
    document.getElementById("new_scan").disabled = false;
}

function createReport() {
    let allRows = [];
    let scanTable = document.getElementById('scan_table');
    for(let i = 0; i < scanTable.tBodies[0].children.length; i++){
        let tr = scanTable.tBodies[0].children[i];
        let row = []
        for(let n = 0; n < tr.cells.length; n++){
            let tmp;
            if(n < tr.cells.length - 1){
                tmp = tr.cells[n].innerText;
            }else {
                tmp = tr.cells[n].getElementsByTagName('img')[0].getAttribute('alt');
            }
            row.push(tmp);
        }
        console.log(row);
        allRows.push(row);
    }
    let report = {
        'user': document.getElementById('user').innerText,
        'date': currentDate(),
        'order_number': orderNumber,
        'scans_amount': layers,
        'scan_rows': allRows,
    }
    let xmlHttpRequest = new XMLHttpRequest();
    let url = "/scan/send_report";
    xmlHttpRequest.open("POST", url, true);
    xmlHttpRequest.setRequestHeader("Content-Type", "application/json");
    xmlHttpRequest.onreadystatechange = function () {
        if(xmlHttpRequest.readyState === 4 && xmlHttpRequest.status === 200){
            clearScanTable();
        }
    }
    let data = JSON.stringify(report);
    xmlHttpRequest.send(data);
}
