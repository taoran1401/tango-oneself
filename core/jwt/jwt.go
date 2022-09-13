package jwt

import (
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"taogin/app/types"
	"time"
	//"taogin/config"
	//"data-center-go/global"
	//"data-center-go/model/system/request"
)

type CustomClaims struct {
	types.BaseClaims
	jwt.RegisteredClaims
}

type JWT struct {
	SigningKey []byte
	ExpiresAt  int64
	Issuer     string
}

var (
	TokenExpired     = errors.New("Token is expired")
	TokenNotValidYet = errors.New("Token not active yet")
	TokenMalformed   = errors.New("That's not even a token")
	TokenInvalid     = errors.New("Couldn't handle this token:")
)

func NewJwt(SigningKey string, ExpiresAt int64, Issuer string) *JWT {
	return &JWT{
		SigningKey: []byte(SigningKey),
		ExpiresAt:  ExpiresAt,
		Issuer:     Issuer,
	}
}

//生成claims
func (this *JWT) CreateClaims(baseClaims types.BaseClaims) (claims CustomClaims) {
	claims = CustomClaims{
		BaseClaims: baseClaims,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),                                                  //签发时间
			NotBefore: jwt.NewNumericDate(time.Now()),                                                  //生效时间
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(this.ExpiresAt) * time.Second)), //有效时间
			Issuer:    this.Issuer,                                                                     //签名发行者
		},
	}
	return claims
}

// 创建一个token
func (this *JWT) CreateToken(claims CustomClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(this.SigningKey)
}

// 解析 token
func (j *JWT) ParseToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (i interface{}, e error) {
		return j.SigningKey, nil
	})
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, TokenMalformed
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				// Token is expired
				return nil, TokenExpired
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, TokenNotValidYet
			} else {
				return nil, TokenInvalid
			}
		}
	}

	if token != nil {
		if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
			return claims, nil
		}
		return nil, TokenInvalid

	} else {
		return nil, TokenInvalid
	}
}
