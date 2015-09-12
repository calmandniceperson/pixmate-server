package http

import (
	"errors"
	"imgturtle/misc"
	"net/http"

	"github.com/gorilla/securecookie"
)

var sCookie = securecookie.New(
	securecookie.GenerateRandomKey(64),
	securecookie.GenerateRandomKey(32),
)

const userSessionCookieName string = "imgt_us"
const userSessionCookieNameUName string = "uname"

func setUserCookie(uName string, response http.ResponseWriter) {
	value := map[string]string{
		userSessionCookieNameUName: uName,
	}
	if encodedCookie, err := sCookie.Encode(userSessionCookieName, value); err == nil {
		cookie := &http.Cookie{
			Name:     userSessionCookieName,
			Value:    encodedCookie,
			Path:     "/",
			HttpOnly: true,
			MaxAge:   5000,
		}
		http.SetCookie(response, cookie)
		misc.PrintMessage(0, "http", "cookie.go", "setUserCookie()", "200. Successfully set user session cookie.")
	} else {
		misc.PrintMessage(1, "http", "cookie.go", "setUserCookie()", "500. Couldn't set user session cookie.\n"+err.Error())
	}
}

func getUserCookieData(req *http.Request) (string, error) {
	var uName string
	if cookie, err := req.Cookie(userSessionCookieName); err == nil {
		cookieValue := make(map[string]string)
		if err = sCookie.Decode(userSessionCookieName, cookie.Value, &cookieValue); err == nil {
			uName = cookieValue[userSessionCookieNameUName]
			return string(uName), nil
		}
		return "", err
	}
	return "", errors.New("Something went wrong trying to read the user session cookie")
}

func clearUserCookie(res http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:   userSessionCookieName,
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	http.SetCookie(res, cookie)
	misc.PrintMessage(0, "http", "cookie.go", "clearUserCookie()", "200. Successfully cleared user session cookie.")
}
