package main

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func handleSearchUsers(ctx *gin.Context) {
	searchTerm := strings.TrimSpace(ctx.Query("q"))

	if searchTerm == "" {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Status: "error",
			Error:  "missing_search_term",
			Msg:    "Search term is required (use ?q=...)",
		})
		return
	}

	if len(searchTerm) < 2 {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Status: "error",
			Error:  "search_term_too_short",
			Msg:    "Search term must be at least 2 characters",
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

	rows, err := SearchUsers(db, searchTerm)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{
			Status: "error",
			Error:  "query_error",
			Msg:    err.Error(),
		})
		return
	}
	defer rows.Close()

	users := []User{}

	for rows.Next() {
		var user User
		if err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.PersonalEmail,
			&user.FirstName,
			&user.LastName,
			&user.Role,
			&user.SchoolName,
			&user.CityName,
		); err == nil {
			users = append(users, user)
		}
	}

	ctx.JSON(http.StatusOK, DataResponse{
		Status: "ok",
		Data:   users,
	})
}

func handleSetUserRole(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Status: "error",
			Error:  "invalid_id",
			Msg:    "Invalid user ID",
		})
		return
	}

	var req SetRoleRequest
	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Status: "error",
			Error:  "malformed_json",
			Msg:    "Invalid request body",
		})
		return
	}

	// Validate role - only allow user and representative for now
	validRoles := map[string]bool{"user": true, "representative": true}
	if !validRoles[req.Role] {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Status: "error",
			Error:  "invalid_role",
			Msg:    "Role must be: user or representative",
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

	if err := SetUserRole(db, id, req.Role); err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{
			Status: "error",
			Error:  "update_error",
			Msg:    "Error updating user role",
		})
		return
	}

	ctx.JSON(http.StatusOK, SuccessResponse{
		Status: "ok",
		Msg:    "User role updated successfully",
	})
}

func handleGetUser(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Status: "error",
			Error:  "invalid_id",
			Msg:    "Invalid user ID",
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

	rows, err := GetUserById(db, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{
			Status: "error",
			Error:  "query_error",
			Msg:    "Error getting user",
		})
		return
	}
	defer rows.Close()

	if !rows.Next() {
		ctx.JSON(http.StatusNotFound, ErrorResponse{
			Status: "error",
			Error:  "user_not_found",
			Msg:    "User not found",
		})
		return
	}

	var user User
	if err := rows.Scan(
		&user.ID, &user.Email, &user.PersonalEmail,
		&user.FirstName, &user.LastName, &user.Role,
		&user.SchoolName, &user.CityName,
	); err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{
			Status: "error",
			Error:  "scan_error",
			Msg:    "Error reading user data",
		})
		return
	}

	ctx.JSON(http.StatusOK, DataResponse{
		Status: "ok",
		Data:   user,
	})
}
