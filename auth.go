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

// Credenziali statiche per accesso limitato (solo gestione utenti/rappresentanti)
const (
	STATIC_USERNAME = "AggiuntaRappresentanti"
	STATIC_PASSWORD = "CONNrap1"
	ROLE_FULL       = "full"
	ROLE_USERS_ONLY = "users_only"
)

func hashPassword(password string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(password)))
}

func generateToken(moderatorId int, role string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"moderator_id": moderatorId,
		"role":         role,
		"exp":          time.Now().Add(24 * time.Hour).Unix(),
		"iat":          time.Now().Unix(),
		"type":         "moderator_session",
	})
	return token.SignedString([]byte(CONF.JWT_SECRET))
}

func validateToken(tokenString string) (int, string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte(CONF.JWT_SECRET), nil
	})

	if err != nil {
		return 0, "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		tokenType, _ := claims["type"].(string)
		if tokenType != "moderator_session" {
			return 0, "", fmt.Errorf("invalid token type")
		}

		moderatorId, ok := claims["moderator_id"].(float64)
		if !ok {
			return 0, "", fmt.Errorf("invalid moderator_id")
		}

		role, _ := claims["role"].(string)
		if role == "" {
			role = ROLE_FULL // Default per retrocompatibilità
		}

		return int(moderatorId), role, nil
	}

	return 0, "", fmt.Errorf("invalid token")
}

func authMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if ctx.Request.Method == "OPTIONS" {
			ctx.AbortWithStatus(http.StatusNoContent)
			return
		}
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

		moderatorId, role, err := validateToken(parts[1])
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
		ctx.Set("role", role)
		ctx.Next()
	}
}

// Middleware per verificare accesso completo (blocca users_only)
func requireFullAccess() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		role, _ := ctx.Get("role")
		if role == ROLE_USERS_ONLY {
			ctx.JSON(http.StatusForbidden, ErrorResponse{
				Status: "error",
				Error:  "access_denied",
				Msg:    "Non hai i permessi per questa sezione",
			})
			ctx.Abort()
			return
		}
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

	// Check static credentials first (users_only role)
	if req.Username == STATIC_USERNAME && req.Password == STATIC_PASSWORD {
		token, err := generateToken(-1, ROLE_USERS_ONLY)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, ErrorResponse{
				Status: "error",
				Error:  "token_error",
				Msg:    "Error generating token",
			})
			return
		}

		ctx.JSON(http.StatusOK, LoginResponse{
			Status: "ok",
			Token:  token,
			Moderator: &Moderator{
				ID:       -1,
				Username: STATIC_USERNAME,
				Name:     "Gestione Rappresentanti",
			},
			Role: ROLE_USERS_ONLY,
		})
		return
	}

	// Check database credentials (full access)
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

	token, err := generateToken(moderator.ID, ROLE_FULL)
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
		Role:      ROLE_FULL,
	})
}

func handleVerify(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, SuccessResponse{
		Status: "ok",
		Msg:    "Token is valid",
	})
}
