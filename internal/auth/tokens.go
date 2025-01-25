package auth

import (
	"log"
	"net/http"
)

func IsCookieSet(r *http.Request, cookieName string) bool {
	cookie, err := r.Cookie(cookieName)
	if err != nil {

		if err == http.ErrNoCookie {
			return false
		}

		log.Println("Error retrieving cookie:", err)
		return false
	}
	
	if cookie.Value == "" {
		return false
	}

	return true
}
