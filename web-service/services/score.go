package services

import (
	"capstone/middleware"
	"database/sql"
	"net/http"
	"strconv"
)

type ScoreService struct {
	DB *sql.DB
}

type Score struct {
	Id       int `json:"id,omitempty"`
	ReportId int `json:"report_id"`
	Score    int `json:"score"`
}

func (s *ScoreService) InsertNewScore(score Score) error {
	_, err := s.DB.Exec(`INSERT INTO SCORE(report_id, score) VALUES ($1, $2)`, score.ReportId, score.Score)
	if err != nil {
		return err
	}
	return nil
}

func (s *ScoreService) GetScoreFromReportId(id int) (int, error) {
	row := s.DB.QueryRow(`SELECT score FROM score WHERE report_id = $1`, id)
	var score int
	err := row.Scan(&score)
	if err != nil {
		return 0, err
	}
	return score, nil
}

func (s *ScoreService) IncrementScoreByReportId(id int) error {
	score, err := s.GetScoreFromReportId(id)
	if err != nil {
		return err
	}
	_, err = s.DB.Exec(`UPDATE SCORE SET score = $1 WHERE report_id = $2`, score+1, id)
	if err != nil {
		return err
	}
	return nil
}

func (s *ScoreService) GetScoreFromIdHandler(w http.ResponseWriter, r *http.Request) {
	middleware.LogRun("GET SCORE - START")
	var res middleware.BaseResponse
	idStr := r.FormValue("id")
	if len(idStr) <= 0 {
		res.FailResponse("Need param Id")
		res.WriteResponse(w)
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		res.FailResponse("Can not parse id")
		res.WriteResponse(w)
		return
	}
	score, err := s.GetScoreFromReportId(id)
	if err != nil {
		middleware.LogError("GET SCORE - PARSE ID", err)
		res.FailResponse("Can not get score")
		res.WriteResponse(w)
		return
	}
	middleware.LogRun("GET SCORE - END")
	res.InsertData(map[string]int{"score": score}, "success")
	res.WriteResponse(w)
	return
}

func (s *ScoreService) InsertScoreHandler(w http.ResponseWriter, r *http.Request) {
	middleware.LogRun("INSERT SCORE - START")
	var res middleware.BaseResponse
	idStr := r.FormValue("id")
	if len(idStr) <= 0 {
		res.FailResponse("Need param Id")
		res.WriteResponse(w)
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		middleware.LogError("INSERT SCORE - PARSE ID", err)
		res.FailResponse("Can not parse id")
		res.WriteResponse(w)
		return
	}
	score := Score{
		ReportId: id,
		Score:    1,
	}
	err = s.InsertNewScore(score)
	if err != nil {
		middleware.LogError("INSERT SCORE - EXEC", err)
		res.FailResponse("Can not insert score")
		res.WriteResponse(w)
		return
	}
	middleware.LogRun("INSERT SCORE - END")
	res.InsertData(nil, "Success")
	res.WriteResponse(w)
	return
}

func (s *ScoreService) UpdateScoreHandler(w http.ResponseWriter, r *http.Request) {
	middleware.LogRun("UPDATE SCORE - START")
	var res middleware.BaseResponse
	idStr := r.FormValue("id")
	if len(idStr) <= 0 {
		res.FailResponse("Need param Id")
		res.WriteResponse(w)
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		middleware.LogError("UPDATE SCORE - PARSE ID", err)
		res.FailResponse("Can not parse id")
		res.WriteResponse(w)
		return
	}
	err = s.IncrementScoreByReportId(id)
	if err != nil {
		middleware.LogError("UPDATE SCORE - QUERY", err)
		res.FailResponse("Can not update score")
		res.WriteResponse(w)
		return
	}

	middleware.LogRun("UPDATE SCORE - END")
	res.InsertData(nil, "Success")
	res.WriteResponse(w)
	return
}
