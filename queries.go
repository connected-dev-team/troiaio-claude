package main

import (
	"context"

	"github.com/jackc/pgx/v5"
)

// ==================== MODERATORS ====================

func QueryModeratorByCredentials(db *pgx.Conn, username, passwdHash string) (pgx.Rows, error) {
	return db.Query(
		context.Background(),
		"SELECT id, username, name FROM moderators WHERE username=$1 AND passwd_hash=$2 LIMIT 1",
		username, passwdHash,
	)
}

// ==================== CITIES ====================

func QueryAllCities(db *pgx.Conn) (pgx.Rows, error) {
	return db.Query(
		context.Background(),
		"SELECT id, name, region FROM cities ORDER BY name",
	)
}

func InsertCity(db *pgx.Conn, name, region string) error {
	_, err := db.Exec(
		context.Background(),
		"INSERT INTO cities (name, region) VALUES ($1, $2)",
		name, region,
	)
	return err
}

func UpdateCity(db *pgx.Conn, id int, name, region string) error {
	_, err := db.Exec(
		context.Background(),
		"UPDATE cities SET name=$1, region=$2 WHERE id=$3",
		name, region, id,
	)
	return err
}

func DeleteCity(db *pgx.Conn, id int) error {
	_, err := db.Exec(
		context.Background(),
		"DELETE FROM cities WHERE id=$1",
		id,
	)
	return err
}

// ==================== SCHOOLS ====================

func QueryAllSchools(db *pgx.Conn) (pgx.Rows, error) {
	return db.Query(
		context.Background(),
		`SELECT s.id, s.name, s.email_domain, s.city, c.name as city_name
		 FROM schools s
		 JOIN cities c ON s.city = c.id
		 ORDER BY c.name, s.name`,
	)
}

func QuerySchoolsByCity(db *pgx.Conn, cityId int) (pgx.Rows, error) {
	return db.Query(
		context.Background(),
		`SELECT s.id, s.name, s.email_domain, s.city, c.name as city_name
		 FROM schools s
		 JOIN cities c ON s.city = c.id
		 WHERE s.city = $1
		 ORDER BY s.name`,
		cityId,
	)
}

func InsertSchool(db *pgx.Conn, name string, cityId int, emailDomain string) error {
	_, err := db.Exec(
		context.Background(),
		"INSERT INTO schools (name, city, email_domain) VALUES ($1, $2, $3)",
		name, cityId, emailDomain,
	)
	return err
}

func UpdateSchool(db *pgx.Conn, id int, name, emailDomain string) error {
	_, err := db.Exec(
		context.Background(),
		"UPDATE schools SET name=$1, email_domain=$2 WHERE id=$3",
		name, emailDomain, id,
	)
	return err
}

func DeleteSchool(db *pgx.Conn, id int) error {
	_, err := db.Exec(
		context.Background(),
		"DELETE FROM schools WHERE id=$1",
		id,
	)
	return err
}

// ==================== POSTS ====================

func QueryPendingPosts(db *pgx.Conn) (pgx.Rows, error) {
	return db.Query(
		context.Background(),
		`SELECT p.id, p.content, p.creator, p.creation_timestamp, p.likes_count,
		        COALESCE(u.first_name, 'N/A'), COALESCE(u.last_name, ''), COALESCE(u.email, 'N/A'),
		        s.name as school_name, c.name as city_name
		 FROM post p
		 JOIN users u ON p.creator = u.id
		 LEFT JOIN schools s ON u.school = s.id
		 LEFT JOIN cities c ON s.city = c.id
		 WHERE p.status = (SELECT id FROM submit_status WHERE description='received')
		 ORDER BY p.creation_timestamp DESC`,
	)
}

func QueryReportedPosts(db *pgx.Conn) (pgx.Rows, error) {
	return db.Query(
		context.Background(),
		`SELECT p.id, p.content, p.creator, p.creation_timestamp, p.likes_count,
		        COALESCE(u.first_name, 'N/A'), COALESCE(u.last_name, ''), COALESCE(u.email, 'N/A'),
		        s.name as school_name, c.name as city_name,
		        COUNT(rp.id) as report_count
		 FROM post p
		 JOIN users u ON p.creator = u.id
		 LEFT JOIN schools s ON u.school = s.id
		 LEFT JOIN cities c ON s.city = c.id
		 JOIN reported_post rp ON p.id = rp.post_id
		 GROUP BY p.id, p.content, p.creator, p.creation_timestamp, p.likes_count,
		          u.first_name, u.last_name, u.email, s.name, c.name
		 ORDER BY report_count DESC, p.creation_timestamp DESC`,
	)
}

func ApprovePost(db *pgx.Conn, postId int) error {
	_, err := db.Exec(
		context.Background(),
		`UPDATE post SET status = (SELECT id FROM submit_status WHERE description='approved'),
		                 approval_timestamp = NOW()
		 WHERE id = $1`,
		postId,
	)
	return err
}

