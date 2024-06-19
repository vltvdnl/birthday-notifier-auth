package jwt

import (
	"fmt"
	"sso/internal/domain/models"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type JWTHandler struct{}

func New() *JWTHandler {
	return &JWTHandler{}
}
func (t *JWTHandler) NewToken(user models.User, app models.App, duration time.Duration) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["uid"] = user.ID
	claims["email"] = user.Email
	claims["exp"] = time.Now().Add(duration).Unix()
	claims["app_id"] = app.ID
	fmt.Println(app.Secret)
	if token_string, err := token.SignedString([]byte(app.Secret)); err != nil {
		return "", err
	} else {
		return token_string, nil
	}
}
func (t *JWTHandler) VerifyPassword(hash []byte, passhash []byte) error {
	if err := bcrypt.CompareHashAndPassword(hash, passhash); err != nil {
		return err
	}
	return nil
}
func (t *JWTHandler) GenerateHash(password []byte) ([]byte, error) {
	passHash, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return passHash, nil
}
