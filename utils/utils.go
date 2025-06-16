package utils

import(
	"fmt"
    "time"
	"net/http"
	"encoding/json"
	"gorm.io/gorm"
	"errors"
	"strings"
	"golang.org/x/crypto/bcrypt"
	"github.com/golang-jwt/jwt/v5"
)

const SecretKey = "2b87694969af1987c4542c2349a0948569f22f63c2da0c01dff1e0404edd37e27d6cb706507cc3fb50edba793041686fcf81749d00457253c0928f3b5b8cfcee"
const TokenExpiration = 4
const RefreshTokenExpiration = 2
var secretKey = []byte(SecretKey)

type ValidationError struct {
	Parameter string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("invalid value for parameter: %s", e.Parameter)
}

func WriteJSONResponse(w http.ResponseWriter, response Response, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}

func IsDuplicateEntryError(err error) bool {
	return errors.Is(err, gorm.ErrDuplicatedKey)
}

func GetIDFromPath(path string,) string {
	// Match "/members/{id}" pattern
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) == 2 {
		return parts[1]
	}
	return ""
}

// Hash password before storing in DB
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// Compare hashed password with input password
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// Create Token generates a JWT token
func CreateToken(name string) (string, error) {
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, 
        jwt.MapClaims{ 
        "name": name, 
        "exp": time.Now().Add(time.Hour * TokenExpiration).Unix(), 
        })

    tokenString, err := token.SignedString(secretKey)
    if err != nil {
    return "", err
    }
	return tokenString, nil
}

//Verify Token 
func VerifyToken(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
	   return secretKey, nil
	})
   
	if err != nil {
	   return err
	}
   
	if !token.Valid {
	   return fmt.Errorf("invalid token")
	}
   
	return nil
}

func GenerateRefreshToken(name string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, 
        jwt.MapClaims{ 
        "name": name, 
        "exp": time.Now().Add(time.Hour * 24 * RefreshTokenExpiration).Unix(), 
        })
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
	return "", err
	}
	return tokenString, nil
}