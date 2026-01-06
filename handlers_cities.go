package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func handleGetCities(ctx *gin.Context) {
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

	rows, err := QueryAllCities(db)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{
			Status: "error",
			Error:  "query_error",
			Msg:    "Error querying cities",
		})
		return
	}
	defer rows.Close()

	var cities []City
	for rows.Next() {
		var city City
		if err := rows.Scan(&city.ID, &city.Name, &city.Region); err != nil {
			continue
		}
		cities = append(cities, city)
	}

	ctx.JSON(http.StatusOK, DataResponse{
		Status: "ok",
		Data:   cities,
	})
}

func handleAddCity(ctx *gin.Context) {
	var req CityRequest
	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Status: "error",
			Error:  "malformed_json",
			Msg:    "Invalid request body",
		})
		return
	}

	if req.Name == "" {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Status: "error",
			Error:  "missing_name",
			Msg:    "City name is required",
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

	if err := InsertCity(db, req.Name, req.Region); err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{
			Status: "error",
			Error:  "insert_error",
			Msg:    "Error inserting city: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, SuccessResponse{
		Status: "ok",
		Msg:    "City added successfully",
	})
}

func handleUpdateCity(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Status: "error",
			Error:  "invalid_id",
			Msg:    "Invalid city ID",
		})
		return
	}

	var req CityRequest
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

	if err := UpdateCity(db, id, req.Name, req.Region); err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{
			Status: "error",
			Error:  "update_error",
			Msg:    "Error updating city",
		})
		return
	}

	ctx.JSON(http.StatusOK, SuccessResponse{
		Status: "ok",
		Msg:    "City updated successfully",
	})
}

func handleDeleteCity(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Status: "error",
			Error:  "invalid_id",
			Msg:    "Invalid city ID",
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

	if err := DeleteCity(db, id); err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{
			Status: "error",
			Error:  "delete_error",
			Msg:    "Error deleting city. Make sure no schools reference this city.",
		})
		return
	}

	ctx.JSON(http.StatusOK, SuccessResponse{
		Status: "ok",
		Msg:    "City deleted successfully",
	})
}
