package util

import "github.com/ministryofjustice/opg-sirius-supervision-firm-deputy-hub/internal/sirius"

type pair struct {
	k string
	v string
}

var validationMappings = map[string]map[string]pair{
	// pii
	"piiRequested": {
		"isEmpty": pair{"pii-requested", "The PII requested date is required and can't be empty"},
	},
	"piiReceived": {
		"isEmpty": pair{"pii-received", "The PII received date is required and can't be empty"},
	},
	"piiExpiry": {
		"isEmpty": pair{"pii-expiry", "The PII expiry date is required and can't be empty"},
	},
	"piiAmount": {
		"isEmpty": pair{"pii-amount", "The PII amount is required and can't be empty"},
	},
	// firm
	"firmName": {
		"stringLengthTooLong": pair{"firmName", "The firm name must be 255 characters or fewer"},
		"isEmpty":             pair{"firmName", "The firm name is required and can't be empty"},
	},
	"firmId": {
		"notGreaterThanInclusive": pair{"existing-firm", "Enter a firm name or number"},
	},
}

func RenameErrors(siriusError sirius.ValidationErrors) sirius.ValidationErrors {
	mappedErrors := sirius.ValidationErrors{}
	for fieldName, value := range siriusError {
		for errorType, errorMessage := range value {
			err := make(map[string]string)
			if mapping, ok := validationMappings[fieldName][errorType]; ok {
				err[errorType] = mapping.v
				mappedErrors[mapping.k] = err
			} else {
				err[errorType] = errorMessage
				mappedErrors[fieldName] = err
			}
		}
	}
	return mappedErrors
}
