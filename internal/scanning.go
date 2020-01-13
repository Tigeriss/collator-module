package internal

import (
	"encoding/json"
	"github.com/recoilme/pudge"
	"log"
	"net/http"
)

type ScanRow struct {
	SerialNumber string `json:"serial_number"`
	ScannedNumbers []string `json:"scanned_numbers"`
	Status         bool     `json:"status"`
}
type Report struct {
	User        string    `json:"user"`
	Date        string    `json:"date"`
	OrderNumber string    `json:"order_number"`
	ScansAmount int       `json:"scans_amount"`
	ScanRows    []ScanRow `json:"scan_rows"`
}

func jsonToReportObject(request *http.Request) error {
	decoder := json.NewDecoder(request.Body)
	report := Report{}
	decoder.Decode(&report)
	err := addReportRecord(report)
	if err != nil {
		return err
	}
	return nil
}

func getReportFromDB(orderNumber string) (Report, error) {
	report := Report{}
	defer closeAllDB()
	err := pudge.Get("./db/reports", orderNumber, &report)
	if err != nil {
		return report, err
	}
	return report, nil
}

func addReportRecord(report Report) error {
	defer closeAllDB()
	err := pudge.Set("./db/reports", report.OrderNumber, report)
	if err != nil {
		return  err
	}
	return nil
}

func DeleteReport(orderNumber string)  error {
	log.Println("going to delete " + orderNumber)
	defer pudge.CloseAll()
	err := pudge.Delete("./db/reports", orderNumber)
	if err != nil {
		return err
	}
	return nil
}