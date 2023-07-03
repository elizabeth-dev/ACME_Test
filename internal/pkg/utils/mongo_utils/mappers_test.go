package mongo_utils

import (
	"github.com/elizabeth-dev/ACME_Test/internal/pkg/utils/query_utils"
	"github.com/elizabeth-dev/ACME_Test/internal/pkg/utils/query_utils/operators"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"testing"
)

func TestMongoMappers(t *testing.T) {
	t.Parallel()

	for name, testGroup := range map[string]map[string]func(t *testing.T){
		"mongo filter mapper": {
			"should return the mapped filter": testMapFilterToBson,
			"should return an empty filter":   testMapFilterToBsonWhenEmpty,
		},
		"mongo operator mapper": {
			"should return the mapped operator": testMapOperatorToMongoOperator,
			"should return empty when invalid":  testMapOperatorToEmptyWhenInvalid,
		},
		"mongo sort mapper": {
			"should return the mapped sort": testMapSortToBson,
			"should return an empty sort":   testMapSortToBsonWhenEmpty,
		},
		"mongo direction mapper": {
			"should return the mapped direction": testMapDirectionToMongoDirection,
			"should return zero when invalid":    testMapDirectionToMongoDirectionWhenInvalid,
		},
	} {
		testGroup := testGroup
		t.Run(
			name, func(t *testing.T) {
				t.Parallel()

				for name, test := range testGroup {
					test := test
					t.Run(
						name, func(t *testing.T) {
							t.Parallel()
							test(t)
						},
					)
				}
			},
		)
	}
}

func testMapFilterToBson(t *testing.T) {
	filter := []query_utils.Filter{
		{
			Field:    "field",
			Operator: operators.EQUALS,
			Value:    "value",
		},
	}

	out := MapFilterToBson(filter)

	bsonFilter := out[filter[0].Field]
	operator := MapComparisonOperatorToBson(filter[0].Operator)
	assert.NotNil(t, bsonFilter)
	assert.NotNil(t, bsonFilter.(bson.M)[operator])
	assert.Equal(t, filter[0].Value, bsonFilter.(bson.M)[operator])
}

func testMapFilterToBsonWhenEmpty(t *testing.T) {
	out := MapFilterToBson([]query_utils.Filter{})
	assert.Empty(t, out)
}

func testMapOperatorToMongoOperator(t *testing.T) {
	out := MapComparisonOperatorToBson(operators.EQUALS)
	assert.Equal(t, "$eq", out)

	out = MapComparisonOperatorToBson(operators.NOT_EQUALS)
	assert.Equal(t, "$ne", out)

	out = MapComparisonOperatorToBson(operators.GREATER_THAN)
	assert.Equal(t, "$gt", out)

	out = MapComparisonOperatorToBson(operators.GREATER_THAN_EQ)
	assert.Equal(t, "$gte", out)

	out = MapComparisonOperatorToBson(operators.LESS_THAN)
	assert.Equal(t, "$lt", out)

	out = MapComparisonOperatorToBson(operators.LESS_THAN_EQ)
	assert.Equal(t, "$lte", out)
}

func testMapOperatorToEmptyWhenInvalid(t *testing.T) {
	out := MapComparisonOperatorToBson(-1234)
	assert.Equal(t, "", out)
}

func testMapSortToBson(t *testing.T) {
	sort := []query_utils.Sort{
		{
			Field:     "field",
			Direction: operators.ASC,
		},
	}

	out := MapSortToBson(sort)

	assert.Equal(t, sort[0].Field, out[0].Key)
	assert.Equal(t, MapSortDirectionToBson(sort[0].Direction), out[0].Value)
}

func testMapSortToBsonWhenEmpty(t *testing.T) {
	out := MapSortToBson([]query_utils.Sort{})
	assert.Empty(t, out)
}

func testMapDirectionToMongoDirection(t *testing.T) {
	out := MapSortDirectionToBson(operators.ASC)
	assert.Equal(t, 1, out)

	out = MapSortDirectionToBson(operators.DESC)
	assert.Equal(t, -1, out)
}

func testMapDirectionToMongoDirectionWhenInvalid(t *testing.T) {
	out := MapSortDirectionToBson(-1234)
	assert.Equal(t, 0, out)
}
