package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	// "github.com/gofiber/fiber/v2/middleware/csrf"
	"github.com/hnimtadd/senditsh/utils"
)

type AuthenticationService struct {
	JWT              *JWTService
	AuthorizationURL string
	AccessTokenURL   string
	ClientID         string
	ClientSecret     string
	RedirectURI      string
	Api *ApiHandlerImpl
}

func (s *AuthenticationService) LoginHandler() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		// token := utils.Must[string](csrf.CsrfFromCookie("csrf_token")(ctx))
		//Check if user already login
		// _, ok := GetClaimsFromContext(ctx)
		//
		// if ok {
		// 	return ctx.Redirect("/", fiber.StatusFound)
		// }

		token := "sampleToken"
		params := url.Values{
			"client_id":    []string{s.ClientID},
			"redirect_uri": []string{s.RedirectURI},
			"scope":        []string{"read:user,user:email"},
			"state":        []string{token},
		}
		logger.Info("param", params)
		u := utils.Must[*url.URL](url.ParseRequestURI(s.AuthorizationURL))

		u.RawQuery = params.Encode()
		cookie := &fiber.Cookie{
			Name:     "csrf_token",
			Value:    token,
			Expires:  time.Now().Add(1 * time.Minute),
			HTTPOnly: true,
		}
		ctx.Cookie(cookie)
		logger.Info("redirect", u.String())
		return ctx.Redirect(u.String(), http.StatusFound)
	}
}

func (s *AuthenticationService) LogoutHandler() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		logger.Info("Logout")
		if err := s.JWT.ForgetToken(ctx); err != nil {
			logger.Error("error while logout", err)
			return ctx.Redirect("/", fiber.StatusInternalServerError)
		}
		return ctx.Redirect("/", fiber.StatusSeeOther)
	}
}

// Should verify user success and then insert user to database if needed
func (s *AuthenticationService) CallBackHandler() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		logger.Info("Callback from github")
		val := ctx.Queries()
		code, ok := val["code"]
		if !ok {
			return nil
		}
		logger.Info("code", code)
		state, ok := val["state"]
		if !ok {
			return nil
		}
		expectedState := ctx.Cookies("csrf_token")
		if expectedState == "" {
			logger.Error("msg", fmt.Sprintf("Expected state is empty: %v", expectedState))
			return ctx.SendStatus(fiber.StatusUnauthorized)
		}

		logger.Info("expectedState", expectedState)
		scopes, ok := val["scope"]
		for _, scope := range strings.Split(scopes, " ") {
			if scope == "user:email" {
				logger.Info("user:email granted")
			}
		}

		if state != expectedState {
			logger.Error("msg", fmt.Sprintf("wanted: %v, got: %v", expectedState, state))
			return ctx.SendStatus(fiber.StatusUnauthorized)
		}

		accessToken, err := s.getAccessToken(s.ClientID, s.ClientSecret, code)
		if err != nil {
			return ctx.SendStatus(fiber.StatusUnauthorized)
		}
		logger.Info("accessToken", accessToken)
		usr, err := getUserInfo(accessToken)
		if err != nil {
			return ctx.SendStatus(fiber.StatusUnauthorized)
		}
		logger.Info("user", usr)
		user := &User{
			Username: usr.Login,
		}
		if err := s.Api.CreateUser(user); err != nil{
			logger.Error("err", err.Error())
		}
		


		if err := s.JWT.GenerateTokenAndStore(ctx, fmt.Sprintf("github:%d", usr.ID), usr.Login); err != nil {
			logger.Error("err while generate token", err)
			return ctx.SendStatus(fiber.StatusUnauthorized)
		}
		return ctx.Redirect("/", fiber.StatusFound)
	}
}

type accessTokenResponse struct {
	AccessToken string `json:"access_token"`
	Scope       string `json:"scope"`
	TokenType   string `json:"token_type"`
}

func (s *AuthenticationService) getAccessToken(
	clientID string,
	clientSecret string,
	code string,
) (string, error) {
	params := url.Values{
		"client_id":     []string{clientID},
		"client_secret": []string{clientSecret},
		"code":          []string{code},
	}
	u := utils.Must(url.ParseRequestURI(s.AccessTokenURL))
	u.RawQuery = params.Encode()
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return "", err
	}
	req.Header.Add("Accept", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var accessTokenResponse accessTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&accessTokenResponse); err != nil {
		return "", err
	}
	return accessTokenResponse.AccessToken, nil
}

type githubUser struct {
	ID    int    `json:"id"`
	Login string `json:"login"`
}

func getUserInfo(accessToken string) (githubUser, error) {
	req, err := http.NewRequest("GET", "https://api.github.com/user", nil)
	if err != nil {
		return githubUser{}, err
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	req.Header.Add("X-GitHub-Api-Version", "2022-11-28")
	req.Header.Add("Accept", "application/vnd.github+json")
	// agent.Get("https://api.github.com/user")
	rsp, err := http.DefaultClient.Do(req)
	if err != nil {
		return githubUser{}, err
	}

	defer rsp.Body.Close()
	var u githubUser
	if err := json.NewDecoder(rsp.Body).Decode(&u); err != nil {
		return githubUser{}, err
	}
	return u, nil
}
