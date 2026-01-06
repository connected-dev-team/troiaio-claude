package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func handleGetSchools(ctx *gin.Context) {
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

	var rows interface{ Next() bool }
	cityIdStr := ctx.Query("city_id")
	if cityIdStr != "" {
		cityId, err := strconv.Atoi(cityIdStr)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, ErrorResponse{
				Status: "error",
				Error:  "invalid_city_id",
				Msg:    "Invalid city ID",
			})
			return
		}
		r, err := QuerySchoolsByCity(db, cityId)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, ErrorResponse{
				Status: "error",
				Error:  "query_error",
				Msg:    "Error querying schools",
			})
			return
		}
		rows = r
		defer r.Close()
	} else {
		r, err := QueryAllSchools(db)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, ErrorResponse{
				Status: "error",
				Error:  "query_error",
				Msg:    "Error querying schools",
			})
			return
		}
		rows = r
		defer r.Close()
	}

	var schools []School
	for rows.Next() {
		var school School
		var emailDomain *string
		if scanner, ok := rows.(interface {
			Scan(dest ...any) error
		}); ok {
			if err := scanner.Scan(&school.ID, &school.Name, &emailDomain, &school.CityID, &school.CityName); err != nil {
				continue
			}
			if emailDomain != nil {
				school.EmailDomain = *emailDomain
			}
			schools = append(schools, school)
		}
	}

	ctx.JSON(http.StatusOK, DataResponse{
		Status: "ok",
		Data:   schools,
	})
}

func handleAddSchool(ctx *gin.Context) {
	var req SchoolRequest
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
			Msg:    "School name is required",
		})
		return
	}

	if req.CityID == 0 {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Status: "error",
			Error:  "missing_city_id",
			Msg:    "City ID is required",
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

	if err := InsertSchool(db, req.Name, req.CityID, req.EmailDomain); err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{
			Status: "error",
			Error:  "insert_error",
			Msg:    "Error inserting school: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, SuccessResponse{
		Status: "ok",
		Msg:    "School added successfully",
	})
}

func handleUpdateSchool(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Status: "error",
			Error:  "invalid_id",
			Msg:    "Invalid school ID",
		})
		return
	}

	var req SchoolUpdateRequest
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

	if err := UpdateSchool(db, id, req.Name, req.EmailDomain); err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{
			Status: "error",
			Error:  "update_error",
			Msg:    "Error updating school",
		})
		return
	}

	ctx.JSON(http.StatusOK, SuccessResponse{
		Status: "ok",
		Msg:    "School updated successfully",
	})
}

func handleDeleteSchool(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Status: "error",
			Error:  "invalid_id",
			Msg:    "Invalid school ID",
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

	if err := DeleteSchool(db, id); err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{
			Status: "error",
			Error:  "delete_error",
			Msg:    "Error deleting school. Make sure no users reference this school.",
		})
		return
	}

	ctx.JSON(http.StatusOK, SuccessResponse{
		Status: "ok",
		Msg:    "School deleted successfully",
	})
}
