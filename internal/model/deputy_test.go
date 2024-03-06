package model

import (
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)

func TestFirmDeputy_GetRAGRating(t *testing.T) {
	tests := []struct {
		firmDeputy FirmDeputy
		want       RAGRating
	}{
		{
			firmDeputy: FirmDeputy{},
			want:       RAGRating{},
		},
		{
			firmDeputy: FirmDeputy{MarkedAsClass: "red"},
			want:       RAGRating{Name: "High risk", Colour: "red"},
		},
		{
			firmDeputy: FirmDeputy{MarkedAsClass: "amber"},
			want:       RAGRating{Name: "Medium risk", Colour: "orange"},
		},
		{
			firmDeputy: FirmDeputy{MarkedAsClass: "green"},
			want:       RAGRating{Name: "Low risk", Colour: "green"},
		},
	}
	for i, test := range tests {
		t.Run("Scenario "+strconv.Itoa(i), func(t *testing.T) {
			assert.Equal(t, test.want, test.firmDeputy.GetRAGRating())
		})
	}
}
