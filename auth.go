package main

import (
	"crypto/md5"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func hashPassword(password string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(password)))
}

func generateToken(moderatorId int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"moderator_id": moderatorId,
		"exp":          time.Now().Add(24 * time.Hour).Unix(),
		"iat":          time.Now().Unix(),
		"type":         "moderator_session",
	})
	return token.SignedString([]byte(CONF.JWT_SECRET))
}

func validateToken(tokenString string) (int, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte(CONF.JWT_SECRET), nil
	})

	if err != nil {
		return 0, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		tokenType, _ := claims["type"].(string)
		if tokenType != "moderator_session" {
			return 0, fmt.Errorf("invalid token type")
		}

		moderatorId, ok := claims["moderator_id"].(float64)
		if !ok {
			return 0, fmt.Errorf("invalid moderator_id")
		}
		return int(moderatorId), nil
	}

	return 0, fmt.Errorf("invalid token")
}

func authMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.JSON(http.StatusUnauthorized, ErrorResponse{
				Status: "error",
				Error:  "missing_auth_header",
				Msg:    "Authorization header required",
			})
			ctx.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			ctx.JSON(http.StatusUnauthorized, ErrorResponse{
				Status: "error",
				Error:  "invalid_auth_header",
				Msg:    "Invalid Authorization header format",
			})
			ctx.Abort()
			return
		}

		moderatorId, err := validateToken(parts[1])
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, ErrorResponse{
				Status: "error",
				Error:  "invalid_token",
				Msg:    "Token is invalid or expired",
			})
			ctx.Abort()
			return
		}

		ctx.Set("moderator_id", moderatorId)
		ctx.Next()
	}
}

func handleLogin(ctx *gin.Context) {
	var req LoginRequest
	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Status: "error",
			Error:  "malformed_json",
			Msg:    "Invalid request body",
		})
		return
	}

	db, err := makeDbaseConnection()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{
			Status: "error",
			Error:  "database_unreachable",
			Msg:    "Cannot connect to database",
		})
		return
	}
	defer db.Close(ctx)

	passwdHash := hashPassword(req.Password)
	rows, err := QueryModeratorByCredentials(db, req.Username, passwdHash)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{
			Status: "error",
			Error:  "query_error",
			Msg:    "Error querying database",
		})
		return
	}
	defer rows.Close()

	if !rows.Next() {
		ctx.JSON(http.StatusUnauthorized, ErrorResponse{
			Status: "error",
			Error:  "invalid_credentials",
			Msg:    "Invalid username or password",
		})
		return
	}

	var moderator Moderator
	if err := rows.Scan(&moderator.ID, &moderator.Username, &moderator.Name); err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{
			Status: "error",
			Error:  "scan_error",
			Msg:    "Error reading moderator data",
		})
		return
	}

	token, err := generateToken(moderator.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{
			Status: "error",
			Error:  "token_error",
			Msg:    "Error generating token",
		})
		return
	}

	ctx.JSON(http.StatusOK, LoginResponse{
		Status:    "ok",
		Token:     token,
		Moderator: &moderator,
	})
}

func handleVerify(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, SuccessResponse{
		Status: "ok",
		Msg:    "Token is valid",
	})
}
