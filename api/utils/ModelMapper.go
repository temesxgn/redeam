package utils

import (
	"github.com/fatih/structs"
	"go.mongodb.org/mongo-driver/bson"
	"strings"
)

type ModelMapper struct{}

// ToMongoDocument maps book KV pair to mongo update fields
func (m *ModelMapper) ToMongoDocument(fields []*structs.Field) bson.D {
	fds := bson.M{}
	for _, f := range fields {
		tagValue := f.Tag("bson")
		if !strings.Contains(tagValue, "_id") {
			fds[tagValue] = f.Value()
		}
	}

	return bson.D{{"$set", fds}}
}
