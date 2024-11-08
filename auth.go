package main

import (
    "time"
    "github.com/gin-gonic/gin"
    "net/http"
    "strings"
    "github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte("secret_key")

type Claims struct {
    Username string `json:"username"`
    jwt.StandardClaims
}
func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        if !strings.HasPrefix(authHeader, "Bearer ") {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing or invalid token"})
            return
        }

        tokenStr := authHeader[len("Bearer "):]
        claims := &Claims{}
        token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
            return jwtKey, nil
        })

        if err != nil || !token.Valid {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
            return
        }

        c.Set("username", claims.Username)
        c.Next()
    }
}

// jwt token generate
func loginHandler(c *gin.Context) {
    var loginDetails struct {
        Username string `json:"username"`
        Password string `json:"password"`
    }
    if err := c.BindJSON(&loginDetails); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
        return
    }

    if loginDetails.Username != "user" || loginDetails.Password != "password" {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
        return
    }

    expirationTime := time.Now().Add(24 * time.Hour)
    claims := &Claims{
        Username: loginDetails.Username,
        StandardClaims: jwt.StandardClaims{
            ExpiresAt: expirationTime.Unix(),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    tokenString, err := token.SignedString(jwtKey)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create token"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"token": tokenString})
}
