package client

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"neft.web/models"
)

type Links struct {
	Before string `json:"before"`
	Self   string `json:"self"`
	Next   string `json:"next"`
	Last   string `json:"last"`
}

func GeneratePaginationFromRequest(c *gin.Context, count int) (models.Pagination, Links) {
	limit := 2
	page := 1
	sort := "created_at asc"

	for key, value := range c.Request.URL.Query() {
		queryValue := value[len(value)-1]
		switch key {
		case "limit":
			limit, _ = strconv.Atoi(queryValue)
		case "page":
			page, _ = strconv.Atoi(queryValue)
		case "sort":
			sort = queryValue
		}
	}

	pagination := models.Pagination{
		Limit: limit,
		Page:  page,
		Sort:  sort,
	}

	pageCount := int(math.Ceil(float64(count) / float64(limit)))
	if pageCount == 0 {
		pageCount = 1
	}

	if page < 1 || page > pageCount {
		handleError(errors.New("test"), c)
		return models.Pagination{}, Links{}
	}

	var nextPage, beforePage string
	if page >= 2 {
		beforePage = fmt.Sprintf("page=%d", page-1)
	}
	if page < pageCount {
		nextPage = fmt.Sprintf("page=%d", page+1)
	}
	url := c.Request.URL.String()
	pageString := fmt.Sprintf("page=%d", page)
	totalPage := fmt.Sprintf("page=%d", pageCount)
	links := Links{
		Before: strings.Replace(url, pageString, beforePage, -1),
		Self:   url,
		Next:   strings.Replace(url, pageString, nextPage, -1),
		Last:   strings.Replace(url, pageString, totalPage, -1),
	}
	return pagination, links

}
