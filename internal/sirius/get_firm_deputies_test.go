package sirius

import (
	"bytes"
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
			"mostRecentlyCompletedAssuranceVisit": {
				"reportReviewDate" : "2023-05-26T00:00:00+00:00",
				"assuranceVisitReportMarkedAs": {
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

	expectedResponse := []FirmDeputy{
		{
			Firstname:            "",
			Surname:              "",
			DeputyId:             76,
			DeputyNumber:         21,
			ActiveClientsCount:   1,
			ExecutiveCaseManager: "PROTeam1 User1",
			OrganisationName:     "pro dept",
			ReviewDate: "05/26/2023",
			MarkedAsLabel: "Green",
			MarkedAsClass: "green",
			AssuranceType: "Visit",
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

func TestGetActiveClientCountOnlyReturnsOneOrderWithSameClient(t *testing.T) {
	testOrders := []orders{
		{
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
		},
		{
			order{
				Id: 6,
				Client: client{
					Id: 99,
				},
				OrderStatus: orderStatus{
					Handle: "ACTIVE",
					Label:  "Active",
				},
			},
		},
	}
	assert.Equal(t, 1, getActiveClientCount(testOrders))
}

func TestGetActiveClientCountReturnsTwoOrdersWithTwoDifferentClients(t *testing.T) {
	testOrders := []orders{
		{
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
		},
		{
			order{
				Id: 6,
				Client: client{
					Id: 44,
				},
				OrderStatus: orderStatus{
					Handle: "ACTIVE",
					Label:  "Active",
				},
			},
		},
	}
	assert.Equal(t, 2, getActiveClientCount(testOrders))
}

func TestGetListOfClientIdsReturnsOnlyActiveClientsOnOrders(t *testing.T) {
	testOrders := []orders{
		{
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
		},
		{
			order{
				Id: 7,
				Client: client{
					Id: 55,
				},
				OrderStatus: orderStatus{
					Handle: "ACTIVE",
					Label:  "Active",
				},
			},
		},
		{
			order{
				Id: 6,
				Client: client{
					Id: 44,
				},
				OrderStatus: orderStatus{
					Handle: "OPEN",
					Label:  "Open",
				},
			},
		},
		{
			order{
				Id: 4,
				Client: client{
					Id: 11,
				},
				OrderStatus: orderStatus{
					Handle: "CLOSED",
					Label:  "Closed",
				},
			},
		},
		{
			order{
				Id: 3,
				Client: client{
					Id: 77,
				},
				OrderStatus: orderStatus{
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
		[]FirmDeputy{
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

	expectedResult := []FirmDeputy{

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