func RejectPost(db *pgx.Conn, postId int) error {
	_, err := db.Exec(
		context.Background(),
		`UPDATE post SET status = (SELECT id FROM submit_status WHERE description='rejected')
		 WHERE id = $1`,
		postId,
	)
	return err
}

func DeletePostById(db *pgx.Conn, postId int) error {
	// First delete likes
	db.Exec(context.Background(), "DELETE FROM post_like WHERE post_id = $1", postId)
	// Then delete reports
	db.Exec(context.Background(), "DELETE FROM reported_post WHERE post_id = $1", postId)
	// Finally delete the post
	_, err := db.Exec(context.Background(), "DELETE FROM post WHERE id = $1", postId)
	return err
}

func QueryAllPosts(db *pgx.Conn) (pgx.Rows, error) {
	return db.Query(
		context.Background(),
		`SELECT p.id, p.content, p.creator, p.creation_timestamp, p.likes_count,
		        COALESCE(u.first_name, 'N/A'), COALESCE(u.last_name, ''), COALESCE(u.email, 'N/A'),
		        s.name as school_name, c.name as city_name,
		        ss.description as status
		 FROM post p
		 JOIN users u ON p.creator = u.id
		 LEFT JOIN schools s ON u.school = s.id
		 LEFT JOIN cities c ON s.city = c.id
		 JOIN submit_status ss ON p.status = ss.id
		 ORDER BY p.creation_timestamp DESC`,
	)
}

func SetPostStatus(db *pgx.Conn, postId int, status string) error {
	_, err := db.Exec(
		context.Background(),
		`UPDATE post SET status = (SELECT id FROM submit_status WHERE description=$1)
		 WHERE id = $2`,
		status, postId,
	)
	return err
}

// ==================== SPOTTED ====================

func QueryPendingSpotted(db *pgx.Conn) (pgx.Rows, error) {
	return db.Query(
		context.Background(),
		`SELECT sp.id, sp.content, sp.creator, sp.creation_timestamp, sp.likes_count,
		        sp.visibility, COALESCE(sp.color, '#6366f1'),
		        COALESCE(u.first_name, 'N/A'), COALESCE(u.last_name, ''), COALESCE(u.email, 'N/A'),
		        s.name as school_name, c.name as city_name,
		        sv.description as visibility_desc
		 FROM spotted sp
		 JOIN users u ON sp.creator = u.id
		 LEFT JOIN schools s ON u.school = s.id
		 LEFT JOIN cities c ON s.city = c.id
		 JOIN spotted_visibility sv ON sp.visibility = sv.id
		 WHERE sp.status = (SELECT id FROM submit_status WHERE description='received')
		 ORDER BY sp.creation_timestamp DESC`,
	)
}

func QueryReportedSpotted(db *pgx.Conn) (pgx.Rows, error) {
	return db.Query(
		context.Background(),
		`SELECT sp.id, sp.content, sp.creator, sp.creation_timestamp, sp.likes_count,
		        sp.visibility, COALESCE(sp.color, '#6366f1'),
		        COALESCE(u.first_name, 'N/A'), COALESCE(u.last_name, ''), COALESCE(u.email, 'N/A'),
		        s.name as school_name, c.name as city_name,
		        sv.description as visibility_desc,
		        COUNT(rs.id) as report_count
		 FROM spotted sp
		 JOIN users u ON sp.creator = u.id
		 LEFT JOIN schools s ON u.school = s.id
		 LEFT JOIN cities c ON s.city = c.id
		 JOIN spotted_visibility sv ON sp.visibility = sv.id
		 JOIN reported_spotted rs ON sp.id = rs.spotted_id
		 GROUP BY sp.id, sp.content, sp.creator, sp.creation_timestamp, sp.likes_count,
		          sp.visibility, sp.color, u.first_name, u.last_name, u.email,
		          s.name, c.name, sv.description
		 ORDER BY report_count DESC, sp.creation_timestamp DESC`,
	)
}

func ApproveSpotted(db *pgx.Conn, spottedId int) error {
	_, err := db.Exec(
		context.Background(),
		`UPDATE spotted SET status = (SELECT id FROM submit_status WHERE description='approved'),
		                    approval_timestamp = NOW()
		 WHERE id = $1`,
		spottedId,
	)
	return err
}

func RejectSpotted(db *pgx.Conn, spottedId int) error {
	_, err := db.Exec(
		context.Background(),
		`UPDATE spotted SET status = (SELECT id FROM submit_status WHERE description='rejected')
		 WHERE id = $1`,
		spottedId,
	)
	return err
}

