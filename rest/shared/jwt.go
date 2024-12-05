package shared

import (
	"crypto/rand"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"time"
)

var JWTManagerInstance *JWTManager = mustNewJWTManager()

const TokenExpirationDuration = time.Hour * 4 // Token expires in x hours

// JWTManager handles signing and verifying JWT tokens
type JWTManager struct {
	secretKey []byte
}

// NewJWTManager creates a new JWT manager with a randomly generated secret key
func newJWTManager() (*JWTManager, error) {
	secretKey, err := generateRandomKey(32) // Generate a 32-byte random key
	if err != nil {
		return nil, err
	}
	return &JWTManager{secretKey: secretKey}, nil
}

func mustNewJWTManager() *JWTManager {
	jwtManager, err := newJWTManager()
	if err != nil {
		panic(err)
	}
	return jwtManager
}

// generateRandomKey generates a secure random key of the specified length
func generateRandomKey(length int) ([]byte, error) {
	key := make([]byte, length)
	if _, err := rand.Read(key); err != nil {
		return nil, err
	}
	return key, nil
}

// CreateToken generates a JWT for the given generic data
// T represents a generic type
func (jwtManager *JWTManager) CreateToken(data *User) (string, error) {
	claims := jwt.MapClaims{
		"sub": data.Name,                                      // Subject (entity that is being issued the token)
		"aud": []string{"magic-wan-rest"},                     // Audience (nodes that are allowed to work based off of this token)
		"exp": time.Now().Add(TokenExpirationDuration).Unix(), // expiration time

		"iss": "magic-wan",       // Issuer
		"iat": time.Now().Unix(), // Issued at
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtManager.secretKey)
}

// ParseToken parses a JWT string into a generic type
func (jwtManager *JWTManager) ParseToken(tokenString string) (*User, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Ensure that the token's signing method conforms to expected method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtManager.secretKey, nil
	})
	if err != nil {
		return nil, err
	}

	if token.Valid {
		subj, err := token.Claims.GetSubject()
		if err != nil {
			return nil, err
		}
		return &User{Name: subj}, nil
	} else {
		return nil, fmt.Errorf("invalid token")
	}
}

func (jwtManager *JWTManager) ParseFromRequest(r *http.Request) *User {
	// Get the Authorization token from the cookie
	authCookie, err := r.Cookie("Authorization")
	if err != nil {
		// No Cookie -> Not Authorized
		return nil
	}

	// Extract the token from the cookie
	token := authCookie.Value

	// (try to) Parse the token
	user, err := JWTManagerInstance.ParseToken(token)
	if err != nil {
		// Failed to parse -> Not Authorized
		return nil
	}
	return user
}
