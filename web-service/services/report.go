package services

import (
	"capstone/helper"
	"capstone/middleware"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const (
	DAY          = 3
	DISTANCE     = 10
	REPORT_COUNT = 3
)

type ReportService struct {
	DB *sql.DB
}

type Report struct {
	Id         int     `json:"id,omitempty"`
	JssId      string  `json:"jss_id"`
	Jenis      string  `json:"jenis"`
	Keterangan string  `json:"keterangan"`
	Latitude   float64 `json:"latitude"`
	Longitude  float64 `json:"longitude"`
	TglLapor   string  `json:"tgl_lapor"`
	Accuracy   float64 `json:"accuracy,omitempty"`
}

func (r *Report) IsValidDistance(sourceLat float64, sourceLon float64, distance float64) bool {
	lat := r.Latitude
	lon := r.Longitude
	d := helper.GetDistanceFromLatLng(sourceLat, sourceLon, lat, lon)
	fmt.Printf("Distance: %v\n", d)
	if d > distance {
		return false
	}
	return true
}

func (s *ReportService) SelectReportFiltered(day int, category string, lat float64, lon float64, filterDistance float64) ([]Report, error) {
	yesterdayDate := time.Now().Add(time.Duration(-24*day) * time.Hour).Format("2006-01-02 15:04:05")
	fmt.Printf("Yesterday Date: %v\n", yesterdayDate)
	row, err := s.DB.Query(`SELECT * FROM REPORT WHERE tgl_lapor >= $1 AND jenis = $2`, yesterdayDate, category)
	if err != nil {
		return nil, err
	}
	var reportList []Report
	for row.Next() {
		var report Report
		err := row.Scan(&report.Id, &report.JssId, &report.Jenis, &report.Keterangan, &report.Latitude, &report.Longitude, &report.TglLapor)
		if err != nil {
			return nil, err
		}
		fmt.Printf("Report: %v\n", report)
		if report.IsValidDistance(lat, lon, filterDistance) == true {
			reportList = append(reportList, report)
		}
	}
	fmt.Printf("Report List: %v\n", reportList)
	return reportList, nil
}

func (s *ReportService) GetReportFiltered(sourceReport Report) ([]Report, error) {
	reportList, err := s.SelectReportFiltered(DAY, sourceReport.Jenis, sourceReport.Latitude, sourceReport.Longitude, DISTANCE)
	if err != nil {
		return nil, err
	}
	return reportList, nil
}

func (s *ReportService) InsertReport(report Report) (int, error) {
	var id int
	err := s.DB.QueryRow(`INSERT INTO REPORT(jss_id, jenis, keterangan, latitude, longitude, tgl_lapor) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`, report.JssId, report.Jenis, report.Keterangan, report.Latitude, report.Longitude, report.TglLapor).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (s *ReportService) InsertReportHandler(w http.ResponseWriter, r *http.Request) {
	middleware.LogRun("REQUEST REPORT - START")
	var res middleware.BaseResponse
	var req Report
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&req)
	if err != nil {
		middleware.LogError("REQUEST REPORT - DECODE", err)
		res.FailResponse("Gagal Mengambil Data")
		res.WriteResponse(w)
		return
	}

	reportList, err := s.GetReportFiltered(req)
	if err != nil {
		middleware.LogError("REPORT - QUERY FILTER", err)
		res.FailResponse("Gagal Mengambil Data")
		res.WriteResponse(w)
		return
	}

	pred := Predict{
		SourceText: req.Keterangan,
		ReportList: reportList,
	}
	pred.Predict()
	if err != nil {
		middleware.LogError("REQUEST REPORT - PREDICT", err)
		res.FailResponse("Fail to predict")
		res.WriteResponse(w)
		return
	}
	pred.SortByAccuracy()
	resTemp := map[string]interface{}{
		"original_report": req,
	}
	if len(pred.ReportList) >= REPORT_COUNT {
		resTemp["report_list"] = pred.ReportList[:REPORT_COUNT]
	} else {
		resTemp["report_list"] = pred.ReportList
	}
	res.InsertData(resTemp, "Sukses")
	res.WriteResponse(w)
	middleware.LogRun("REQUEST REPORT - END")
}

func (s *ReportService) InsertNewReportHandler(w http.ResponseWriter, r *http.Request) {
	middleware.LogRun("INSERT NEW REPORT - START")
	var res middleware.BaseResponse
	var req Report
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&req)
	if err != nil {
		middleware.LogError("INSERT NEW REPORT - DECODE", err)
		res.FailResponse("Fail to insert report")
		res.WriteResponse(w)
		return
	}

	id, err := s.InsertReport(req)
	if err != nil {
		middleware.LogError("INSERT NEW REPORT - EXEC DB", err)
		res.FailResponse("Fail to insert report")
		res.WriteResponse(w)
		return
	}

	scoreService := ScoreService{
		DB: s.DB,
	}

	score := Score{
		ReportId: int(id),
		Score:    1,
	}
	err = scoreService.InsertNewScore(score)
	if err != nil {
		middleware.LogError("INSERT NEW REPORT - INSERT SCORE", err)
		res.FailResponse("Fail to insert new score")
		res.WriteResponse(w)
		return
	}

	res.InsertData(nil, "succes")
	res.WriteResponse(w)
	middleware.LogRun("INSERT NEW REPORT - END")
}
