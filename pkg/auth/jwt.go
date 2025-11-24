package auth

import (
    "errors"
    "strings"
    "time"

    "github.com/golang-jwt/jwt/v5"
)

type JWTManager struct {
    secret     []byte
    expiration time.Duration
}

func NewJWTManager(secret []byte, expiration time.Duration) *JWTManager {
    return &JWTManager{
        secret:     secret,
        expiration: expiration,
    }
}

func (j *JWTManager) GenerateToken(userID int, email, role, tenantID string) (string, error) {
    claims := jwt.MapClaims{
        "user_id":   userID,
        "email":     email,
        "role":      role,
        "tenant_id": tenantID,
        "exp":       time.Now().Add(j.expiration).Unix(),
        "iat":       time.Now().Unix(),
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(j.secret)
}

func (j *JWTManager) ValidateToken(tokenString string) (*jwt.MapClaims, error) {
    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, errors.New("unexpected signing method")
        }
        return j.secret, nil
    })

    if err != nil {
        return nil, err
    }

    if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
        return &claims, nil
    }

    return nil, errors.New("invalid token")
}

func (j *JWTManager) ExtractTokenFromHeader(authHeader string) (string, error) {
    if authHeader == "" {
        return "", errors.New("authorization header missing")
    }

    parts := strings.Split(authHeader, " ")
    if len(parts) != 2 || parts[0] != "Bearer" {
        return "", errors.New("invalid authorization header format")
    }

    return parts[1], nil
}

func (j *JWTManager) HasPermission(userRole string, requiredRole string) bool {
    roleHierarchy := map[string]int{
        "user":     1,
        "agent":    2,
        "supervisor": 3,
        "admin":    4,
    }

    userLevel, userExists := roleHierarchy[userRole]
    requiredLevel, requiredExists := roleHierarchy[requiredRole]

    if !userExists || !requiredExists {
        return false
    }

    return userLevel >= requiredLevel
}
