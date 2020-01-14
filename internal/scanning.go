package internal

import (
	"encoding/json"
	"github.com/recoilme/pudge"
	"net/http"
)

type Field struct {
	Value string `json:"value"`
	Valid bool `json:"valid"`
}

type Report struct {
	User        string    `json:"user"`
	Date        string    `json:"date"`
	OrderNumber string    `json:"order_number"`
	ScansAmount string       `json:"scans_amount"`
	ScanRows    [][]Field `json:"scan_rows"`
}

func jsonToReportObject(request *http.Request) error {
	report := Report{}
	decoder := json.NewDecoder(request.Body)
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
	defer pudge.CloseAll()
	err := pudge.Delete("./db/reports", orderNumber)
	if err != nil {
		return err
	}
	return nil
}