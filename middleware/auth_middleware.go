package middleware

import (
	"context"
	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"github.com/ivansukach/gql-tutorial/models"
	"github.com/ivansukach/gql-tutorial/repository"
	"github.com/pkg/errors"
	"net/http"
	"os"
	"strings"
)

const bearerPrefix = "Bearer "

var errUserIsAbsentInContext = errors.New("user is absent in the context")

var authHeaderExtractor = &request.PostExtractionFilter{
	Extractor: request.HeaderExtractor{"Authorization"},
	Filter:    stripBearerPrefixFromToken,
}

func stripBearerPrefixFromToken(token string) (string, error) {
	if !strings.HasPrefix(token, bearerPrefix) {
		return "", errors.New("Invalid token value in Authorization HTTP header")
	}
	return token[len(bearerPrefix):], nil
}

var authExtractor = &request.MultiExtractor{
	authHeaderExtractor,
	request.ArgumentExtractor{"access_token"},
}

const userKey = "userData"

func AuthMiddleware(repo repository.UsersRepo) func(handler http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token, err := parseToken(r)
			if err != nil {
				//if we delete this line, we won't get UI
				next.ServeHTTP(w, r)
				//log error
				//or we could use next.ServeHTTP to go next, because of not every handler require authentication
				//it seems like http middleware useless in this cases
				return
			}
			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok && !token.Valid {
				//if we delete this line, we won't get UI
				next.ServeHTTP(w, r)
				//log error
				//or we could use next.ServeHTTP to go next, because of not every handler require authentication
				//it seems like http middleware useless in this cases
				return
			}
			user, err := repo.GetUserByID(claims["jti"].(string))
			if err != nil {
				//if we delete this line, we won't get UI
				next.ServeHTTP(w, r)
				//log error
				//or we could use next.ServeHTTP to go next, because of not every handler require authentication
				//it seems like http middleware useless in this cases
				return
			}
			ctx := context.WithValue(r.Context(), userKey, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func parseToken(r *http.Request) (*jwt.Token, error) {
	jwtToken, err := request.ParseFromRequest(r, authExtractor, func(token *jwt.Token) (interface{}, error) {
		jwtSecret := []byte(os.Getenv("JWT_SECRET"))
		return jwtSecret, nil
	})

	return jwtToken, errors.Wrap(err, "parseToken error: ")
}

func GetCurrentUserFromCTX(ctx context.Context) (*models.User, error) {
	if ctx.Value(userKey) == nil {
		return nil, errUserIsAbsentInContext
	}
	//check if it is possible to have a pointer to empty data in that case
	user, ok := ctx.Value(userKey).(*models.User)
	if !ok || user.ID == "" {
		return nil, errUserIsAbsentInContext
	}
	return user, nil
}
