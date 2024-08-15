package jwt

import (
	"aid-server/configs"
	"aid-server/pkg/res"
	"errors"
	"github.com/golang-jwt/jwt/v4"
	jwtlib "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"time"
)

type UserClaims struct {
	ID string `json:"id"`
	jwtlib.RegisteredClaims
}

func GenerateToken(id string) (string, error) {
	Claims := UserClaims{
		ID: id,
		RegisteredClaims: jwtlib.RegisteredClaims{
			ExpiresAt: jwtlib.NewNumericDate(time.Now().Add(configs.Configs.Jwt.Duration)),
			IssuedAt:  jwtlib.NewNumericDate(time.Now()),
			Subject:   "auth",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &Claims)
	return token.SignedString([]byte(configs.Configs.Jwt.Secret))
}

func (uc *UserClaims) Valid() error {
	if time.Now().After(uc.ExpiresAt.Time) {
		return errors.New("token is expired")
	}
	if uc.Subject != "auth" {
		return errors.New("invalid subject")
	}
	// parse id to UUID
	if _, err := uuid.Parse(uc.ID); err != nil {
		return errors.New("invalid user id")
	}
	return nil
}

func ParseToken(tokenString string) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(configs.Configs.Jwt.Secret), nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*UserClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}
	return claims, nil
}

// GenerateParseJwtMiddle is a middleware function to parse jwt token
// and set claims to context by key "claims"
func GenerateParseJwtMiddle(resFunc func(bool, string) res.Response) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			tokenString := c.Request().Header.Get("Authorization")
			if tokenString == "" {
				return c.JSON(401, resFunc(false, "no token"))
			}
			claims, err := ParseToken(tokenString)
			if err != nil {
				return c.JSON(401, resFunc(false, "invalid token"))
			}
			if err := claims.Valid(); err != nil {
				return c.JSON(401, resFunc(false, err.Error()))
			}
			c.Set("claims", claims)
			return next(c)
		}
	}
}
