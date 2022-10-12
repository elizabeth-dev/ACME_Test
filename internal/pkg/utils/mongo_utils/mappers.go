package mongo_utils

import (
	"github.com/elizabeth-dev/FACEIT_Test/internal/pkg/utils/query_utils"
	"github.com/elizabeth-dev/FACEIT_Test/internal/pkg/utils/query_utils/operators"
	"go.mongodb.org/mongo-driver/bson"
)

/* Comparison mappers */

func MapFilterToBson(filter []query_utils.Filter) bson.M {
	// When applying a filter, the document order for the filters doesn't matter, so we use bson.M
	var bsonFilter bson.M

	for _, f := range filter {
		operator := MapComparisonOperatorToBson(f.Operator)

		if operator != "" {
			bsonFilter[f.Field] = bson.M{operator: f.Value}
		}
	}

	return bsonFilter
}

func MapComparisonOperatorToBson(op operators.Comparison) string {
	switch op {
	case operators.EQUALS:
		return "$eq"
	case operators.NOT_EQUALS:
		return "$ne"
	case operators.GREATER_THAN:
		return "$gt"
	case operators.GREATER_THAN_EQ:
		return "$gte"
	case operators.LESS_THAN:
		return "$lt"
	case operators.LESS_THAN_EQ:
		return "$lte"
	default:
		return ""
	}
}

/* Sorting mappers */

func MapSortToBson(sort []query_utils.Sort) bson.D {
	// When sorting, the document order for the sorters matters, because the first sort parameter has the highest priority, so we use bson.D
	var bsonSort bson.D

	for _, s := range sort {
		operator := MapSortDirectionToBson(s.Direction)

		if operator != 0 {
			bsonSort = append(bsonSort, bson.E{Key: s.Field, Value: operator})
		}
	}

	return bsonSort
}

func MapSortDirectionToBson(op operators.Sort) int {
	switch op {
	case operators.ASC:
		return 1
	case operators.DESC:
		return -1
	default:
		return 0
	}
}
