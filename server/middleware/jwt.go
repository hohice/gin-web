package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/hohice/gin-web/server/ex"

	"github.com/gin-gonic/gin"
	"github.com/hohice/gin-web/pkg/setting"

	jwt "github.com/dgrijalva/jwt-go"
)

var (
	enable      bool
	tokenLookup string
	authScheme  string
	jwtSecret   []byte
)

func init() {
	registerSelf(func(conf *setting.Configs) (error, Closeble) {
		enable = conf.Auth.Enable
		jwtSecret = []byte(conf.Auth.JwtSecret)
		tokenLookup = conf.Auth.TokenLookup
		authScheme = conf.Auth.AuthScheme
		return nil, func() {}
	})
}

func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		var code int

		code = ex.SUCCESS
		token := getExtractor()(c)
		if token == "" {
			code = ex.INVALID_PARAMS
		} else {
			claims, err := ParseToken(token)
			if err != nil {
				code = ex.ERROR_AUTH_CHECK_TOKEN_FAIL
			} else if time.Now().Unix() > claims.ExpiresAt {
				code = ex.ERROR_AUTH_CHECK_TOKEN_TIMEOUT
			}
		}

		if code != ex.SUCCESS {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": code,
				"msg":  ex.GetMsg(code),
			})

			c.Abort()
			return
		}

		c.Next()
	}
}

type Claims struct {
	Username string `json:"username"`
	Password string `json:"password"`
	jwt.StandardClaims
}

func GenerateToken(username, password string) (string, error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(3 * time.Hour)

	claims := Claims{
		username,
		password,
		jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    "ginS",
		},
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString(jwtSecret)

	return token, err
}

func ParseToken(token string) (*Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}

	return nil, err
}

type jwtExtractor func(c *gin.Context) string

func getExtractor() jwtExtractor {
	if len(tokenLookup) == 0 {
		tokenLookup = "header:" + "token"
	}
	if len(authScheme) == 0 {
		authScheme = "Bearer"
	}
	parts := strings.Split(tokenLookup, ":")
	extractor := jwtFromHeader(parts[1], authScheme)
	switch parts[0] {
	case "query":
		extractor = jwtFromQuery(parts[1])
	case "cookie":
		extractor = jwtFromCookie(parts[1])
	}

	return extractor
}

// jwtFromHeader returns a `jwtExtractor` that extracts token from the request header.
func jwtFromHeader(header string, authScheme string) jwtExtractor {
	return func(c *gin.Context) string {
		auth := c.GetHeader(header)
		l := len(authScheme)
		if len(auth) > l+1 && auth[:l] == authScheme {
			return auth[l+1:]
		}
		return ""
	}
}

// jwtFromQuery returns a `jwtExtractor` that extracts token from the query string.
func jwtFromQuery(param string) jwtExtractor {
	return func(c *gin.Context) string {
		token := c.Query(param)
		if token == "" {
			return ""
		}
		return token
	}
}

func jwtFromCookie(name string) jwtExtractor {
	return func(c *gin.Context) string {
		cookie, err := c.Cookie(name)
		if err != nil {
			return ""
		}
		return cookie
	}
}
