package data

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type Token struct {
	ID           int64  `json:"-"`
	UserID       int64  `json:"-"`
	RefreshToken string `json:"refresh_token"`
	AccessToken  string `json:"access_token,omitempty"`
}

type TokenModel struct {
	DB *sql.DB
}

func GenerateTokens(userID int64) (map[string]string, error) {
	// refreshToken := jwt.New(jwt.SigningMethodHS512)
	// refreshClaims := refreshToken.Claims.(jwt.MapClaims)
	// refreshClaims["user_id"] = userID
	// refreshClaims["exp"] = time.Now().Add(time.Hour * 24 * 30).Unix()

	refreshClaims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24 * 30).Unix(),
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS512, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString([]byte(os.Getenv("REFRESH_SK")))
	if err != nil {
		return nil, err
	}
	tokensMap := make(map[string]string)
	tokensMap["refreshToken"] = refreshTokenString

	accessClaims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 4).Unix(),
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS512, accessClaims)
	accessTokenString, err := accessToken.SignedString([]byte(os.Getenv("ACCESS_SK")))

	if err != nil {
		return nil, err
	}
	tokensMap["accessToken"] = accessTokenString
	return tokensMap, nil
}

func (t TokenModel) SaveToken(userID int64, refreshToken string) {
	//find token, if exists => save
	//token does not exist => create new token => generateTokens(userID64)
	//save refreshTOken
}

func (t TokenModel) RemoveToken(refreshToken string) {

}

func (t TokenModel) FindToken(refreshToken string) {

}

func DecodeRefreshToken(refreshToken string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("incorrect signing method: %v", token.Header["alg"])
		}

		return []byte(os.Getenv("REFRESH_SK")), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token yoooo")
}

func DecodeAccessToken(accessToken string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("incorrect signing method: %v", token.Header["alg"])
		}

		return []byte(os.Getenv("ACCESS_SK")), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token yoooo")
}
