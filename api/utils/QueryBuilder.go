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
	defaultPage  = 1
	defaultOrder = 1
	defaultSort  = "_id"
)

// GetQueryParams returns mongodb filters and findOptions to restrict query results
func (builder *QueryBuilder) GetQueryParams(r *http.Request) (bson.M, *options.FindOptions) {
	queries := r.URL.Query()
	filters := bson.M{}
	findOptions := options.Find().SetLimit(defaultSize).SetSkip(defaultPage).SetSort(bson.M{defaultSort: defaultOrder})

	// TODO Split into multiple functions and use custom Query model in entity_models.go file
	for k, v := range queries {
		switch k {
		case "size":
			if size, _ := strconv.ParseInt(v[0], 10, 64); size > 0 {
				findOptions = findOptions.SetLimit(size)
			}

		case "page":
			if page, _ := strconv.ParseInt(v[0], 10, 64); page > 0 {
				findOptions = findOptions.SetSkip(page)
			}

		case "author":
			filters["author"] = strings.Replace(v[0], "+", " ", -1)

		case "status", "rating":
			value, _ := strconv.ParseInt(v[0], 10, 64)
			filters[k] = value
		}
	}

	skip := *findOptions.Limit * (*findOptions.Skip - 1)
	findOptions = findOptions.SetSkip(skip)
	return filters, findOptions
}
