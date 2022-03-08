package graph

import (
	"context"
	"github.com/99designs/gqlgen/graphql"
	"github.com/ivansukach/gql-tutorial/validator"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

func Validation(ctx context.Context, v validator.Validation) bool {
	valid, errors := v.Validate()
	if !valid {
		for field, err := range errors {
			graphql.AddError(ctx, &gqlerror.Error{
				Message: err,
				Extensions: map[string]interface{}{
					"field": field,
				},
			})
		}
	}

	return valid
}
