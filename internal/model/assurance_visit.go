package model

type AssuranceVisit struct {
	ReportReviewDate    string  `json:"reportReviewDate"`
	VisitReportMarkedAs RefData `json:"assuranceVisitReportMarkedAs"`
	AssuranceType       RefData `json:"assuranceType"`
}
