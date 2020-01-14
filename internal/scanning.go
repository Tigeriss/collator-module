package internal

import (
	"encoding/json"
	"github.com/recoilme/pudge"
	"log"
	"net/http"
	"path"
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
	err := pudge.Get(path.Join(".", "db", "reports"), orderNumber, &report)
	log.Println(report.ScanRows[0][0].Value)
	if err != nil {
		return report, err
	}
	return report, nil
}

func addReportRecord(report Report) error {
	defer closeAllDB()
	log.Println(report.ScanRows[0][0].Value)
	err := pudge.Set(path.Join(".", "db", "reports"), report.OrderNumber, report)
	if err != nil {
		return  err
	}
	return nil
}

func DeleteReport(orderNumber string)  error {
	defer pudge.CloseAll()
	err := pudge.Delete(path.Join(".", "db", "reports"), orderNumber)
	if err != nil {
		return err
	}
	return nil
}