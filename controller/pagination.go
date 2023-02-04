package controller

import (
	"encoding/json"
	"fmt"
	engine "github.com/JoanGTSQ/api"
	"math"
	"neft.web/models"
)

type Links struct {
	Before string `json:"before"`
	Self   string `json:"self"`
	Next   string `json:"next"`
	Last   string `json:"last"`
}

func (client *Client) GeneratePaginationFromRequest(count int) (models.Pagination, Links) {

	pagination := models.Pagination{
		Limit: 2,
		Page:  1,
		Sort:  "created_at asc",
	}

	jsonStr, err := json.Marshal(client.IncomingMessage.Data["pagination"])
	if err != nil {
		engine.Debug.Println(err)
		client.LastMessage.Data["error"] = err.Error()
		client.SendMessage()
		return models.Pagination{}, Links{}
	}

	// Obtain the body in the request and parse to the user
	if err := json.Unmarshal(jsonStr, &pagination); err != nil {
		engine.Warning.Println(err)
		client.LastMessage.Data["error"] = engine.ERR_INVALID_JSON
		client.SendMessage()
		return models.Pagination{}, Links{}
	}

	pageCount := int(math.Ceil(float64(count) / float64(pagination.Limit)))
	if pageCount == 0 {
		pageCount = 1
	}

	if pagination.Page < 1 || pagination.Page > pageCount {
		engine.Warning.Println("could not translate the pagination request", pagination.Page, pageCount)
		return models.Pagination{}, Links{}
	}

	var nextPage, beforePage string
	if pagination.Page >= 2 {
		beforePage = fmt.Sprintf("page=%d", pagination.Page-1)
	}
	if pagination.Page < pageCount {
		nextPage = fmt.Sprintf("page=%d", pagination.Page+1)
	}
	pageString := fmt.Sprintf("page=%d", pagination.Page)
	totalPage := fmt.Sprintf("page=%d", pageCount)
	links := Links{
		Before: beforePage,
		Self:   pageString,
		Next:   nextPage,
		Last:   totalPage,
	}
	return pagination, links

}
