package api

import (
	"fmt"
	"net/http"

	"compton-scat/internal/constants"
	"compton-scat/internal/kinematics"
	"compton-scat/internal/runbook"
	"compton-scat/internal/section"
)

type kinematicsRequest struct {
	Name      string  `json:"name"`
	EnergyKEV float64 `json:"energy_keV"`
	AngleDeg  float64 `json:"angle_deg"`
}

type kinematicsResponse struct {
	Name          string  `json:"name"`
	EnergyKEV     float64 `json:"energy_keV"`
	AngleDeg      float64 `json:"angle_deg"`
	LambdaNm      float64 `json:"lambda_nm"`
	LambdaPrimeNm float64 `json:"lambda_prime_nm"`
	DeltaLambdaNm float64 `json:"delta_lambda_nm"`
	ScatteredKEV  float64 `json:"scattered_keV"`
	RecoilKEV     float64 `json:"recoil_keV"`
	EnergyDeficit float64 `json:"energy_deficit"`
}

type sectionRequest struct {
	EnergyKEV float64 `json:"energy_keV"`
}

type sectionResponse struct {
	EnergyKEV         float64 `json:"energy_keV"`
	DifferentialAt90  float64 `json:"differential_90_m2_per_sr"`
	DifferentialAt180 float64 `json:"differential_180_m2_per_sr"`
	Total             float64 `json:"total_m2"`
	Thomson           float64 `json:"thomson_m2"`
	TotalRatio        float64 `json:"total_ratio"`
}

type trendResponse struct {
	EnergiesKEV []float64 `json:"energies_keV"`
	Totals      []float64 `json:"totals_m2"`
	Ratios      []float64 `json:"ratios"`
	Decreasing  bool      `json:"decreasing"`
}

func (s *Server) handleKinematics(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "POST only")
		return
	}
	var req kinematicsRequest
	if err := readJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	res, err := kinematics.Solve(kinematics.Input{
		EnergyJ:  constants.KEVToJoules(req.EnergyKEV),
		ThetaDeg: req.AngleDeg,
	})
	if err != nil {
		writeError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}
	res = kinematics.HoldKineLive(res)
	sum := res.Summary()
	out := kinematicsResponse{
		Name:          req.Name,
		EnergyKEV:     sum.EnergyKEV,
		AngleDeg:      sum.ThetaDeg,
		LambdaNm:      sum.LambdaNm,
		LambdaPrimeNm: sum.LambdaPrimeNm,
		DeltaLambdaNm: sum.DeltaLambdaNm,
		ScatteredKEV:  sum.ScatteredKEV,
		RecoilKEV:     sum.RecoilKEV,
		EnergyDeficit: constants.KEVToJoules(constants.JoulesToKEV(res.EnergyDeficit())),
	}
	entry := runbook.Entry{
		ID:            fmt.Sprintf("run-%d", s.book.NextSeq()+1),
		EnergyKEV:     sum.EnergyKEV,
		AngleDeg:      sum.ThetaDeg,
		LambdaNm:      sum.LambdaNm,
		LambdaPrimeNm: sum.LambdaPrimeNm,
		DeltaLambdaNm: sum.DeltaLambdaNm,
		ScatteredKEV:  sum.ScatteredKEV,
		RecoilKEV:     sum.RecoilKEV,
		Note:          req.Name,
	}
	if err := s.book.Add(entry); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, out)
}

func (s *Server) handleSection(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "POST only")
		return
	}
	var req sectionRequest
	if err := readJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	res, err := section.CrossSection(constants.KEVToJoules(req.EnergyKEV))
	if err != nil {
		writeError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, sectionResponse{
		EnergyKEV:         req.EnergyKEV,
		DifferentialAt90:  res.DifferentialAt90,
		DifferentialAt180: res.DifferentialAt180,
		Total:             res.Total,
		Thomson:           res.Thomson,
		TotalRatio:        res.TotalRatio,
	})
}

func (s *Server) handleTrend(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "GET only")
		return
	}
	tr := section.ComputeTrend(section.DefaultIntegrationNodes)
	writeJSON(w, http.StatusOK, trendResponse{
		EnergiesKEV: tr.EnergiesKEV,
		Totals:      tr.Totals,
		Ratios:      tr.Ratios,
		Decreasing:  tr.Decreasing,
	})
}
