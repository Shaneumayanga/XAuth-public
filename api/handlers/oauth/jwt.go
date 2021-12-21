package oauth

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

//TODO get from env
var key string

func init() {
	if os.Getenv("JWT_SECRET") == "" {
		key = "temporarykey"
	} else {
		key = os.Getenv("JWT_SECRET")
	}
}

func GenerateJWTWithCode(code string) string {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"code": code,
		"exp":  time.Now().Add(time.Minute * 10).Unix(),
	})

	tokenString, err := token.SignedString([]byte(key))
	if err != nil {
		fmt.Printf("err.Error(): %v\n", err.Error())
		return ""
	}
	fmt.Printf("tokenString: %v\n", tokenString)
	return tokenString
}

func ValidateTokenAndGetCode(tokenReq string) string {
	token, err := jwt.Parse(tokenReq, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(key), nil
	})
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		fmt.Println(claims["code"])
		return claims["code"].(string)
	} else {
		fmt.Println("Error in ValidateTokenAndGetCode :" + err.Error())
		return ""
	}
}

func GenerateJWTAccessTokenWithUserId(userid string) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userid": userid,
		"exp":    time.Now().Add(time.Minute * 60).Unix(),
	})

	tokenString, err := token.SignedString([]byte(key))
	if err != nil {
		fmt.Printf("err.Error(): %v\n", err.Error())
		return ""
	}
	fmt.Printf("tokenString: %v\n", tokenString)
	return tokenString
}

func ValidateTokenAndGetUserID(tokenReq string) string {
	token, err := jwt.Parse(tokenReq, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(key), nil
	})
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		fmt.Println(claims["userid"])
		return claims["userid"].(string)
	} else {
		fmt.Println(err)
	}
	return ""
}
