package sirius

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ministryofjustice/opg-sirius-supervision-firm-deputy-hub/internal/mocks"
	"github.com/stretchr/testify/assert"
)

func TestGetFirmDeputiesReturned(t *testing.T) {
	mockClient := &mocks.MockClient{}
	client, _ := NewClient(mockClient, "http://localhost:3000")

	json := `[
		{"id":76,"orders":[{"order":{"id":63,"client":{"id":74,"firstname":"Louis","surname":"Dauphin"},"orderStatus":{"handle":"ACTIVE","label":"Active"}}}],"deputyNumber":21,"organisationName":"pro dept","executiveCaseManager":{"id":94,"displayName":"PROTeam1 User1"},"firm":{"id":1}},
		{"id":77, "firstName":"Louis", "surname":"Devito", "orders":[{"order":{"id":49,"client":{"id":99,"firstname":"Bob","surname":"Mortimer"},"orderStatus":{"handle":"ACTIVE","label":"Active"}}}],"deputyNumber":25,"organisationName":"","executiveCaseManager":{"id":94,"displayName":"PROTeam1 User1"},"firm":{"id":1}}
	]`

	r := ioutil.NopCloser(bytes.NewReader([]byte(json)))

	mocks.GetDoFunc = func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 200,
			Body:       r,
		}, nil
	}

	expectedResponse := []FirmDeputy{
		{
			Firstname:            "",
			Surname:              "",
			DeputyId:             76,
			DeputyNumber:         21,
			ActiveClientsCount:         1,
			ExecutiveCaseManager: "PROTeam1 User1",
			OrganisationName:     "pro dept",
		},
		{
			Firstname:            "Louis",
			Surname:              "Devito",
			DeputyId:             77,
			DeputyNumber:         25,
			ActiveClientsCount:         1,
			ExecutiveCaseManager: "PROTeam1 User1",
			OrganisationName:     "",
		},
	}

	firmDetails, err := client.GetFirmDeputies(getContext(nil), 1)

	assert.Equal(t, expectedResponse, firmDetails)
	assert.Nil(t, err)
}

func TestGetFirmDeputiesReturnsNewStatusError(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}))
	defer svr.Close()

	client, _ := NewClient(http.DefaultClient, svr.URL)

	firmDetails, err := client.GetFirmDeputies(getContext(nil), 1)

	var expectedResponse []FirmDeputy

	assert.Equal(t, expectedResponse, firmDetails)
	assert.Equal(t, StatusError{
		Code:   http.StatusMethodNotAllowed,
		URL:    svr.URL + "/api/v1/firms/1/deputies",
		Method: http.MethodGet,
	}, err)
}

func TestGetFirmDeputiesReturnsUnauthorisedClientError(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer svr.Close()

	client, _ := NewClient(http.DefaultClient, svr.URL)

	firmDetails, err := client.GetFirmDeputies(getContext(nil), 1)

	var expectedResponse []FirmDeputy

	assert.Equal(t, ErrUnauthorized, err)
	assert.Equal(t, expectedResponse, firmDetails)
}


func TestGetActiveClientCount(t *testing.T) {
	testOrders := []orders{
		order{
			Id: 5,
			Client: client{
				Id: 99,
			},
			OrderStatus: orderStatus{
				Handle: "ACTIVE",
				Label:  "Active",
			},
		},
	}
	assert.Equal(t, 3, getActiveClientCount(testOrders))
}

//type T struct {
//	Orders []struct {
//		Order struct {
//			Id     int `json:"id"`
//			Client struct {
//				Id        int    `json:"id"`
//				Firstname string `json:"firstname"`
//				Surname   string `json:"surname"`
//			} `json:"client"`
//			OrderStatus struct {
//				Handle string `json:"handle"`
//				Label  string `json:"label"`
//			} `json:"orderStatus"`
//		} `json:"order"`
//	} `json:"orders"`
//}