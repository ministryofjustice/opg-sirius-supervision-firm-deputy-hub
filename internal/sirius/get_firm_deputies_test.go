package sirius

import (
	"bytes"
	"github.com/ministryofjustice/opg-sirius-supervision-firm-deputy-hub/internal/model"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ministryofjustice/opg-sirius-supervision-firm-deputy-hub/internal/mocks"
	"github.com/stretchr/testify/assert"
)

func TestGetFirmDeputiesReturned(t *testing.T) {
	mockClient := &mocks.MockClient{}
	client, _ := NewClient(mockClient, "http://localhost:3000")

	json := `[{
			"id":76,
			"orders":[
				{
					"order":{
						"id":63,
						"client":{
							"id":74,
							"firstname":"Louis",
							"surname":"Dauphin"
						},
						"orderStatus":{
							"handle":"ACTIVE",
							"label":"Active"
						}
					}
				}
			],
			"deputyNumber":21,
			"organisationName":"pro dept",
			"executiveCaseManager":{
				"id":94,
				"displayName":"PROTeam1 User1"
			},
			"firm":{
				"id":1
			},
			"mostRecentlyCompletedAssurance": {
				"reportReviewDate" : "2023-05-26T00:00:00+00:00",
				"reportMarkedAs": {
					"handle": "GREEN",
					"label": "Green"
				},
				"assuranceType": {
					"handle": "VISIT",
					"label": "Visit"
				}
			}
			
		},
		{
			"id":77, 
			"firstName":"Louis", 
			"surname":"Devito", 
			"orders":[
				{
					"order":{
						"id":49,
						"client":{
							"id":99,
							"firstname":"Bob",
							"surname":"Mortimer"
						},
						"orderStatus":{
							"handle":"ACTIVE",
							"label":"Active"
						}
					}
				}
			],
			"deputyNumber":25,
			"organisationName":"",
			"executiveCaseManager":{"id":94,"displayName":"PROTeam1 User1"},
			"firm":{
				"id":1
			}
		}
	]`

	r := io.NopCloser(bytes.NewReader([]byte(json)))

	mocks.GetDoFunc = func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 200,
			Body:       r,
		}, nil
	}

	expectedResponse := []model.FirmDeputy{
		{
			Firstname:            "",
			Surname:              "",
			DeputyId:             76,
			DeputyNumber:         21,
			ActiveClientsCount:   1,
			ExecutiveCaseManager: "PROTeam1 User1",
			OrganisationName:     "pro dept",
			ReviewDate:           "26/05/2023",
			MarkedAsLabel:        "Green",
			MarkedAsClass:        "green",
			AssuranceType:        "Visit",
		},
		{
			Firstname:            "Louis",
			Surname:              "Devito",
			DeputyId:             77,
			DeputyNumber:         25,
			ActiveClientsCount:   1,
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

	var expectedResponse []model.FirmDeputy

	assert.Equal(t, expectedResponse, firmDetails)
	assert.Equal(t, StatusError{
		Code:   http.StatusMethodNotAllowed,
		URL:    svr.URL + SupervisionAPIPath + "/v1/firms/1/deputies",
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

	var expectedResponse []model.FirmDeputy

	assert.Equal(t, ErrUnauthorized, err)
	assert.Equal(t, expectedResponse, firmDetails)
}

func TestGetActiveClientCountOnlyReturnsOneOrderWithSameClient(t *testing.T) {
	testOrders := []model.Orders{
		{
			Order: model.Order{
				Id: 5,
				Client: model.Client{
					Id: 99,
				},
				OrderStatus: model.OrderStatus{
					Handle: "ACTIVE",
					Label:  "Active",
				},
			},
		},
		{
			Order: model.Order{
				Id: 6,
				Client: model.Client{
					Id: 99,
				},
				OrderStatus: model.OrderStatus{
					Handle: "ACTIVE",
					Label:  "Active",
				},
			},
		},
	}
	assert.Equal(t, 1, getActiveClientCount(testOrders))
}

func TestGetActiveClientCountReturnsTwoOrdersWithTwoDifferentClients(t *testing.T) {
	testOrders := []model.Orders{
		{
			Order: model.Order{
				Id: 5,
				Client: model.Client{
					Id: 99,
				},
				OrderStatus: model.OrderStatus{
					Handle: "ACTIVE",
					Label:  "Active",
				},
			},
		},
		{
			Order: model.Order{
				Id: 6,
				Client: model.Client{
					Id: 44,
				},
				OrderStatus: model.OrderStatus{
					Handle: "ACTIVE",
					Label:  "Active",
				},
			},
		},
	}
	assert.Equal(t, 2, getActiveClientCount(testOrders))
}

func TestGetListOfClientIdsReturnsOnlyActiveClientsOnOrders(t *testing.T) {
	testOrders := []model.Orders{
		{
			Order: model.Order{
				Id: 5,
				Client: model.Client{
					Id: 99,
				},
				OrderStatus: model.OrderStatus{
					Handle: "ACTIVE",
					Label:  "Active",
				},
			},
		},
		{
			Order: model.Order{
				Id: 7,
				Client: model.Client{
					Id: 55,
				},
				OrderStatus: model.OrderStatus{
					Handle: "ACTIVE",
					Label:  "Active",
				},
			},
		},
		{
			Order: model.Order{
				Id: 6,
				Client: model.Client{
					Id: 44,
				},
				OrderStatus: model.OrderStatus{
					Handle: "OPEN",
					Label:  "Open",
				},
			},
		},
		{
			Order: model.Order{
				Id: 4,
				Client: model.Client{
					Id: 11,
				},
				OrderStatus: model.OrderStatus{
					Handle: "CLOSED",
					Label:  "Closed",
				},
			},
		},
		{
			Order: model.Order{
				Id: 3,
				Client: model.Client{
					Id: 77,
				},
				OrderStatus: model.OrderStatus{
					Handle: "DUPLICATE",
					Label:  "Duplicate",
				},
			},
		},
	}
	assert.Equal(t, []int{99, 55}, getListOfClientIds(testOrders))
}

func TestRemoveDuplicateIDs(t *testing.T) {
	clientIds := []int{99, 99, 55, 12, 46, 87, 99, 12}
	assert.Equal(t, []int{99, 55, 12, 46, 87}, removeDuplicateIDs(clientIds))
}

func TestSortTheDeputiesByNumberOfClients(t *testing.T) {
	firmDeputy :=
		[]model.FirmDeputy{
			{
				DeputyId:             12,
				Firstname:            "Missy",
				Surname:              "Longstocking",
				DeputyNumber:         25,
				ActiveClientsCount:   0,
				ExecutiveCaseManager: "Manager Name",
				OrganisationName:     "",
			},
			{
				DeputyId:             85,
				Firstname:            "",
				Surname:              "",
				DeputyNumber:         52,
				ActiveClientsCount:   44,
				ExecutiveCaseManager: "",
				OrganisationName:     "OrganisationName",
			},
			{
				DeputyId:             1,
				Firstname:            "One",
				Surname:              "Deadpool",
				DeputyNumber:         99,
				ActiveClientsCount:   1,
				ExecutiveCaseManager: "",
				OrganisationName:     "",
			},
			{
				DeputyId:             1,
				Firstname:            "Turtle",
				Surname:              "Pillow",
				DeputyNumber:         99,
				ActiveClientsCount:   4,
				ExecutiveCaseManager: "",
				OrganisationName:     "",
			},
		}

	expectedResult := []model.FirmDeputy{

		{
			DeputyId:             85,
			Firstname:            "",
			Surname:              "",
			DeputyNumber:         52,
			ActiveClientsCount:   44,
			ExecutiveCaseManager: "",
			OrganisationName:     "OrganisationName",
		},
		{
			DeputyId:             1,
			Firstname:            "Turtle",
			Surname:              "Pillow",
			DeputyNumber:         99,
			ActiveClientsCount:   4,
			ExecutiveCaseManager: "",
			OrganisationName:     "",
		},
		{
			DeputyId:             1,
			Firstname:            "One",
			Surname:              "Deadpool",
			DeputyNumber:         99,
			ActiveClientsCount:   1,
			ExecutiveCaseManager: "",
			OrganisationName:     "",
		},
		{
			DeputyId:             12,
			Firstname:            "Missy",
			Surname:              "Longstocking",
			DeputyNumber:         25,
			ActiveClientsCount:   0,
			ExecutiveCaseManager: "Manager Name",
			OrganisationName:     "",
		},
	}

	assert.Equal(t, expectedResult, sortTheDeputiesByNumberOfClients(firmDeputy))
}
