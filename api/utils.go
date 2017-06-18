package api

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
)

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func authenticated(req *http.Request) bool {
	token := req.Header.Get("Token")
	if token == "" {
		return false
	}
	fmt.Println(token)
	fmt.Println(sessions)
	//validate token
	if _, ok := sessions[token]; ok {
		return true
	}
	return false
}

func generateSessionId(s int) (string, error) {
	b, err := generateRandomBytes(s)
	return base64.URLEncoding.EncodeToString(b), err
}

func generateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}
