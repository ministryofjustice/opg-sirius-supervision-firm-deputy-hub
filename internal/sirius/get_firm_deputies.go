package sirius

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type orderStatus struct {
	Handle string `json:"handle"`
	Label  string `json:"label"`
}

type client struct {
	Id int `json:"id"`
}

type orders struct {
	Order struct {
		Id          int         `json:"id"`
		Client      client      `json:"client"`
		OrderStatus orderStatus `json:"orderStatus"`
	} `json:"order"`
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
}

type FirmDeputy struct {
	DeputyId             int
	Firstname            string
	Surname              string
	DeputyNumber         int
	ActiveClient         int
	ExecutiveCaseManager string
	OrganisationName     string
}

type FirmDeputiesDetails []FirmDeputy

func (c *Client) GetFirmDeputies(ctx Context, firmId int) (FirmDeputiesDetails, error) {
	req, err := c.newRequest(ctx, http.MethodGet, fmt.Sprintf("/api/v1/firms/%d/deputies", firmId), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	// io.Copy(os.Stdout, resp.Body)
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

	var deputies FirmDeputiesDetails

	for _, t := range v {
		var deputy = FirmDeputy{
			DeputyId:             t.DeputyId,
			Firstname:            t.Firstname,
			Surname:              t.Surname,
			DeputyNumber:         t.DeputyNumber,
			ActiveClient:         getActiveClientCount(t.Orders),
			ExecutiveCaseManager: t.ExecutiveCaseManager.EcmName,
			OrganisationName:     t.OrganisationName,
		}

		deputies = append(deputies, deputy)
	}

	return deputies, err
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