func DeleteSpottedById(db *pgx.Conn, spottedId int) error {
	// First delete likes
	db.Exec(context.Background(), "DELETE FROM spotted_like WHERE spotted_id = $1", spottedId)
	// Then delete reports
	db.Exec(context.Background(), "DELETE FROM reported_spotted WHERE spotted_id = $1", spottedId)
	// Finally delete the spotted
	_, err := db.Exec(context.Background(), "DELETE FROM spotted WHERE id = $1", spottedId)
	return err
}

func QueryAllSpotted(db *pgx.Conn) (pgx.Rows, error) {
	return db.Query(
		context.Background(),
		`SELECT sp.id, sp.content, sp.creator, sp.creation_timestamp, sp.likes_count,
		        sp.visibility, COALESCE(sp.color, '#6366f1'),
		        COALESCE(u.first_name, 'N/A'), COALESCE(u.last_name, ''), COALESCE(u.email, 'N/A'),
		        s.name as school_name, c.name as city_name,
		        sv.description as visibility_desc,
		        ss.description as status
		 FROM spotted sp
		 JOIN users u ON sp.creator = u.id
		 LEFT JOIN schools s ON u.school = s.id
		 LEFT JOIN cities c ON s.city = c.id
		 JOIN spotted_visibility sv ON sp.visibility = sv.id
		 JOIN submit_status ss ON sp.status = ss.id
		 ORDER BY sp.creation_timestamp DESC`,
	)
}

func SetSpottedStatus(db *pgx.Conn, spottedId int, status string) error {
	_, err := db.Exec(
		context.Background(),
		`UPDATE spotted SET status = (SELECT id FROM submit_status WHERE description=$1)
		 WHERE id = $2`,
		status, spottedId,
	)
	return err
}

// ==================== USERS ====================

func SearchUsers(db *pgx.Conn, searchTerm string) (pgx.Rows, error) {
	searchPattern := "%" + searchTerm + "%"
	return db.Query(
		context.Background(),
		`SELECT u.id, COALESCE(u.email, ''), u.personal_email,
		        COALESCE(u.first_name, ''), COALESCE(u.last_name, ''),
		        COALESCE(ur.description, 'user') as role,
		        s.name as school_name, c.name as city_name
		 FROM users u
		 LEFT JOIN user_role ur ON u.user_role = ur.id
		 LEFT JOIN users_to_school_cities utsc ON u.id = utsc.user_id
		 LEFT JOIN schools s ON utsc.school = s.id
		 LEFT JOIN cities c ON utsc.city = c.id
		 WHERE LOWER(COALESCE(u.first_name, '')) LIKE LOWER($1)
		    OR LOWER(COALESCE(u.last_name, '')) LIKE LOWER($1)
		    OR LOWER(COALESCE(u.email, '')) LIKE LOWER($1)
		    OR LOWER(COALESCE(u.personal_email, '')) LIKE LOWER($1)
		 ORDER BY u.last_name, u.first_name
		 LIMIT 50`,
		searchPattern,
	)
}

func SetUserRole(db *pgx.Conn, userId int, role string) error {
	_, err := db.Exec(
		context.Background(),
		`UPDATE users SET user_role = (SELECT id FROM user_role WHERE description = $1)
		 WHERE id = $2`,
		role, userId,
	)
	return err
}

func GetUserById(db *pgx.Conn, userId int) (pgx.Rows, error) {
	return db.Query(
		context.Background(),
		`SELECT u.id, COALESCE(u.email, ''), u.personal_email,
		        COALESCE(u.first_name, ''), COALESCE(u.last_name, ''),
		        COALESCE(ur.description, 'user') as role,
		        s.name as school_name, c.name as city_name
		 FROM users u
		 LEFT JOIN user_role ur ON u.user_role = ur.id
		 LEFT JOIN users_to_school_cities utsc ON u.id = utsc.user_id
		 LEFT JOIN schools s ON utsc.school = s.id
		 LEFT JOIN cities c ON utsc.city = c.id
		 WHERE u.id = $1`,
		userId,
	)
}

// ==================== STATISTICS ====================

// Totali generali
func QueryTotalStats(db *pgx.Conn) (pgx.Rows, error) {
	return db.Query(
		context.Background(),
		`SELECT
			(SELECT COUNT(*) FROM users) as total_users,
			(SELECT COUNT(*) FROM post) as total_posts,
			(SELECT COUNT(*) FROM spotted) as total_spotted,
			(SELECT COUNT(*) FROM post WHERE status = (SELECT id FROM submit_status WHERE description='approved')) as approved_posts,
			(SELECT COUNT(*) FROM spotted WHERE status = (SELECT id FROM submit_status WHERE description='approved')) as approved_spotted,
			(SELECT COUNT(*) FROM post_like) as total_post_likes,
			(SELECT COUNT(*) FROM spotted_like) as total_spotted_likes,
			(SELECT COUNT(*) FROM cities) as total_cities,
			(SELECT COUNT(*) FROM schools) as total_schools`,
	)
}

