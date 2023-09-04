package sirius

import (
	"encoding/json"
	"fmt"
	"github.com/ministryofjustice/opg-sirius-supervision-firm-deputy-hub/internal/model"
	"net/http"
	"sort"
	"strings"
)

func (c *Client) GetFirmDeputies(ctx Context, firmId int) ([]model.FirmDeputy, error) {
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

	var v []model.Deputies
	if err = json.NewDecoder(resp.Body).Decode(&v); err != nil {
		return nil, err
	}

	var deputies []model.FirmDeputy

	for _, t := range v {
		var deputy = model.FirmDeputy{
			DeputyId:             t.DeputyId,
			Firstname:            t.Firstname,
			Surname:              t.Surname,
			DeputyNumber:         t.DeputyNumber,
			ActiveClientsCount:   getActiveClientCount(t.Orders),
			ExecutiveCaseManager: t.ExecutiveCaseManager.DisplayName,
			OrganisationName:     t.OrganisationName,
			Town:                 t.Town,
			AssuranceType:        t.AssuranceVisit.AssuranceType.Label,
			MarkedAsLabel:        t.AssuranceVisit.VisitReportMarkedAs.Label,
			MarkedAsClass:        strings.ToLower(t.AssuranceVisit.VisitReportMarkedAs.Label),
			ReviewDate:           FormatDateAndTime(DateTimeFormat, t.AssuranceVisit.ReportReviewDate, DateTimeDisplayFormat),
			PanelDeputy:          t.DeputyImportantInformation.PanelDeputy,
		}

		deputies = append(deputies, deputy)
	}
	sortedDeputies := sortTheDeputiesByNumberOfClients(deputies)
	return sortedDeputies, err
}

func getActiveClientCount(orders []model.Orders) int {
	clientId := getListOfClientIds(orders)
	uniqueClientIds := removeDuplicateIDs(clientId)

	return len(uniqueClientIds)
}

func getListOfClientIds(orders []model.Orders) []int {
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

func sortTheDeputiesByNumberOfClients(firmDeputies []model.FirmDeputy) []model.FirmDeputy {
	sort.Slice(firmDeputies, func(i, j int) bool {
		return firmDeputies[i].ActiveClientsCount > firmDeputies[j].ActiveClientsCount
	})
	return firmDeputies
}
