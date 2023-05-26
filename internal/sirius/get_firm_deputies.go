package sirius

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strings"
)

type orderStatus struct {
	Handle string `json:"handle"`
	Label  string `json:"label"`
}

type client struct {
	Id int `json:"id"`
}

type order struct {
	Id          int         `json:"id"`
	Client      client      `json:"client"`
	OrderStatus orderStatus `json:"orderStatus"`
}

type orders struct {
	Order order `json:"order"`
}

type assuranceVisit struct {
	ReportReviewDate     string  `json:"reportReviewDate"`
	VisitReportMarkedAs  RefData `json:"assuranceVisitReportMarkedAs"` 
	AssuranceType        RefData `json:"assuranceType"`  
}

type executiveCaseManager struct {
	EcmId   int    `json:"id"`
	EcmName string `json:"displayName"`
}

type Deputies struct {
	DeputyId             int                  `json:"id"`
	Firstname            string               `json:"firstname"`
	Surname              string               `json:"surname"`
	DeputyNumber         int                  `json:"deputyNumber"`
	Orders               []orders             `json:"orders"`
	ExecutiveCaseManager executiveCaseManager `json:"executiveCaseManager"`
	OrganisationName     string               `json:"organisationName"`
	Town                 string               `json:"town"`
	AssuranceVisit       assuranceVisit       `json:"mostRecentlyCompletedAssuranceVisit"`
}

type FirmDeputy struct {
	DeputyId             int
	Firstname            string
	Surname              string
	DeputyNumber         int
	ActiveClientsCount   int
	ExecutiveCaseManager string
	OrganisationName     string
	Town                 string  
	ReviewDate           string
	MarkedAsLabel        string
	MarkedAsClass        string
	AssuranceType        string
}



type RefData struct {
	Handle string `json:"handle"`
	Label  string `json:"label"`
}

func (c *Client) GetFirmDeputies(ctx Context, firmId int) ([]FirmDeputy, error) {
	req, err := c.newRequest(ctx, http.MethodGet, fmt.Sprintf("/api/v1/firms/%d/deputies", firmId), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	
	if resp.StatusCode == http.StatusUnauthorized {
		return nil, ErrUnauthorized
	}

	if resp.StatusCode != http.StatusOK {
		return nil, newStatusError(resp)
	}

	var v []Deputies
	if err = json.NewDecoder(resp.Body).Decode(&v); err != nil {
		return nil, err
	}

	var deputies []FirmDeputy

	for _, t := range v {
		var deputy = FirmDeputy{
			DeputyId:             t.DeputyId,
			Firstname:            t.Firstname,
			Surname:              t.Surname,
			DeputyNumber:         t.DeputyNumber,
			ActiveClientsCount:   getActiveClientCount(t.Orders),
			ExecutiveCaseManager: t.ExecutiveCaseManager.EcmName,
			OrganisationName:     t.OrganisationName,
			Town:                 t.Town,
			AssuranceType:        t.AssuranceVisit.AssuranceType.Label,
			MarkedAsLabel:        t.AssuranceVisit.VisitReportMarkedAs.Label,
			MarkedAsClass:        strings.ToLower(t.AssuranceVisit.VisitReportMarkedAs.Label),
			ReviewDate:           FormatDateAndTime(DateTimeFormat, t.AssuranceVisit.ReportReviewDate, DateTimeDisplayFormat),
		}

		deputies = append(deputies, deputy)
	}
	sortedDeputies := sortTheDeputiesByNumberOfClients(deputies)
	return sortedDeputies, err
}

func getActiveClientCount(orders []orders) int {
	clientId := getListOfClientIds(orders)
	uniqueClientIds := removeDuplicateIDs(clientId)

	return len(uniqueClientIds)
}

func getListOfClientIds(orders []orders) []int {
	listClientIds := []int{}
	for _, k := range orders {
		if k.Order.OrderStatus.Handle == "ACTIVE" {
			listClientIds = append(listClientIds, k.Order.Client.Id)
		}
	}
	return listClientIds
}

func removeDuplicateIDs(clientIds []int) []int {
	allKeys := make(map[int]bool)
	uniqueClientIds := []int{}
	for _, k := range clientIds {
		if _, value := allKeys[k]; !value {
			allKeys[k] = true
			uniqueClientIds = append(uniqueClientIds, k)
		}
	}
	return uniqueClientIds
}

func sortTheDeputiesByNumberOfClients(firmDeputies []FirmDeputy) []FirmDeputy {
	sort.Slice(firmDeputies, func(i, j int) bool {
		return firmDeputies[i].ActiveClientsCount > firmDeputies[j].ActiveClientsCount
	})
	return firmDeputies
}
