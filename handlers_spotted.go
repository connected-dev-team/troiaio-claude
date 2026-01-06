package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func handleGetPendingSpotted(ctx *gin.Context) {
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

	rows, err := QueryPendingSpotted(db)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{
			Status: "error",
			Error:  "query_error",
			Msg:    "Error querying pending spotted",
		})
		return
	}
	defer rows.Close()

	var spotted []PendingSpotted
	for rows.Next() {
		var s PendingSpotted
		if err := rows.Scan(
			&s.ID, &s.Content, &s.CreatorID, &s.CreationTimestamp,
			&s.LikesCount, &s.Visibility, &s.Color,
			&s.CreatorFirstName, &s.CreatorLastName, &s.CreatorEmail,
			&s.SchoolName, &s.CityName, &s.VisibilityDesc,
		); err != nil {
			continue
		}
		spotted = append(spotted, s)
	}

	ctx.JSON(http.StatusOK, DataResponse{
		Status: "ok",
		Data:   spotted,
	})
}

func handleGetReportedSpotted(ctx *gin.Context) {
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

	rows, err := QueryReportedSpotted(db)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{
			Status: "error",
			Error:  "query_error",
			Msg:    "Error querying reported spotted",
		})
		return
	}
	defer rows.Close()

	var spotted []ReportedSpotted
	for rows.Next() {
		var s ReportedSpotted
		if err := rows.Scan(
			&s.ID, &s.Content, &s.CreatorID, &s.CreationTimestamp,
			&s.LikesCount, &s.Visibility, &s.Color,
			&s.CreatorFirstName, &s.CreatorLastName, &s.CreatorEmail,
			&s.SchoolName, &s.CityName, &s.VisibilityDesc, &s.ReportCount,
		); err != nil {
			continue
		}
		spotted = append(spotted, s)
	}

	ctx.JSON(http.StatusOK, DataResponse{
		Status: "ok",
		Data:   spotted,
	})
}

func handleApproveSpotted(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Status: "error",
			Error:  "invalid_id",
			Msg:    "Invalid spotted ID",
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

	if err := ApproveSpotted(db, id); err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{
			Status: "error",
			Error:  "update_error",
			Msg:    "Error approving spotted",
		})
		return
	}

	ctx.JSON(http.StatusOK, SuccessResponse{
		Status: "ok",
		Msg:    "Spotted approved successfully",
	})
}

func handleRejectSpotted(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Status: "error",
			Error:  "invalid_id",
			Msg:    "Invalid spotted ID",
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

	if err := RejectSpotted(db, id); err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{
			Status: "error",
			Error:  "update_error",
			Msg:    "Error rejecting spotted",
		})
		return
	}

	ctx.JSON(http.StatusOK, SuccessResponse{
		Status: "ok",
		Msg:    "Spotted rejected successfully",
	})
}

func handleDeleteSpotted(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Status: "error",
			Error:  "invalid_id",
			Msg:    "Invalid spotted ID",
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

	if err := DeleteSpottedById(db, id); err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{
			Status: "error",
			Error:  "delete_error",
			Msg:    "Error deleting spotted",
		})
		return
	}

	ctx.JSON(http.StatusOK, SuccessResponse{
		Status: "ok",
		Msg:    "Spotted deleted successfully",
	})
}
