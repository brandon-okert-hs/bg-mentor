package handler

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/bcokert/bg-mentor/internal/pkg/database"
	"github.com/bcokert/bg-mentor/internal/pkg/requesterror"
	"github.com/bcokert/bg-mentor/pkg/model"
	"github.com/bcokert/bg-mentor/pkg/request"
	jwt "github.com/dgrijalva/jwt-go"
	"golang.org/x/oauth2"
)

type Auth0Config struct {
	Domain          string
	ClientID        string
	ClientSecret    string
	RedirectURLRoot string
	JWTSecret       []byte
	CookieName      string
	oauthConfig     *oauth2.Config
}

type AuthHandler struct {
	AuthConfig    *Auth0Config
	DB            *database.Database
	pendingStates map[string]bool
}

func (h *AuthHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	var head string
	head, req.URL.Path = request.ShiftURL(req.URL.Path)

	switch head {
	case "login":
		switch req.Method {
		case http.MethodGet:
			h.GetLogin(res, req)
			return
		default:
			RespondJSON("Status", http.StatusNotFound, requesterror.MethodNotFound("Auth", head, req), res)
			return
		}
	case "confirm":
		switch req.Method {
		case http.MethodGet:
			h.GetConfirm(res, req)
			return
		default:
			RespondJSON("Status", http.StatusNotFound, requesterror.MethodNotFound("Auth", head, req), res)
			return
		}
	case "logout":
		switch req.Method {
		case http.MethodGet:
			h.GetLogout(res, req)
			return
		default:
			RespondJSON("Status", http.StatusNotFound, requesterror.MethodNotFound("Auth", head, req), res)
			return
		}
	default:
		RespondJSON("Status", http.StatusNotFound, requesterror.PathNotFound("Auth", head, req), res)
		return
	}
}

func (h *AuthHandler) getOauthConfig() *oauth2.Config {
	if h.AuthConfig.oauthConfig == nil {
		h.AuthConfig.oauthConfig = &oauth2.Config{
			ClientID:     h.AuthConfig.ClientID,
			ClientSecret: h.AuthConfig.ClientSecret,
			RedirectURL:  h.AuthConfig.RedirectURLRoot + "/auth/confirm",
			Scopes:       []string{"openid", "profile"},
			Endpoint: oauth2.Endpoint{
				AuthURL:  "https://" + h.AuthConfig.Domain + "/authorize",
				TokenURL: "https://" + h.AuthConfig.Domain + "/oauth/token",
			},
		}
	}

	return h.AuthConfig.oauthConfig
}

func createOauthState() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.StdEncoding.EncodeToString(b)
}

func (h *AuthHandler) addStateToPending(state string) {
	if h.pendingStates == nil {
		h.pendingStates = make(map[string]bool)
	}

	h.pendingStates[state] = true
}

func (h *AuthHandler) isStatePending(state string) bool {
	_, ok := h.pendingStates[state]
	return ok
}

func (h *AuthHandler) removeStateFromPending(state string) {
	if h.pendingStates != nil {
		delete(h.pendingStates, state)
	}
}

func (h *AuthHandler) GetLogin(w http.ResponseWriter, r *http.Request) {
	// Generate random state, works like a salt to confirm redirects
	state := createOauthState()

	// Store pending state in state list
	h.addStateToPending(state)

	logger.Infow("Logging in a user", "state", state)

	// Get oauth config
	conf := h.getOauthConfig()
	aud := "https://" + h.AuthConfig.Domain + "/userinfo"

	// Redirect to the Auth0 login
	audience := oauth2.SetAuthURLParam("audience", aud)
	url := conf.AuthCodeURL(state, audience)
	http.Redirect(w, r, url, http.StatusFound)
}

