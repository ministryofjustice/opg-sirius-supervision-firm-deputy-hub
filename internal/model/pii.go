package model

type PiiDetails struct {
	FirmId       int     `json:"firmId"`
	PiiReceived  string  `json:"piiReceived"`
	PiiExpiry    string  `json:"piiExpiry"`
	PiiAmount    float64 `json:"piiAmount,omitempty"`
	PiiRequested string  `json:"piiRequested"`
}
