package jwt

import (
	jsonWebToken "github.com/golang-jwt/jwt/v5"
)

type ECBearTokenClaims struct {
	Uuid string `json:"uuid"` // EC user UUID
	Scp  string `json:"scp"`
	jsonWebToken.RegisteredClaims
}

// SignAccessToken - 根據傳入的 jwtSecret, 簽出具有 subject: memberID, Id: jti, uuid: uuid 的 jwtToken
func SignAccessToken(jwtSecret string, memberID string, uuid string, jti string) (string, error) {
	claims := &ECBearTokenClaims{
		Uuid: uuid,
		Scp:  "member",
		RegisteredClaims: jsonWebToken.RegisteredClaims{
			ID:      jti,
			Subject: memberID,
		},
	}
	token := jsonWebToken.NewWithClaims(jsonWebToken.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtSecret))
}
