package graph

import (
	"context"
	"fmt"
	"github.com/go-pg/pg/v10"
	"github.com/ivansukach/gql-tutorial/models"
	"net/http"
	"time"
)

const userLoaderKey = "userLoader"

func DataLoaderMiddleware(db *pg.DB, next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		userLoader := UserLoader{
			maxBatch: 100,
			wait:     100 * time.Millisecond,
			fetch: func(ids []string) ([]*models.User, []error) {
				var users []*models.User
				fmt.Println(ids)
				err := db.Model(&users).Where("id IN (?)", pg.In(ids)).Select()
				if err != nil {
					return nil, []error{err}
				}
				orderedUsers := make([]*models.User, len(users))
				for i, id := range ids {
					for _, user := range users {
						if user.ID == id {
							orderedUsers[i] = user
						}
					}
				}
				return orderedUsers, nil
			},
		}
		ctx := context.WithValue(request.Context(), userLoaderKey, &userLoader)
		next.ServeHTTP(writer, request.WithContext(ctx))
	})
}

func getUserLoader(ctx context.Context) *UserLoader {
	return ctx.Value(userLoaderKey).(*UserLoader)
}
