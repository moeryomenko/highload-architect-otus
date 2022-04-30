package loadtest

import (
	"context"
	"log"

	gofakeit "github.com/brianvoe/gofakeit/v6"
	uuid "github.com/satori/go.uuid"

	"github.com/moeryomenko/highload-architect-otus/social/internal/config"
	"github.com/moeryomenko/highload-architect-otus/social/internal/domain"
	"github.com/moeryomenko/highload-architect-otus/social/internal/repository"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	pool, err := repository.InitConnPool(context.Background(), cfg)
	if err != nil {
		panic(err)
	}

	users := repository.NewUsers(pool)

	for i := 0; i < 1_000_000; i++ {
		err := users.Save(context.Background(), generateUser())
		if err != nil {
			log.Printf("couldn't save user profile, reason: %v", err)
		}
	}
	log.Println("succesful store test data")
}

func generateUser() *domain.User {
	return &domain.User{
		ID: uuid.NewV4(),
		Info: &domain.Profile{
			FirstName: gofakeit.FirstName(),
			LastName:  gofakeit.LastName(),
			Age:       gofakeit.IntRange(18, 50),
			Gender:    domain.Gender(gofakeit.RandomString([]string{domain.Male, domain.Female})),
			City:      gofakeit.City(),
		},
	}
}