// Statistiche per città
func QueryStatsByCity(db *pgx.Conn) (pgx.Rows, error) {
	return db.Query(
		context.Background(),
		`SELECT
			c.id, c.name, c.region,
			COUNT(DISTINCT utsc.user_id) as user_count,
			(SELECT COUNT(*) FROM schools WHERE city = c.id) as school_count,
			COALESCE((SELECT COUNT(*) FROM post p
				JOIN users u ON p.creator = u.id
				JOIN users_to_school_cities utsc2 ON u.id = utsc2.user_id
				WHERE utsc2.city = c.id), 0) as post_count,
			COALESCE((SELECT COUNT(*) FROM spotted s
				JOIN users u ON s.creator = u.id
				JOIN users_to_school_cities utsc2 ON u.id = utsc2.user_id
				WHERE utsc2.city = c.id), 0) as spotted_count
		 FROM cities c
		 LEFT JOIN users_to_school_cities utsc ON utsc.city = c.id
		 GROUP BY c.id, c.name, c.region
		 ORDER BY user_count DESC`,
	)
}

// Statistiche per scuola
func QueryStatsBySchool(db *pgx.Conn) (pgx.Rows, error) {
	return db.Query(
		context.Background(),
		`SELECT
			s.id, s.name, c.id as city_id, c.name as city_name,
			COUNT(DISTINCT utsc.user_id) as user_count,
			COALESCE((SELECT COUNT(*) FROM post p
				JOIN users u ON p.creator = u.id
				JOIN users_to_school_cities utsc2 ON u.id = utsc2.user_id
				WHERE utsc2.school = s.id), 0) as post_count,
			COALESCE((SELECT COUNT(*) FROM spotted sp
				JOIN users u ON sp.creator = u.id
				JOIN users_to_school_cities utsc2 ON u.id = utsc2.user_id
				WHERE utsc2.school = s.id), 0) as spotted_count
		 FROM schools s
		 JOIN cities c ON s.city = c.id
		 LEFT JOIN users_to_school_cities utsc ON utsc.school = s.id
		 GROUP BY s.id, s.name, c.id, c.name
		 ORDER BY user_count DESC`,
	)
}

// Statistiche temporali - utenti per mese
func QueryUsersOverTime(db *pgx.Conn) (pgx.Rows, error) {
	return db.Query(
		context.Background(),
		`SELECT
			DATE_TRUNC('month', creation_timestamp) as month,
			COUNT(*) as count
		 FROM users
		 WHERE creation_timestamp IS NOT NULL
		 GROUP BY DATE_TRUNC('month', creation_timestamp)
		 ORDER BY month DESC
		 LIMIT 12`,
	)
}

// Statistiche temporali - post per mese
func QueryPostsOverTime(db *pgx.Conn) (pgx.Rows, error) {
	return db.Query(
		context.Background(),
		`SELECT
			DATE_TRUNC('month', creation_timestamp) as month,
			COUNT(*) as count
		 FROM post
		 WHERE creation_timestamp IS NOT NULL
		 GROUP BY DATE_TRUNC('month', creation_timestamp)
		 ORDER BY month DESC
		 LIMIT 12`,
	)
}

// Statistiche temporali - spotted per mese
func QuerySpottedOverTime(db *pgx.Conn) (pgx.Rows, error) {
	return db.Query(
		context.Background(),
		`SELECT
			DATE_TRUNC('month', creation_timestamp) as month,
			COUNT(*) as count
		 FROM spotted
		 WHERE creation_timestamp IS NOT NULL
		 GROUP BY DATE_TRUNC('month', creation_timestamp)
		 ORDER BY month DESC
		 LIMIT 12`,
	)
}

// Top 10 città per utenti
func QueryTopCities(db *pgx.Conn, limit int) (pgx.Rows, error) {
	return db.Query(
		context.Background(),
		`SELECT
			c.id, c.name, c.region,
			COUNT(DISTINCT utsc.user_id) as user_count
		 FROM cities c
		 LEFT JOIN users_to_school_cities utsc ON utsc.city = c.id
		 GROUP BY c.id, c.name, c.region
		 ORDER BY user_count DESC
		 LIMIT $1`,
		limit,
	)
}

// Top 10 scuole per utenti
func QueryTopSchools(db *pgx.Conn, limit int) (pgx.Rows, error) {
	return db.Query(
		context.Background(),
		`SELECT
			s.id, s.name, c.name as city_name,
			COUNT(DISTINCT utsc.user_id) as user_count
		 FROM schools s
		 JOIN cities c ON s.city = c.id
		 LEFT JOIN users_to_school_cities utsc ON utsc.school = s.id
		 GROUP BY s.id, s.name, c.name
		 ORDER BY user_count DESC
		 LIMIT $1`,
		limit,
	)
}