func (h *AuthHandler) GetConfirm(w http.ResponseWriter, r *http.Request) {
	var err error

	// Compare the state to the token state to prevent man in the middle attack
	stateFromAuth0 := r.URL.Query().Get("state")
	if !h.isStatePending(stateFromAuth0) {
		id := LogId()
		logger.Errorw("State received was not pending confirmation, when confirming login", "stateFromAuth0", stateFromAuth0, "pendingStates", h.pendingStates, "id", id)
		fmt.Fprintf(w, "Login Confirmation Error. Please give this code to an admin: %s", id)
		return
	}

	logger.Infow("Confirming user login via Auth0", "stateFromAuth0", stateFromAuth0)

	// Get oauth config
	conf := h.getOauthConfig()

	// Convert oauth code into an oauthToken
	oauthCode := r.URL.Query().Get("code")
	var oauthToken *oauth2.Token
	oauthToken, err = conf.Exchange(context.TODO(), oauthCode)
	if err != nil {
		id := LogId()
		logger.Errorw("Error converting oauth code to oauth token when confirming login", "oauthCode", oauthCode, "error", err, "id", id)
		fmt.Fprintf(w, "Login Confirmation Error. Please give this code to an admin: %s", id)
		return
	}

	// Get the user info by using the oauth token
	client := conf.Client(context.TODO(), oauthToken)
	var oauthResp *http.Response
	oauthResp, err = client.Get("https://" + h.AuthConfig.Domain + "/userinfo")
	if err != nil {
		id := LogId()
		logger.Errorw("Error fetching user info from Auth0 with oauth token, while confirming login", "url", "https://"+h.AuthConfig.Domain+"/userinfo", "error", err, "id", id)
		fmt.Fprintf(w, "Login Confirmation Error. Please give this code to an admin: %s", id)
		return
	}
	defer oauthResp.Body.Close()

	// Extract the user info into a map
	var userInfo map[string]interface{}
	if err = json.NewDecoder(oauthResp.Body).Decode(&userInfo); err != nil {
		id := LogId()
		logger.Errorw("Error extracting user info from Auth0 oauth response, while confirming login", "error", err, "id", id)
		fmt.Fprintf(w, "Login Confirmation Error. Please give this code to an admin: %s", id)
		return
	}

	// Create the user if they don't exist yet
	email, ok1 := userInfo["name"].(string)
	name, ok2 := userInfo["nickname"].(string)
	avatarURL, ok3 := userInfo["picture"].(string)
	if !ok1 || !ok2 || !ok3 {
		id := LogId()
		logger.Errorw("Error extracting user data from Auth0 oauth response, while confirming login", "error", err, "email", userInfo["name"], "name", userInfo["nickname"], "avatarURL", userInfo["picture"], "id", id)
		fmt.Fprintf(w, "Login Confirmation Error. Please give this code to an admin: %s", id)
		return
	}
	post := &model.MemberPost{
		Email:     email,
		Name:      name,
		AvatarURL: avatarURL,
	}
	_, err = h.DB.CreateMemberIfNotExisting(post)
	if err != nil {
		id := LogId()
		logger.Errorw("Error creating or fetching member after login, while confirming login", "error", err, "post", post, "id", id)
		fmt.Fprintf(w, "Login Confirmation Error. Please give this code to an admin: %s", id)
		return
	}

	// Add the email to the token
	token := jwt.New(jwt.SigningMethodHS512)
	claims := token.Claims.(jwt.MapClaims)
	claims["email"] = email
	tokenString, err := token.SignedString(h.AuthConfig.JWTSecret)
	if err != nil {
		id := LogId()
		logger.Errorw("Login Error while signing new jwt token, while confirming login", "error", err, "email", email, "id", id)
		fmt.Fprintf(w, "Login Error. Please give this code to an admin: %s", id)
		return
	}

	// Tell browser to store the new JWT token
	http.SetCookie(w, &http.Cookie{
		Name:     h.AuthConfig.CookieName,
		Value:    tokenString,
		MaxAge:   604800,
		HttpOnly: true,
		Path:     "/",
	})

	// Redirect to index, now that they're logged in
	http.Redirect(w, r, "/", http.StatusFound)
}

func (h *AuthHandler) GetLogout(w http.ResponseWriter, r *http.Request) {
	var URL *url.URL
	URL, err := URL.Parse("https://" + h.AuthConfig.Domain + "/v2/logout")
	if err != nil {
		id := LogId()
		logger.Errorw("Logout error while creating url for logout", "error", err, "domain", h.AuthConfig.Domain, "id", id)
		fmt.Fprintf(w, "Logout Error. Please give this code to an admin: %s", id)
		return
	}

	parameters := url.Values{}
	parameters.Add("returnTo", h.AuthConfig.RedirectURLRoot+"/")
	parameters.Add("client_id", h.AuthConfig.ClientID)
	URL.RawQuery = parameters.Encode()

	clearCookie(w, h.AuthConfig.CookieName)

	http.Redirect(w, r, URL.String(), http.StatusTemporaryRedirect)
}
