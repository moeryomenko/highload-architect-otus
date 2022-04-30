package repository

import (
	"fmt"
	// "reflect"
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
		query := searchQuery([]predicate{{"first", firstName}, {"second", lastName}})
		return assert.Equal(t, fmt.Sprintf(" first LIKE '%s' AND second LIKE '%s' ", firstName, lastName), query)
	}, genFirstName(), genLastName()))

	// corner case.
	query := searchQuery(nil)
	assert.Equal(t, "", query)

	properties.TestingRun(t)
}

func Test_getQuery(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	properties := gopter.NewProperties(parameters)

	properties.Property("take first page query without search predicates", prop.ForAll(func(pageSize int) bool {
		query := (&pageQuery{pageSize: pageSize}).getQuery()
		return assert.Equal(t, fmt.Sprintf(paginatedListQuery, "", pageSize), query)
	}, gen.IntRange(10, 20)))

	properties.Property("take next page query without search predicates", prop.ForAll(func(pageSize int, at time.Time) bool {
		query := (&pageQuery{pageSize: pageSize, at: &at}).getQuery()
		return assert.Equal(t, fmt.Sprintf(paginatedListQuery, fmt.Sprintf(" WHERE  created_at < '%s'  ", at.Format("2006-01-02 15:04:05")), pageSize), query)
	}, gen.IntRange(10, 20), gen.Time()))

	properties.Property("take first page query with search predicates", prop.ForAll(func(pageSize int, firstName, lastName string) bool {
		pageOpts := &pageQuery{
			pageSize: pageSize,
			searchPredicates: []predicate{
				{"first", firstName},
				{"last", lastName},
			},
		}
		query := pageOpts.getQuery()
		return assert.Equal(t, fmt.Sprintf(paginatedListQuery, fmt.Sprintf(" WHERE  first LIKE '%s' AND last LIKE '%s'  ", firstName, lastName), pageSize), query)
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
		query := pageOpts.getQuery()
		return assert.Equal(t, fmt.Sprintf(paginatedListQuery, fmt.Sprintf(" WHERE  first LIKE '%s' AND last LIKE '%s' AND created_at < '%s'  ", firstName, lastName, at.Format("2006-01-02 15:04:05")), pageSize), query)
	}, gen.IntRange(10, 20), gen.Time(), genFirstName(), genLastName()))

	properties.TestingRun(t)
}
