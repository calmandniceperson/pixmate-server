package http

import (
	"errors"
	"net/http"
)

const userSessionCookie string = "imgt_us"
const userSessionCookieUName string = "uname"

func setUserCookie(uName string, response http.ResponseWriter) {
	value := map[string]string{
		userSessionCookieUName: uName,
	}
	if encodedCookie, err := cookieHandler.Encode(userSessionCookie, value); err == nil {
		cookie := &http.Cookie{
			Name:     userSessionCookie,
			Value:    encodedCookie,
			Path:     "/",
			HttpOnly: true,
			MaxAge:   5000,
		}
		http.SetCookie(response, cookie)
	}
}

func getUserCookieData(req *http.Request) (string, error) {
	var uName string
	if cookie, err := req.Cookie(userSessionCookie); err == nil {
		cookieValue := make(map[string]string)
		if err = cookieHandler.Decode(userSessionCookie, cookie.Value, &cookieValue); err == nil {
			uName = cookieValue[userSessionCookieUName]
			return string(uName), nil
		}
		return "", err
	}
	return "", errors.New("Something went wrong trying to read the user session cookie")
}

func clearUserCookie(res http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:   userSessionCookie,
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	http.SetCookie(res, cookie)
}
