package grpc_utils

import (
	"github.com/elizabeth-dev/FACEIT_Test/internal/pkg/utils/query_utils"
	"github.com/elizabeth-dev/FACEIT_Test/internal/pkg/utils/query_utils/operators"
	apiV1 "github.com/elizabeth-dev/FACEIT_Test/pkg/api/v1"
)

func MapGrpcFilterToFilter(filter *apiV1.Filter) query_utils.Filter {
	return query_utils.Filter{
		Field:    filter.Field,
		Value:    MapGrpcFilterValueToFilterValue(filter),
		Operator: MapGrpcOperatorToOperator(filter.Operator),
	}
}

func MapGrpcOperatorToOperator(operator apiV1.Filter_Operator) operators.Comparison {
	switch operator {
	case apiV1.Filter_EQUALS:
		return operators.EQUALS
	case apiV1.Filter_NOT_EQUALS:
		return operators.NOT_EQUALS
	case apiV1.Filter_GREATER_THAN:
		return operators.GREATER_THAN
	case apiV1.Filter_GREATER_THAN_EQ:
		return operators.GREATER_THAN_EQ
	case apiV1.Filter_LESS_THAN:
		return operators.LESS_THAN
	case apiV1.Filter_LESS_THAN_EQ:
		return operators.LESS_THAN_EQ
	default:
		return operators.EQUALS // If an invalid operator is passed, we default to EQUALS
	}
}

func MapGrpcFilterValueToFilterValue(filter *apiV1.Filter) interface{} {
	switch filter.Value.(type) {
	case *apiV1.Filter_StringValue:
		return filter.GetStringValue()
	case *apiV1.Filter_IntValue:
		return filter.GetIntValue()
	case *apiV1.Filter_BoolValue:
		return filter.GetBoolValue()
	case *apiV1.Filter_DoubleValue:
		return filter.GetDoubleValue()
	case *apiV1.Filter_TimestampValue:
		return filter.GetTimestampValue().AsTime()
	default:
		return nil
	}
}

func MapGrpcSortToSort(sort *apiV1.Sort) query_utils.Sort {
	return query_utils.Sort{
		Field:     sort.Field,
		Direction: MapGrpcDirectionToDirection(sort.Direction),
	}
}

func MapGrpcDirectionToDirection(direction apiV1.Sort_Direction) operators.Sort {
	switch direction {
	case apiV1.Sort_ASC:
		return operators.ASC
	case apiV1.Sort_DESC:
		return operators.DESC
	default:
		return operators.ASC // If an invalid direction is passed, we default to ASC
	}
}
