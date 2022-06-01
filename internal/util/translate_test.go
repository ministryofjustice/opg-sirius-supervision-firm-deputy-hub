package util

import (
	"github.com/ministryofjustice/opg-sirius-supervision-firm-deputy-hub/internal/sirius"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRenameErrors(t *testing.T) {
	siriusErrors := sirius.ValidationErrors{
		"piiRequested": map[string]string{"isEmpty": "isEmpty"},
		"firmName":     map[string]string{"stringLengthTooLong": "isEmpty"},
	}
	expected := sirius.ValidationErrors{
		"pii-requested": map[string]string{"isEmpty": "The PII requested date is required and can't be empty"},
		"firmName":      map[string]string{"stringLengthTooLong": "The firm name must be 255 characters or fewer"},
	}

	assert.Equal(t, expected, RenameErrors(siriusErrors))
}

func TestRenameErrors_default(t *testing.T) {
	siriusErrors := sirius.ValidationErrors{
		"x": map[string]string{"y": "z"},
	}

	assert.Equal(t, siriusErrors, RenameErrors(siriusErrors))
}
