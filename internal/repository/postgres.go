package repository

import (
	"database/sql"

	"mini_url/internal/models"

	"github.com/redis/go-redis/v9"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *models.User) error {
	query := `INSERT INTO users (username, email, password_hash, role, created_at)
			  VALUES ($1, $2, $3, $4, $5) RETURNING id`
	err := r.db.QueryRow(query, user.Username, user.Email, user.Password, user.Role, user.CreatedAt).Scan(&user.ID)
	return err
}

func (r *UserRepository) FindByID(id int) (*models.User, error) {
	query := `SELECT id, username, email, password_hash, role, created_at, last_login FROM users WHERE id = $1`
	row := r.db.QueryRow(query, id)

	user := &models.User{}
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.Role, &user.CreatedAt, &user.LastLogin)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) FindByUsername(username string) (*models.User, error) {
	query := `SELECT id, username, email, password_hash, role, created_at, last_login FROM users WHERE username = $1`
	row := r.db.QueryRow(query, username)

	user := &models.User{}
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.Role, &user.CreatedAt, &user.LastLogin)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
	query := `SELECT id, username, email, password_hash, role, created_at, last_login FROM users WHERE email = $1`
	row := r.db.QueryRow(query, email)

	user := &models.User{}
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.Role, &user.CreatedAt, &user.LastLogin)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) Update(user *models.User) error {
	query := `UPDATE users SET username = $1, email = $2, password_hash = $3, last_login = $4 WHERE id = $5`
	_, err := r.db.Exec(query, user.Username, user.Email, user.Password, user.LastLogin, user.ID)
	return err
}

func (r *UserRepository) Delete(id int) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

type LinkRepository struct {
	db    *sql.DB
	redis *redis.Client
}

func NewLinkRepository(db *sql.DB, redis *redis.Client) *LinkRepository {
	return &LinkRepository{db: db, redis: redis}
}

func (r *LinkRepository) Create(link *models.Link) error {
	query := `INSERT INTO links (original_url, short_id, user_id, created_at, updated_at, status)
			  VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`
	err := r.db.QueryRow(query, link.OriginalURL, link.ShortID, link.UserID, link.CreatedAt, link.UpdatedAt, link.Status).Scan(&link.ID)
	return err
}

func (r *LinkRepository) FindByID(id int) (*models.Link, error) {
	query := `SELECT id, original_url, short_id, user_id, created_at, updated_at, last_accessed_at, click_count, status, expires_at
			  FROM links WHERE id = $1`
	row := r.db.QueryRow(query, id)

	link := &models.Link{}
	err := row.Scan(&link.ID, &link.OriginalURL, &link.ShortID, &link.UserID, &link.CreatedAt, &link.UpdatedAt,
		&link.LastAccessedAt, &link.ClickCount, &link.Status, &link.ExpiresAt)
	if err != nil {
		return nil, err
	}
	return link, nil
}

func (r *LinkRepository) FindByShortID(shortID string) (*models.Link, error) {
	query := `SELECT id, original_url, short_id, user_id, created_at, updated_at, last_accessed_at, click_count, status, expires_at
			  FROM links WHERE short_id = $1`
	row := r.db.QueryRow(query, shortID)

	link := &models.Link{}
	err := row.Scan(&link.ID, &link.OriginalURL, &link.ShortID, &link.UserID, &link.CreatedAt, &link.UpdatedAt,
		&link.LastAccessedAt, &link.ClickCount, &link.Status, &link.ExpiresAt)
	if err != nil {
		return nil, err
	}
	return link, nil
}

func (r *LinkRepository) FindByUserID(userID int) ([]*models.Link, error) {
	query := `SELECT id, original_url, short_id, user_id, created_at, updated_at, last_accessed_at, click_count, status, expires_at
			  FROM links WHERE user_id = $1`
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var links []*models.Link
	for rows.Next() {
		link := &models.Link{}
		err := rows.Scan(&link.ID, &link.OriginalURL, &link.ShortID, &link.UserID, &link.CreatedAt, &link.UpdatedAt,
			&link.LastAccessedAt, &link.ClickCount, &link.Status, &link.ExpiresAt)
		if err != nil {
			return nil, err
		}
		links = append(links, link)
	}
	return links, nil
}

func (r *LinkRepository) Update(link *models.Link) error {
	query := `UPDATE links SET original_url = $1, updated_at = $2, last_accessed_at = $3, click_count = $4, status = $5, expires_at = $6
			  WHERE id = $7`
	_, err := r.db.Exec(query, link.OriginalURL, link.UpdatedAt, link.LastAccessedAt, link.ClickCount, link.Status, link.ExpiresAt, link.ID)
	return err
}

func (r *LinkRepository) Delete(id int) error {
	query := `DELETE FROM links WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

type StatsRepository struct {
	db *sql.DB
}

func NewStatsRepository(db *sql.DB) *StatsRepository {
	return &StatsRepository{db: db}
}

func (r *StatsRepository) FindByShortID(shortID string) ([]*models.Stats, error) {
	query := `SELECT s.id, s.link_id, s.accessed_at, s.ip_address, s.user_agent, s.referer
			  FROM stats s JOIN links l ON s.link_id = l.id WHERE l.short_id = $1`
	rows, err := r.db.Query(query, shortID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stats []*models.Stats
	for rows.Next() {
		st := &models.Stats{}
		err := rows.Scan(&st.ID, &st.LinkID, &st.AccessedAt, &st.IPAddress, &st.UserAgent, &st.Referer)
		if err != nil {
			return nil, err
		}
		stats = append(stats, st)
	}
	return stats, nil
}

func (r *StatsRepository) Create(stats *models.Stats) error {
	query := `INSERT INTO stats (link_id, accessed_at, ip_address, user_agent, referer) VALUES ($1, $2, $3, $4, $5)`
	_, err := r.db.Exec(query, stats.LinkID, stats.AccessedAt, stats.IPAddress, stats.UserAgent, stats.Referer)
	return err
}