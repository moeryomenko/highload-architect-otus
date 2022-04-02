package repository

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
	"github.com/stretchr/testify/assert"
)

func genFirstName() gopter.Gen {
	return gen.SliceOf(gen.Rune()).Map(func(_ []rune) string { return gofakeit.FirstName() })
}

func genLastName() gopter.Gen {
	return gen.SliceOf(gen.Rune()).Map(func(_ []rune) string { return gofakeit.LastName() })
}

func Test_searchQuery(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	properties := gopter.NewProperties(parameters)

	properties.Property("correct map predicates to query", prop.ForAll(func(firstName, lastName string) bool {
		query, values := searchQuery([]predicate{{"first", firstName}, {"second", lastName}})
		return assert.Equal(t, " first LIKE ? AND second LIKE ? ", query) &&
			reflect.DeepEqual([]any{firstName, lastName}, values)
	}, genFirstName(), genLastName()))

	// corner case.
	query, values := searchQuery(nil)
	assert.Equal(t, "", query)
	assert.Empty(t, values)

	properties.TestingRun(t)
}

func Test_getQuery(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	properties := gopter.NewProperties(parameters)

	properties.Property("take first page query without search predicates", prop.ForAll(func(pageSize int) bool {
		query, values := (&pageQuery{pageSize: pageSize}).getQuery()
		return assert.Equal(t, fmt.Sprintf(paginatedListQuery, ""), query) && reflect.DeepEqual([]any{pageSize}, values)
	}, gen.IntRange(10, 20)))

	properties.Property("take next page query without search predicates", prop.ForAll(func(pageSize int, at time.Time) bool {
		query, values := (&pageQuery{pageSize: pageSize, at: &at}).getQuery()
		return assert.Equal(t, fmt.Sprintf(paginatedListQuery, " WHERE  created_at < ?  "), query) &&
			reflect.DeepEqual([]any{at, pageSize}, values)
	}, gen.IntRange(10, 20), gen.Time()))

	properties.Property("take first page query with search predicates", prop.ForAll(func(pageSize int, firstName, lastName string) bool {
		pageOpts := &pageQuery{
			pageSize: pageSize,
			searchPredicates: []predicate{
				{"first", firstName},
				{"last", lastName},
			},
		}
		query, values := pageOpts.getQuery()
		return assert.Equal(t, fmt.Sprintf(paginatedListQuery, " WHERE  first LIKE ? AND last LIKE ?  "), query) &&
			reflect.DeepEqual([]any{firstName, lastName, pageSize}, values)
	}, gen.IntRange(10, 20), genFirstName(), genLastName()))

	properties.Property("take next page query with search predicates", prop.ForAll(func(pageSize int, at time.Time, firstName, lastName string) bool {
		pageOpts := &pageQuery{
			pageSize: pageSize,
			at:       &at,
			searchPredicates: []predicate{
				{"first", firstName},
				{"last", lastName},
			},
		}
		query, values := pageOpts.getQuery()
		return assert.Equal(t, fmt.Sprintf(paginatedListQuery, " WHERE  first LIKE ? AND last LIKE ? AND created_at < ?  "), query) &&
			reflect.DeepEqual([]any{firstName, lastName, at, pageSize}, values)
	}, gen.IntRange(10, 20), gen.Time(), genFirstName(), genLastName()))

	properties.TestingRun(t)
}
