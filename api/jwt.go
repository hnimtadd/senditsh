package api

import (
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

const (
	TokenCookieKey  = "session_token"
	expiresDuration = 15 * time.Minute
)

type Claims struct {
	jwt.RegisteredClaims
	UserId   string `json:"user_id"`
	UserName string `json:"user_name"`
}

type JWTService struct {
	SecretKey []byte
	Api       *ApiHandlerImpl
}

func NewJWTService(SecretKey []byte, api *ApiHandlerImpl) *JWTService {
	service := &JWTService{
		SecretKey: SecretKey,
		Api:       api,
	}
	return service
}

func (s *JWTService) GenerateToken(userId, userName string) (string, error) {
	claims := &Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresDuration)),
		},
		UserId:   userId,
		UserName: userName,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Create the Claims

	tokenString, err := token.SignedString(s.SecretKey)
	if err != nil {
		log.Printf("error: %v\n", err)
		return "", err
	}
	return tokenString, nil
}

func (s *JWTService) VerifyToken(tokenString string) (*Claims, error) {
	var claims = Claims{}
	token, err := jwt.ParseWithClaims(
		tokenString,
		&claims,
		func(t *jwt.Token) (any, error) {
			// Make sure the signing method is HMAC
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unxpected signing method: %v", t.Header["alg"])
			}
			return []byte(s.SecretKey), nil
		},
	)
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, err
	}

	return &claims, nil
}

func (s *JWTService) CheckTokenExpireAndStore(claims *Claims) (bool, error) {
	expired, err := claims.GetExpirationTime()
	if err != nil {
		return false, err
	}
	if expired.Time.Before(time.Now()) {
		return true, nil
	}
	return false, nil
}
func (s *JWTService) GenerateTokenAndStore(
	ctx *fiber.Ctx,
	userID string,
	userName string,
) error {
	token, err := s.GenerateToken(userID, userName)
	if err != nil {
		return err
	}
	log.Println(token)
	cookie := &fiber.Cookie{
		Name:     TokenCookieKey,
		Value:    token,
		Expires:  time.Now().Add(expiresDuration),
		HTTPOnly: true,
	}
	ctx.Cookie(cookie)
	return nil
}

func (s *JWTService) ForgetToken(ctx *fiber.Ctx) error {
	cookie := ctx.Cookies(TokenCookieKey)
	if cookie == "" {
		return nil
	}
	ctx.ClearCookie(TokenCookieKey)
	return nil
}

func (s *JWTService) AuthMiddleware() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		log.Println("JWTAuth")
		cookie := ctx.Cookies(TokenCookieKey)
		// Get the JWT token from the request header
		if cookie == "" {
			return ctx.Next()
			// return ctx.Redirect("/", fiber.StatusUnauthorized)
		}
		claims, err := s.VerifyToken(cookie)
		if err != nil {
			log.Printf("token verification failed: %v\n", err)
			return ctx.Next()
		}

		expired, err := s.CheckTokenExpireAndStore(claims)
		if err != nil {
			log.Printf("Error while check token expired: %v\n", err)
			return ctx.Next()
		}
		if expired {
			if err := s.GenerateTokenAndStore(ctx, claims.UserId, claims.UserName); err != nil {
				log.Printf("regenerate token")
				return ctx.RestartRouting()
			}
		}
		ctx.Locals("claims", claims)

		// log.Println("flashData at middleware: ", flashData)
		usr, err := s.Api.GetUserByUserName(claims.UserName)
		if err != nil {
			return ctx.Next()
		}
		ctx.Locals("user", usr)
		return ctx.Next()
	}
}
func GetClaimsFromContext(ctx *fiber.Ctx) (*Claims, bool) {
	claims, ok := ctx.Locals("claims").(*Claims)
	return claims, ok
}
