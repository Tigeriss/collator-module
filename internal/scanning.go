package internal

type ScanRow struct {
	ScannedNumbers []string `json:"scannedNumbers"`
	Status         bool     `json:"status"`
}
type Report struct {
	User        string    `json:"user"`
	Date        string    `json:"date"`
	OrderNumber string    `json:"orderNumber"`
	ScansAmount int       `json:"scansAmount"`
	ScanRows    []ScanRow `json:"scanRows"`
}

// var allScans[] ScanRow

// func CreateScanReport() {
//	defer pudge.CloseAll()
//	// get data from front: orderNum, scans amount
//	rep := Report{
//		User: "user1",
//		Date:        time.Now().Format("DD.MM.YYYY"),
//		OrderNumber: "order_num",
//		ScansAmount: 5,
//		ScanRows:    allScans,
//	}
//
//	err := pudge.Set("./db/reports", rep.OrderNumber, rep)
//	check(err)
// }

// func AddScannedRow(data ScanRow){
//	allScans = append(allScans, data)
// }
