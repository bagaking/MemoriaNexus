package passport

import (
	"net/http"
	"time"
)

func setCookie(w http.ResponseWriter, token string) {
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    token,
		Expires:  time.Now().Add(72 * time.Hour),
		HttpOnly: true, // HttpOnly标志确保Javascript无法读取该cookie
	})
}

func (svr *Service) genJWTTokenAndSetCookie(w http.ResponseWriter, userID uint) (token string, err error) {
	token, err = svr.JWT.GenerateToken(userID)
	if err != nil {
		return "", err
	}
	setCookie(w, token)
	return token, nil
}
