package auth

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Config struct {
	JWTSecret []byte `mapstructure:"jwtsecret"`
	Lifetime  int    `mapstructure:"lifetime"`
	Issuer    string `mapstructure:"issuer"`
}

type AppClaims struct {
	jwt.RegisteredClaims
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
}

type Auth struct {
	cfg Config
}

func MustLoad(cfg Config) *Auth {
	return &Auth{cfg: cfg}
}

func (a *Auth) GenerateJWT(userID int, username string) (string, error) {
	claims := AppClaims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * time.Duration(a.cfg.Lifetime))),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    a.cfg.Issuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(a.cfg.JWTSecret)
}

var excluded = []*regexp.Regexp{
	regexp.MustCompile(`^/api/v\d+/users/register$`),
	regexp.MustCompile(`^/api/v\d+/users/login$`),
	regexp.MustCompile(`^/status$`),
	regexp.MustCompile(`^/docs.*$`),
	regexp.MustCompile(`^/metrics$`),
}

func isExcluded(path string) bool {
	for _, re := range excluded {
		if re.MatchString(path) {
			return true
		}
	}
	return false
}

func (a *Auth) MiddlewareHandleJWT(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		path := req.URL.Path

		if isExcluded(path) {
			next.ServeHTTP(res, req)
			return
		}

		raw := req.Header.Get("Authorization")
		if len(raw) == 0 {
			res.Header().Set("WWW-Authenticate", `Bearer`)
			http.Error(res, "Must provide credentials for this request, use Authorization 'Bearer [token]'", http.StatusUnauthorized)
			return
		}

		if !strings.HasPrefix(raw, "Bearer ") {
			res.Header().Set("WWW-Authenticate", `Bearer`)
			http.Error(res, "Authorization must be 'Bearer [token]'", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(raw, "Bearer ")

		token, err := jwt.ParseWithClaims(tokenString, &AppClaims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return a.cfg.JWTSecret, nil
		})

		if err != nil || !token.Valid {
			res.Header().Set("WWW-Authenticate", `Bearer`)
			http.Error(res, "Token not valid", http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(*AppClaims)
		if !ok {
			res.Header().Set("WWW-Authenticate", `Bearer`)
			http.Error(res, "Invalid claims", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(req.Context(), "user_id", claims.UserID)
		next.ServeHTTP(res, req.WithContext(ctx))
	})
}
