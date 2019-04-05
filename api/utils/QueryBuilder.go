package utils

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/http"
	"strconv"
	"strings"
)

type QueryBuilder struct{}

const (
	defaultSize  = 10
	defaultOrder = "asc"
)

// GetQueryParams returns mongodb filters and findOptions to restrict query results
func (builder *QueryBuilder) GetQueryParams(r *http.Request) (bson.M, *options.FindOptions) {
	queries := r.URL.Query()
	filters := bson.M{}
	findOptions := options.Find().SetLimit(defaultSize)

	// TODO Split into multiple functions and use custom Query model in entity_models.go file
	for k, v := range queries {
		switch k {
		case "size":
			size, _ := strconv.ParseInt(r.URL.Query().Get("size"), 10, 64)
			if size == 0 {
				size = 10
			}

			findOptions = findOptions.SetLimit(size)
		case "page":
			page, _ := strconv.ParseInt(v[0], 10, 64)
			if page < 1 {
				page = 1
			}

			skip := *findOptions.Limit * (page - 1)
			findOptions = findOptions.SetSkip(skip)
		case "author":
			filters["author"] = strings.Replace(v[0], "+", " ", -1)

		case "status", "rating":
			value, _ := strconv.ParseInt(v[0], 10, 64)
			filters[k] = value
		}
	}

	return filters, findOptions
}
