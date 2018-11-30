package handler

import (
	"context"
	"fmt"
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
)

// An AuthContextKey represents something stored in a context.Context.Values
// It's meant to ensure context keys are unique to the application (so eg: "user" does not collide with another library)
type authContextKey string

var contextKeyAuthedUser = authContextKey("authedUser")

func GetMemberEmailFromContext(ctx context.Context) (string, error) {
	emailRaw := ctx.Value(contextKeyAuthedUser)
	email, ok := emailRaw.(string)
	if !ok {
		return "", fmt.Errorf("Did not find email in context")
	}
	return email, nil
}

func getEmailFromToken(jwtToken string, authConfig *Auth0Config) (string, error) {
	// Parse token and verify it has the correct algorithm.
	// The parser takes a function that validates the public metadata and then returns the token secret
	// Thus the public information must be correct before the private token is unpacked with the secret
	parsedToken, err := jwt.Parse(jwtToken, func(tk *jwt.Token) (interface{}, error) {
		if tk.Method.Alg() != jwt.SigningMethodHS512.Alg() {
			return nil, fmt.Errorf("Unexpected signing method: %v", tk.Method.Alg())
		}
		return []byte(authConfig.JWTSecret), nil
	})
	if err != nil {
		return "", fmt.Errorf("Failed to parse token: %s", err)
	}

	// Refuse the data if the token is not valid
	// During the parse, the token and claims were validated; this is what this is checking
	if !parsedToken.Valid {
		return "", fmt.Errorf("The parsed token was invalid. Cannot trust its content")
	}

	// Now that we have a legal token, we can trust its data, so extract the claims from it
	// This error probably can't happen, since failing to parse a jwt.MapClaims fails earlier at the "Failed to parse token" error
	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("Failed to extract claims from token")
	}

	// Extract the claims into an AuthedUser, and return
	email, ok := claims["email"].(string)
	if !ok {
		return "", fmt.Errorf("Failed to parse email inside token as a string")
	}
	return email, nil
}

func clearCookie(w http.ResponseWriter, cookieName string) {
	http.SetCookie(w, &http.Cookie{
		Name:     cookieName,
		Value:    "",
		MaxAge:   -1,
		HttpOnly: true,
		Path:     "/",
	})
}

func Authenticated(next http.Handler, authConfig *Auth0Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var err error

		// Look for an auth JWT token
		var cookie *http.Cookie
		cookie, err = r.Cookie(authConfig.CookieName)
		if err == http.ErrNoCookie {
			id := LogId()
			logger.Infow("No cookie found when authenticating user for request. Assuming they just need to login.", "CookieName", authConfig.CookieName, "id", id)

			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprintf(w, "You need to login before performing this action")
			return
		}
		if err != nil {
			id := LogId()
			logger.Infow("Unkown cookie error occured while checking authentication", "error", err, "id", id)
			fmt.Fprintf(w, "Authentication error. Please give this code to an admin, and try logging in again: %s", id)

			// Tell the browser to clear the cookie, so we don't keep having this issue
			clearCookie(w, authConfig.CookieName)
			return
		}
		token := cookie.Value

		// Validate the token and get the email out of it
		var email string
		email, err = getEmailFromToken(token, authConfig)
		if err != nil {
			id := LogId()
			logger.Infow("Failed to get email from token while checking authentication", "error", err, "id", id)
			fmt.Fprintf(w, "Authentication error. Please give this code to an admin, and try logging in again: %s", id)

			// Tell the browser to clear the cookie, so we don't keep having this issue
			clearCookie(w, authConfig.CookieName)
			return
		}

		// No error means auth successful, proceed to handler
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), contextKeyAuthedUser, email)))
	})
}
