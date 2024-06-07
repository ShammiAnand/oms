package auth

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/shammianand/oms/services/common/config"
	"github.com/shammianand/oms/services/common/util"
)

// NOTE: userID referes to customerID
func CreateJWT(secret []byte, userID int) (string, error) {

	expiration := time.Second * time.Duration(config.Envs.TokenExpiry)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID":    strconv.Itoa(userID),
		"expiredAt": time.Now().Add(expiration).Unix(),
	})

	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func WithJWTAuth(handlerFunc http.HandlerFunc, store UserStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// get the token from the user Request
		tokenString := getTokenFromRequest(r)
		// validate the JWT
		token, err := validateToken(tokenString)
		if err != nil {
			log.Printf("failed to validate token %v", err)
			permissionDenied(w)
			return
		}
		if !token.Valid {
			log.Printf("invalid token")
			permissionDenied(w)
			return
		}

		// fetch user id and expiration time
		claims := token.Claims.(jwt.MapClaims)
		userID, _ := strconv.Atoi(claims["userID"].(string))
		expirationTime := int64(claims["expiredAt"].(float64))

		if time.Now().Unix() > expirationTime {
			log.Println("TOKEN EXPIRED")
			permissionDenied(w)
			return
		}

		// get user from user id
		user, err := store.GetUserByID(userID)
		if err != nil {
			log.Printf("failed to get user %v", err)
			permissionDenied(w)
			return
		}

		// set context with userID
		ctx := r.Context()
		ctx = context.WithValue(ctx, "userID", user.ID)
		r = r.WithContext(ctx)

		handlerFunc(w, r)
	}
}

func getTokenFromRequest(r *http.Request) string {
	tokenAuth := r.Header.Get("Authorization")
	return tokenAuth
}

func validateToken(token string) (*jwt.Token, error) {
	return jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(config.Envs.Secret), nil
	})
}

func permissionDenied(w http.ResponseWriter) {
	util.WriteError(w, http.StatusForbidden, fmt.Errorf("permission denied"))
}

func GetUserIdFromContext(ctx context.Context) int {
	userID, ok := ctx.Value("userID").(int)
	if !ok {
		return -1
	}
	return userID

}
