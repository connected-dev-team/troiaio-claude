package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func handleGetStatistics(ctx *gin.Context) {
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

	var stats FullStatistics

	// 1. Get total stats
	rows, err := QueryTotalStats(db)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{
			Status: "error",
			Error:  "query_error",
			Msg:    "Error querying total stats",
		})
		return
	}
	if rows.Next() {
		rows.Scan(
			&stats.Totals.TotalUsers,
			&stats.Totals.TotalPosts,
			&stats.Totals.TotalSpotted,
			&stats.Totals.ApprovedPosts,
			&stats.Totals.ApprovedSpotted,
			&stats.Totals.TotalPostLikes,
			&stats.Totals.TotalSpottedLikes,
			&stats.Totals.TotalCities,
			&stats.Totals.TotalSchools,
		)
		stats.Totals.TotalInteractions = stats.Totals.TotalPostLikes + stats.Totals.TotalSpottedLikes
	}
	rows.Close()

	// 2. Get stats by city
	rows, err = QueryStatsByCity(db)
	if err == nil {
		for rows.Next() {
			var city CityStats
			rows.Scan(&city.ID, &city.Name, &city.Region, &city.UserCount, &city.SchoolCount, &city.PostCount, &city.SpottedCount)
			stats.CitiesStats = append(stats.CitiesStats, city)
		}
		rows.Close()
	}

	// 3. Get stats by school
	rows, err = QueryStatsBySchool(db)
	if err == nil {
		for rows.Next() {
			var school SchoolStats
			rows.Scan(&school.ID, &school.Name, &school.CityID, &school.CityName, &school.UserCount, &school.PostCount, &school.SpottedCount)
			stats.SchoolsStats = append(stats.SchoolsStats, school)
		}
		rows.Close()
	}

	// 4. Get users over time
	rows, err = QueryUsersOverTime(db)
	if err == nil {
		for rows.Next() {
			var ts TimeStats
			var month time.Time
			rows.Scan(&month, &ts.Count)
			ts.Month = month.Format("2006-01")
			stats.UsersOverTime = append(stats.UsersOverTime, ts)
		}
		rows.Close()
	}

	// 5. Get posts over time
	rows, err = QueryPostsOverTime(db)
	if err == nil {
		for rows.Next() {
			var ts TimeStats
			var month time.Time
			rows.Scan(&month, &ts.Count)
			ts.Month = month.Format("2006-01")
			stats.PostsOverTime = append(stats.PostsOverTime, ts)
		}
		rows.Close()
	}

	// 6. Get spotted over time
	rows, err = QuerySpottedOverTime(db)
	if err == nil {
		for rows.Next() {
			var ts TimeStats
			var month time.Time
			rows.Scan(&month, &ts.Count)
			ts.Month = month.Format("2006-01")
			stats.SpottedOverTime = append(stats.SpottedOverTime, ts)
		}
		rows.Close()
	}

	// 7. Get top cities
	rows, err = QueryTopCities(db, 10)
	if err == nil {
		for rows.Next() {
			var city TopCity
			rows.Scan(&city.ID, &city.Name, &city.Region, &city.UserCount)
			stats.TopCities = append(stats.TopCities, city)
		}
		rows.Close()
	}

	// 8. Get top schools
	rows, err = QueryTopSchools(db, 10)
	if err == nil {
		for rows.Next() {
			var school TopSchool
			rows.Scan(&school.ID, &school.Name, &school.CityName, &school.UserCount)
			stats.TopSchools = append(stats.TopSchools, school)
		}
		rows.Close()
	}

	// Initialize empty slices if nil
	if stats.CitiesStats == nil {
		stats.CitiesStats = []CityStats{}
	}
	if stats.SchoolsStats == nil {
		stats.SchoolsStats = []SchoolStats{}
	}
	if stats.UsersOverTime == nil {
		stats.UsersOverTime = []TimeStats{}
	}
	if stats.PostsOverTime == nil {
		stats.PostsOverTime = []TimeStats{}
	}
	if stats.SpottedOverTime == nil {
		stats.SpottedOverTime = []TimeStats{}
	}
	if stats.TopCities == nil {
		stats.TopCities = []TopCity{}
	}
	if stats.TopSchools == nil {
		stats.TopSchools = []TopSchool{}
	}

	ctx.JSON(http.StatusOK, DataResponse{
		Status: "ok",
		Data:   stats,
	})
}
