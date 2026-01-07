package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func handleGetPendingPosts(ctx *gin.Context) {
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

	rows, err := QueryPendingPosts(db)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{
			Status: "error",
			Error:  "query_error",
			Msg:    "Error querying pending posts",
		})
		return
	}
	defer rows.Close()

	var posts []PendingPost
	for rows.Next() {
		var post PendingPost
		if err := rows.Scan(
			&post.ID, &post.Content, &post.CreatorID, &post.CreationTimestamp,
			&post.LikesCount, &post.CreatorFirstName, &post.CreatorLastName,
			&post.CreatorEmail, &post.SchoolName, &post.CityName,
		); err != nil {
			continue
		}
		posts = append(posts, post)
	}

	ctx.JSON(http.StatusOK, DataResponse{
		Status: "ok",
		Data:   posts,
	})
}

func handleGetReportedPosts(ctx *gin.Context) {
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

	rows, err := QueryReportedPosts(db)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{
			Status: "error",
			Error:  "query_error",
			Msg:    "Error querying reported posts",
		})
		return
	}
	defer rows.Close()

	var posts []ReportedPost
	for rows.Next() {
		var post ReportedPost
		if err := rows.Scan(
			&post.ID, &post.Content, &post.CreatorID, &post.CreationTimestamp,
			&post.LikesCount, &post.CreatorFirstName, &post.CreatorLastName,
			&post.CreatorEmail, &post.SchoolName, &post.CityName, &post.ReportCount,
		); err != nil {
			continue
		}
		posts = append(posts, post)
	}

	ctx.JSON(http.StatusOK, DataResponse{
		Status: "ok",
		Data:   posts,
	})
}

func handleApprovePost(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Status: "error",
			Error:  "invalid_id",
			Msg:    "Invalid post ID",
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

	if err := ApprovePost(db, id); err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{
			Status: "error",
			Error:  "update_error",
			Msg:    "Error approving post",
		})
		return
	}

	ctx.JSON(http.StatusOK, SuccessResponse{
		Status: "ok",
		Msg:    "Post approved successfully",
	})
}

func handleRejectPost(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Status: "error",
			Error:  "invalid_id",
			Msg:    "Invalid post ID",
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

	if err := RejectPost(db, id); err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{
			Status: "error",
			Error:  "update_error",
			Msg:    "Error rejecting post",
		})
		return
	}

	ctx.JSON(http.StatusOK, SuccessResponse{
		Status: "ok",
		Msg:    "Post rejected successfully",
	})
}

func handleDeletePost(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Status: "error",
			Error:  "invalid_id",
			Msg:    "Invalid post ID",
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

	if err := DeletePostById(db, id); err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{
			Status: "error",
			Error:  "delete_error",
			Msg:    "Error deleting post",
		})
		return
	}

	ctx.JSON(http.StatusOK, SuccessResponse{
		Status: "ok",
		Msg:    "Post deleted successfully",
	})
}

func handleGetAllPosts(ctx *gin.Context) {
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

	rows, err := QueryAllPosts(db)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{
			Status: "error",
			Error:  "query_error",
			Msg:    "Error querying posts",
		})
		return
	}
	defer rows.Close()

	var posts []AllPost
	for rows.Next() {
		var post AllPost
		if err := rows.Scan(
			&post.ID, &post.Content, &post.CreatorID, &post.CreationTimestamp,
			&post.LikesCount, &post.CreatorFirstName, &post.CreatorLastName,
			&post.CreatorEmail, &post.SchoolName, &post.CityName, &post.Status,
		); err != nil {
			continue
		}
		posts = append(posts, post)
	}

	ctx.JSON(http.StatusOK, DataResponse{
		Status: "ok",
		Data:   posts,
	})
}

func handleSetPostStatus(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Status: "error",
			Error:  "invalid_id",
			Msg:    "Invalid post ID",
		})
		return
	}

	var req SetStatusRequest
	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Status: "error",
			Error:  "malformed_json",
			Msg:    "Invalid request body",
		})
		return
	}

	// Validate status
	validStatuses := map[string]bool{"received": true, "approved": true, "rejected": true}
	if !validStatuses[req.Status] {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Status: "error",
			Error:  "invalid_status",
			Msg:    "Status must be: received, approved, or rejected",
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

	if err := SetPostStatus(db, id, req.Status); err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{
			Status: "error",
			Error:  "update_error",
			Msg:    "Error updating post status",
		})
		return
	}

	ctx.JSON(http.StatusOK, SuccessResponse{
		Status: "ok",
		Msg:    "Post status updated successfully",
	})
}
