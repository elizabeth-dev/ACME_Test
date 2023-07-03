package grpc_utils

import (
	"github.com/elizabeth-dev/ACME_Test/internal/pkg/utils/query_utils/operators"
	apiV1 "github.com/elizabeth-dev/ACME_Test/pkg/api/v1"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/timestamppb"
	"testing"
	"time"
)

func TestGrpcMappers(t *testing.T) {
	t.Parallel()

	for name, testGroup := range map[string]map[string]func(t *testing.T){
		"grpc filter mapper": {
			"should return the mapped filter": testMapGrpcFilterToFilter,
		},
		"grpc operator mapper": {
			"should return the mapped operator":               testMapGrpcOperatorToOperator,
			"should return the default operator when invalid": testMapGrpcOperatorToOperatorWithInvalidOperator,
		},
		"grpc filter value mapper": {
			"should return the mapped value": testMapGrpcFilterValueToValue,
		},
		"grpc sort mapper": {
			"should return the mapped sort": testMapGrpcSortToSort,
		},
		"grpc direction mapper": {
			"should return the mapped direction":               testMapGrpcDirectionToDirection,
			"should return the default direction when invalid": testMapGrpcDirectionToDirectionWithInvalidDirection,
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

func testMapGrpcFilterToFilter(t *testing.T) {
	filter := &apiV1.Filter{
		Field:    "field",
		Operator: apiV1.Filter_EQUALS,
		Value: &apiV1.Filter_StringValue{
			StringValue: "value",
		},
	}

	out := MapGrpcFilterToFilter(filter)

	assert.Equal(t, filter.Field, out.Field)
	assert.Equal(t, MapGrpcFilterValueToFilterValue(filter), out.Value)
	assert.Equal(t, MapGrpcOperatorToOperator(filter.Operator), out.Operator)
}

func testMapGrpcOperatorToOperator(t *testing.T) {
	out := MapGrpcOperatorToOperator(apiV1.Filter_EQUALS)
	assert.Equal(t, operators.EQUALS, out)

	out = MapGrpcOperatorToOperator(apiV1.Filter_NOT_EQUALS)
	assert.Equal(t, operators.NOT_EQUALS, out)

	out = MapGrpcOperatorToOperator(apiV1.Filter_GREATER_THAN)
	assert.Equal(t, operators.GREATER_THAN, out)

	out = MapGrpcOperatorToOperator(apiV1.Filter_GREATER_THAN_EQ)
	assert.Equal(t, operators.GREATER_THAN_EQ, out)

	out = MapGrpcOperatorToOperator(apiV1.Filter_LESS_THAN)
	assert.Equal(t, operators.LESS_THAN, out)

	out = MapGrpcOperatorToOperator(apiV1.Filter_LESS_THAN_EQ)
	assert.Equal(t, operators.LESS_THAN_EQ, out)
}

func testMapGrpcOperatorToOperatorWithInvalidOperator(t *testing.T) {
	out := MapGrpcOperatorToOperator(-1234)
	assert.Equal(t, operators.EQUALS, out)
}

func testMapGrpcFilterValueToValue(t *testing.T) {
	out := MapGrpcFilterValueToFilterValue(
		&apiV1.Filter{
			Value: &apiV1.Filter_StringValue{
				StringValue: "value",
			},
		},
	)
	assert.Equal(t, "value", out)

	out = MapGrpcFilterValueToFilterValue(
		&apiV1.Filter{
			Value: &apiV1.Filter_IntValue{
				IntValue: 1234,
			},
		},
	)
	assert.Equal(t, int64(1234), out)

	out = MapGrpcFilterValueToFilterValue(
		&apiV1.Filter{
			Value: &apiV1.Filter_BoolValue{
				BoolValue: true,
			},
		},
	)
	assert.Equal(t, true, out)

	out = MapGrpcFilterValueToFilterValue(
		&apiV1.Filter{
			Value: &apiV1.Filter_DoubleValue{
				DoubleValue: 1234.1234,
			},
		},
	)
	assert.Equal(t, 1234.1234, out)

	out = MapGrpcFilterValueToFilterValue(
		&apiV1.Filter{
			Value: &apiV1.Filter_TimestampValue{
				TimestampValue: &timestamppb.Timestamp{
					Seconds: 1665913529,
					Nanos:   328,
				},
			},
		},
	)
	assert.Equal(t, time.Unix(1665913529, 328).UTC(), out)
}

func testMapGrpcSortToSort(t *testing.T) {
	sort := &apiV1.Sort{
		Field:     "field",
		Direction: apiV1.Sort_ASC,
	}

	out := MapGrpcSortToSort(sort)

	assert.Equal(t, sort.Field, out.Field)
	assert.Equal(t, MapGrpcDirectionToDirection(sort.Direction), out.Direction)
}

func testMapGrpcDirectionToDirection(t *testing.T) {
	out := MapGrpcDirectionToDirection(apiV1.Sort_ASC)
	assert.Equal(t, operators.ASC, out)

	out = MapGrpcDirectionToDirection(apiV1.Sort_DESC)
	assert.Equal(t, operators.DESC, out)
}

func testMapGrpcDirectionToDirectionWithInvalidDirection(t *testing.T) {
	out := MapGrpcDirectionToDirection(-1234)
	assert.Equal(t, operators.ASC, out)
}
