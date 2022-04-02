package repository

import (
	"testing"

	uuid "github.com/satori/go.uuid"

	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
)

func genUUID() gopter.Gen {
	return gen.SliceOfN(20, gen.Rune()).Map(func(r []rune) string { return uuid.NewV4().String() })
}

func Test_uuid(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	properties := gopter.NewProperties(parameters)

	properties.Property("invariant uuid conversion", prop.ForAll(func(uuidStr string) bool {
		id, err := uuid.FromString(uuidStr)
		if err != nil {
			t.Fatal(err)
		}

		return binaryToUUID(uuidToBinary(id)).String() == id.String()
	}, genUUID()))

	properties.TestingRun(t)
}
