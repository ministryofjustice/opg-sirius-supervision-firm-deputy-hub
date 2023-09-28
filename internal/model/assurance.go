package model

type Assurance struct {
	ReportReviewDate    string  `json:"reportReviewDate"`
	VisitReportMarkedAs RefData `json:"reportMarkedAs"`
	AssuranceType       RefData `json:"assuranceType"`
}
