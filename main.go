package main

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	_ "modernc.org/sqlite"

	"github.com/golang-jwt/jwt/v5"
)

// Application dependencies container to enable simple dependency injection across layers.
type Application struct {
	Database      *sql.DB
	JWTSigningKey []byte
	HTTPServer    *http.Server
}

func main() {
	// Configure base logger output
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	cfg := loadConfigFromEnv()

	db, err := openSQLiteDatabase(cfg.DatabasePath)
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}
	defer func() { _ = db.Close() }()

	if err := runMigrations(db, "migrations"); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	app := &Application{
		Database:      db,
		JWTSigningKey: []byte(cfg.JWTSecret),
	}

	mux := http.NewServeMux()

	// Middlewares: security headers then logging
	handler := withSecurityHeaders(withLogging(mux))

	// Register routes
	registerRoutes(mux, app)

	srv := &http.Server{
		Addr:              fmt.Sprintf(":%d", cfg.HTTPPort),
		Handler:           handler,
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
		ErrorLog:          log.Default(),
	}
	app.HTTPServer = srv

	// Graceful shutdown support
	go func() {
		log.Printf("server starting on :%d", cfg.HTTPPort)
		if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("server error: %v", err)
		}
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("graceful shutdown failed: %v", err)
	}
}

// -------------------- Configuration --------------------

type config struct {
	HTTPPort     int
	DatabasePath string
	JWTSecret    string
}

func loadConfigFromEnv() config {
	port := 8080
	if v := os.Getenv("PORT"); v != "" {
		var parsed int
		_, _ = fmt.Sscanf(v, "%d", &parsed)
		if parsed > 0 {
			port = parsed
		}
	}
	dbPath := os.Getenv("DATABASE_PATH")
	if dbPath == "" {
		dbPath = filepath.Join(".", "new-api.db")
	}
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		// Development default only. Override in production via env.
		jwtSecret = "dev-secret-change-me"
	}
	return config{HTTPPort: port, DatabasePath: dbPath, JWTSecret: jwtSecret}
}

// -------------------- Database --------------------

func openSQLiteDatabase(dbPath string) (*sql.DB, error) {
	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(dbPath), 0o755); err != nil {
		return nil, fmt.Errorf("create db dir: %w", err)
	}
	dsn := fmt.Sprintf("file:%s?_pragma=foreign_keys(1)&_pragma=busy_timeout(5000)", dbPath)
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("open sqlite: %w", err)
	}
	if err := db.Ping(); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("ping sqlite: %w", err)
	}
	return db, nil
}

// -------------------- JWT Helpers --------------------

type jwtClaims struct {
	UserID int64 `json:"userId"`
	jwt.RegisteredClaims
}

// -------------------- Middleware --------------------

func withLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		duration := time.Since(start)
		log.Printf("%s %s %s", r.Method, r.URL.Path, duration)
	})
}

func withSecurityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Minimal secure headers for API responses
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-XSS-Protection", "0")
		w.Header().Set("Referrer-Policy", "no-referrer")
		w.Header().Set("Content-Security-Policy", "default-src 'none'")
		next.ServeHTTP(w, r)
	})
}

// -------------------- Routing --------------------

func registerRoutes(mux *http.ServeMux, app *Application) {
	// Health check
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte("{\"status\":\"ok\"}"))
	})

	// Auth and Users
	mux.HandleFunc("/register", app.handleRegister)
	mux.Handle("/users", app.authenticated(http.HandlerFunc(app.handleUsersCollection)))
	mux.Handle("/users/", app.authenticated(http.HandlerFunc(app.handleUserByID)))

	// Habits
	mux.Handle("/habits", app.authenticated(http.HandlerFunc(app.handleHabitsCollection)))
	mux.Handle("/habits/", app.authenticated(http.HandlerFunc(app.handleHabitByID)))

	// Logs
	mux.Handle("/logs", app.authenticated(http.HandlerFunc(app.handleLogsCollection)))
	mux.Handle("/logs/", app.authenticated(http.HandlerFunc(app.handleLogByID)))
}

// -------------------- Auth Middleware --------------------

func (a *Application) authenticated(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		const prefix = "Bearer "
		if len(header) <= len(prefix) || header[:len(prefix)] != prefix {
			http.Error(w, "missing or invalid Authorization header", http.StatusUnauthorized)
			return
		}
		tokenString := header[len(prefix):]
		token, err := jwt.ParseWithClaims(tokenString, &jwtClaims{}, func(token *jwt.Token) (interface{}, error) {
			return a.JWTSigningKey, nil
		})
		if err != nil || !token.Valid {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}
		claims, ok := token.Claims.(*jwtClaims)
		if !ok {
			http.Error(w, "invalid token claims", http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), contextKeyUserID, claims.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// -------------------- Context Keys --------------------

type contextKey string

const contextKeyUserID contextKey = "userID"

func getUserIDFromContext(ctx context.Context) (int64, bool) {
	v := ctx.Value(contextKeyUserID)
	if v == nil {
		return 0, false
	}
	id, ok := v.(int64)
	return id, ok
}

// -------------------- HTTP Helpers --------------------

func writeJSON(w http.ResponseWriter, status int, body string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, _ = w.Write([]byte(body))
}

func writeJSONError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, fmt.Sprintf("{\"error\":\"%s\"}", escapeForJSON(message)))
}

func escapeForJSON(s string) string {
	// Minimal escaping for error messages
	b := make([]rune, 0, len(s))
	for _, r := range s {
		switch r {
		case '\\':
			b = append(b, '\\', '\\')
		case '"':
			b = append(b, '\\', '"')
		case '\n':
			b = append(b, '\\', 'n')
		case '\r':
			b = append(b, '\\', 'r')
		case '\t':
			b = append(b, '\\', 't')
		default:
			b = append(b, r)
		}
	}
	return string(b)
}

// -------------------- Models --------------------

type User struct {
	ID              int64     `json:"id"`
	ProfileImageURL string    `json:"profileImageUrl"`
	Name            string    `json:"name"`
	TimeZone        string    `json:"timeZone"`
	PhoneNumber     string    `json:"phoneNumber"`
	CreatedAt       time.Time `json:"createdAt"`
}

type Habit struct {
	ID          int64     `json:"id"`
	UserID      int64     `json:"userId"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"createdAt"`
}

type LogEntry struct {
	ID        int64     `json:"id"`
	HabitID   int64     `json:"habitId"`
	Notes     string    `json:"notes"`
	CreatedAt time.Time `json:"createdAt"`
}

// -------------------- Handlers --------------------

// handleRegister creates a user and returns a JWT for subsequent authenticated calls.
func (a *Application) handleRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSONError(w, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}
	type payload struct {
		ProfileImageURL string `json:"profileImageUrl"`
		Name            string `json:"name"`
		TimeZone        string `json:"timeZone"`
		PhoneNumber     string `json:"phoneNumber"`
	}
	var p payload
	if err := jsonDecodeStrict(r, &p); err != nil {
		log.Printf("register decode error: %v", err)
		writeJSONError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if p.Name == "" || p.PhoneNumber == "" || p.TimeZone == "" {
		writeJSONError(w, http.StatusBadRequest, "name, phoneNumber, and timeZone are required")
		return
	}
	now := time.Now().UTC()
	res, err := a.Database.Exec(
		`INSERT INTO users (profile_image_url, name, time_zone, phone_number, created_at) VALUES (?, ?, ?, ?, ?)`,
		p.ProfileImageURL, p.Name, p.TimeZone, p.PhoneNumber, now,
	)
	if err != nil {
		log.Printf("register insert error: %v", err)
		writeJSONError(w, http.StatusInternalServerError, "failed to create user")
		return
	}
	id, _ := res.LastInsertId()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtClaims{
		UserID: id,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(24 * time.Hour)),
		},
	})
	tokenString, err := token.SignedString(a.JWTSigningKey)
	if err != nil {
		log.Printf("jwt sign error: %v", err)
		writeJSONError(w, http.StatusInternalServerError, "failed to issue token")
		return
	}
	writeJSON(w, http.StatusCreated, fmt.Sprintf(`{"userId":%d,"token":"%s"}`, id, tokenString))
}

// Users collection supports GET (list self only) not required, but keep for completeness; POST not allowed here (use /register)
func (a *Application) handleUsersCollection(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		userID, _ := getUserIDFromContext(r.Context())
		row := a.Database.QueryRow(`SELECT id, profile_image_url, name, time_zone, phone_number, created_at FROM users WHERE id = ?`, userID)
		var u User
		if err := row.Scan(&u.ID, &u.ProfileImageURL, &u.Name, &u.TimeZone, &u.PhoneNumber, &u.CreatedAt); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				writeJSONError(w, http.StatusNotFound, "user not found")
				return
			}
			log.Printf("user fetch error: %v", err)
			writeJSONError(w, http.StatusInternalServerError, "failed to fetch user")
			return
		}
		b := fmt.Sprintf(`{"id":%d,"profileImageUrl":"%s","name":"%s","timeZone":"%s","phoneNumber":"%s","createdAt":"%s"}`,
			u.ID, escapeForJSON(u.ProfileImageURL), escapeForJSON(u.Name), escapeForJSON(u.TimeZone), escapeForJSON(u.PhoneNumber), u.CreatedAt.Format(time.RFC3339))
		writeJSON(w, http.StatusOK, b)
	default:
		writeJSONError(w, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
	}
}

// User by ID allows GET, PUT, DELETE but only for self.
func (a *Application) handleUserByID(w http.ResponseWriter, r *http.Request) {
	id, ok := parseIDParam(r.URL.Path)
	if !ok {
		writeJSONError(w, http.StatusBadRequest, "invalid id")
		return
	}
	userID, _ := getUserIDFromContext(r.Context())
	if id != userID {
		writeJSONError(w, http.StatusForbidden, "cannot access other users")
		return
	}
	switch r.Method {
	case http.MethodGet:
		a.handleUsersCollection(w, r)
	case http.MethodPut:
		type payload struct {
			ProfileImageURL *string `json:"profileImageUrl"`
			Name            *string `json:"name"`
			TimeZone        *string `json:"timeZone"`
			PhoneNumber     *string `json:"phoneNumber"`
		}
		var p payload
		if err := jsonDecodeStrict(r, &p); err != nil {
			writeJSONError(w, http.StatusBadRequest, "invalid request body")
			return
		}
		// Build dynamic update
		fields := make([]string, 0, 4)
		args := make([]any, 0, 5)
		if p.ProfileImageURL != nil {
			fields = append(fields, "profile_image_url = ?")
			args = append(args, *p.ProfileImageURL)
		}
		if p.Name != nil {
			fields = append(fields, "name = ?")
			args = append(args, *p.Name)
		}
		if p.TimeZone != nil {
			fields = append(fields, "time_zone = ?")
			args = append(args, *p.TimeZone)
		}
		if p.PhoneNumber != nil {
			fields = append(fields, "phone_number = ?")
			args = append(args, *p.PhoneNumber)
		}
		if len(fields) == 0 {
			writeJSONError(w, http.StatusBadRequest, "no fields to update")
			return
		}
		args = append(args, id)
		q := "UPDATE users SET " + joinComma(fields) + " WHERE id = ?"
		if _, err := a.Database.Exec(q, args...); err != nil {
			log.Printf("user update error: %v", err)
			writeJSONError(w, http.StatusInternalServerError, "failed to update user")
			return
		}
		w.WriteHeader(http.StatusNoContent)
	case http.MethodDelete:
		if _, err := a.Database.Exec(`DELETE FROM users WHERE id = ?`, id); err != nil {
			log.Printf("user delete error: %v", err)
			writeJSONError(w, http.StatusInternalServerError, "failed to delete user")
			return
		}
		w.WriteHeader(http.StatusNoContent)
	default:
		writeJSONError(w, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
	}
}

// Habits collection: GET (list self), POST (create)
func (a *Application) handleHabitsCollection(w http.ResponseWriter, r *http.Request) {
	userID, _ := getUserIDFromContext(r.Context())
	switch r.Method {
	case http.MethodGet:
		rows, err := a.Database.Query(`SELECT id, user_id, name, description, created_at FROM habits WHERE user_id = ? ORDER BY id DESC`, userID)
		if err != nil {
			log.Printf("habits list error: %v", err)
			writeJSONError(w, http.StatusInternalServerError, "failed to list habits")
			return
		}
		defer rows.Close()
		first := true
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("["))
		for rows.Next() {
			var h Habit
			if err := rows.Scan(&h.ID, &h.UserID, &h.Name, &h.Description, &h.CreatedAt); err != nil {
				continue
			}
			if !first {
				_, _ = w.Write([]byte(","))
			} else {
				first = false
			}
			_, _ = w.Write([]byte(fmt.Sprintf(`{"id":%d,"userId":%d,"name":"%s","description":"%s","createdAt":"%s"}`,
				h.ID, h.UserID, escapeForJSON(h.Name), escapeForJSON(h.Description), h.CreatedAt.Format(time.RFC3339))))
		}
		_, _ = w.Write([]byte("]"))
	case http.MethodPost:
		type payload struct {
			Name        string `json:"name"`
			Description string `json:"description"`
		}
		var p payload
		if err := jsonDecodeStrict(r, &p); err != nil {
			writeJSONError(w, http.StatusBadRequest, "invalid request body")
			return
		}
		if p.Name == "" {
			writeJSONError(w, http.StatusBadRequest, "name is required")
			return
		}
		now := time.Now().UTC()
		res, err := a.Database.Exec(`INSERT INTO habits (user_id, name, description, created_at) VALUES (?, ?, ?, ?)`, userID, p.Name, p.Description, now)
		if err != nil {
			log.Printf("habit insert error: %v", err)
			writeJSONError(w, http.StatusInternalServerError, "failed to create habit")
			return
		}
		id, _ := res.LastInsertId()
		writeJSON(w, http.StatusCreated, fmt.Sprintf(`{"id":%d}`, id))
	default:
		writeJSONError(w, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
	}
}

// Habit by ID: GET, PUT, DELETE
func (a *Application) handleHabitByID(w http.ResponseWriter, r *http.Request) {
	id, ok := parseIDParam(r.URL.Path)
	if !ok {
		writeJSONError(w, http.StatusBadRequest, "invalid id")
		return
	}
	userID, _ := getUserIDFromContext(r.Context())
	// Ownership check
	var ownerID int64
	row := a.Database.QueryRow(`SELECT user_id FROM habits WHERE id = ?`, id)
	if err := row.Scan(&ownerID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			writeJSONError(w, http.StatusNotFound, "habit not found")
			return
		}
		log.Printf("habit owner fetch error: %v", err)
		writeJSONError(w, http.StatusInternalServerError, "failed to fetch habit")
		return
	}
	if ownerID != userID {
		writeJSONError(w, http.StatusForbidden, "cannot access other users' habits")
		return
	}
	switch r.Method {
	case http.MethodGet:
		row := a.Database.QueryRow(`SELECT id, user_id, name, description, created_at FROM habits WHERE id = ?`, id)
		var h Habit
		if err := row.Scan(&h.ID, &h.UserID, &h.Name, &h.Description, &h.CreatedAt); err != nil {
			writeJSONError(w, http.StatusInternalServerError, "failed to fetch habit")
			return
		}
		writeJSON(w, http.StatusOK, fmt.Sprintf(`{"id":%d,"userId":%d,"name":"%s","description":"%s","createdAt":"%s"}`,
			h.ID, h.UserID, escapeForJSON(h.Name), escapeForJSON(h.Description), h.CreatedAt.Format(time.RFC3339)))
	case http.MethodPut:
		type payload struct {
			Name        *string `json:"name"`
			Description *string `json:"description"`
		}
		var p payload
		if err := jsonDecodeStrict(r, &p); err != nil {
			writeJSONError(w, http.StatusBadRequest, "invalid request body")
			return
		}
		fields := make([]string, 0, 2)
		args := make([]any, 0, 3)
		if p.Name != nil {
			fields = append(fields, "name = ?")
			args = append(args, *p.Name)
		}
		if p.Description != nil {
			fields = append(fields, "description = ?")
			args = append(args, *p.Description)
		}
		if len(fields) == 0 {
			writeJSONError(w, http.StatusBadRequest, "no fields to update")
			return
		}
		args = append(args, id)
		q := "UPDATE habits SET " + joinComma(fields) + " WHERE id = ?"
		if _, err := a.Database.Exec(q, args...); err != nil {
			log.Printf("habit update error: %v", err)
			writeJSONError(w, http.StatusInternalServerError, "failed to update habit")
			return
		}
		w.WriteHeader(http.StatusNoContent)
	case http.MethodDelete:
		if _, err := a.Database.Exec(`DELETE FROM habits WHERE id = ?`, id); err != nil {
			log.Printf("habit delete error: %v", err)
			writeJSONError(w, http.StatusInternalServerError, "failed to delete habit")
			return
		}
		w.WriteHeader(http.StatusNoContent)
	default:
		writeJSONError(w, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
	}
}

// Logs collection: GET with optional habitId, POST
func (a *Application) handleLogsCollection(w http.ResponseWriter, r *http.Request) {
	userID, _ := getUserIDFromContext(r.Context())
	switch r.Method {
	case http.MethodGet:
		// Optional filter by habitId
		habitID, hasHabit := parseQueryID(r, "habitId")
		if hasHabit {
			// Verify ownership
			var ownerID int64
			row := a.Database.QueryRow(`SELECT user_id FROM habits WHERE id = ?`, habitID)
			if err := row.Scan(&ownerID); err != nil {
				writeJSONError(w, http.StatusNotFound, "habit not found")
				return
			}
			if ownerID != userID {
				writeJSONError(w, http.StatusForbidden, "cannot access other users' habits")
				return
			}
			rows, err := a.Database.Query(`SELECT id, habit_id, notes, created_at FROM logs WHERE habit_id = ? ORDER BY id DESC`, habitID)
			if err != nil {
				log.Printf("logs list error: %v", err)
				writeJSONError(w, http.StatusInternalServerError, "failed to list logs")
				return
			}
			defer rows.Close()
			first := true
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("["))
			for rows.Next() {
				var le LogEntry
				if err := rows.Scan(&le.ID, &le.HabitID, &le.Notes, &le.CreatedAt); err != nil {
					continue
				}
				if !first {
					_, _ = w.Write([]byte(","))
				} else {
					first = false
				}
				_, _ = w.Write([]byte(fmt.Sprintf(`{"id":%d,"habitId":%d,"notes":"%s","createdAt":"%s"}`,
					le.ID, le.HabitID, escapeForJSON(le.Notes), le.CreatedAt.Format(time.RFC3339))))
			}
			_, _ = w.Write([]byte("]"))
			return
		}
		// No habit filter: list all logs for user's habits
		rows, err := a.Database.Query(`SELECT l.id, l.habit_id, l.notes, l.created_at FROM logs l JOIN habits h ON h.id = l.habit_id WHERE h.user_id = ? ORDER BY l.id DESC`, userID)
		if err != nil {
			log.Printf("logs list error: %v", err)
			writeJSONError(w, http.StatusInternalServerError, "failed to list logs")
			return
		}
		defer rows.Close()
		first := true
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("["))
		for rows.Next() {
			var le LogEntry
			if err := rows.Scan(&le.ID, &le.HabitID, &le.Notes, &le.CreatedAt); err != nil {
				continue
			}
			if !first {
				_, _ = w.Write([]byte(","))
			} else {
				first = false
			}
			_, _ = w.Write([]byte(fmt.Sprintf(`{"id":%d,"habitId":%d,"notes":"%s","createdAt":"%s"}`,
				le.ID, le.HabitID, escapeForJSON(le.Notes), le.CreatedAt.Format(time.RFC3339))))
		}
		_, _ = w.Write([]byte("]"))
	case http.MethodPost:
		type payload struct {
			HabitID int64  `json:"habitId"`
			Notes   string `json:"notes"`
		}
		var p payload
		if err := jsonDecodeStrict(r, &p); err != nil {
			writeJSONError(w, http.StatusBadRequest, "invalid request body")
			return
		}
		if p.HabitID == 0 {
			writeJSONError(w, http.StatusBadRequest, "habitId is required")
			return
		}
		// Ownership verify
		var ownerID int64
		row := a.Database.QueryRow(`SELECT user_id FROM habits WHERE id = ?`, p.HabitID)
		if err := row.Scan(&ownerID); err != nil {
			writeJSONError(w, http.StatusNotFound, "habit not found")
			return
		}
		if ownerID != userID {
			writeJSONError(w, http.StatusForbidden, "cannot access other users' habits")
			return
		}
		now := time.Now().UTC()
		res, err := a.Database.Exec(`INSERT INTO logs (habit_id, notes, created_at) VALUES (?, ?, ?)`, p.HabitID, p.Notes, now)
		if err != nil {
			log.Printf("log insert error: %v", err)
			writeJSONError(w, http.StatusInternalServerError, "failed to create log")
			return
		}
		id, _ := res.LastInsertId()
		writeJSON(w, http.StatusCreated, fmt.Sprintf(`{"id":%d}`, id))
	default:
		writeJSONError(w, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
	}
}

// Log by ID: GET, PUT, DELETE with ownership across join to habits
func (a *Application) handleLogByID(w http.ResponseWriter, r *http.Request) {
	id, ok := parseIDParam(r.URL.Path)
	if !ok {
		writeJSONError(w, http.StatusBadRequest, "invalid id")
		return
	}
	userID, _ := getUserIDFromContext(r.Context())
	// Verify ownership via join
	var habitID int64
	row := a.Database.QueryRow(`SELECT habit_id FROM logs WHERE id = ?`, id)
	if err := row.Scan(&habitID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			writeJSONError(w, http.StatusNotFound, "log not found")
			return
		}
		writeJSONError(w, http.StatusInternalServerError, "failed to fetch log")
		return
	}
	var ownerID int64
	row2 := a.Database.QueryRow(`SELECT user_id FROM habits WHERE id = ?`, habitID)
	if err := row2.Scan(&ownerID); err != nil {
		writeJSONError(w, http.StatusInternalServerError, "failed to verify ownership")
		return
	}
	if ownerID != userID {
		writeJSONError(w, http.StatusForbidden, "cannot access other users' logs")
		return
	}
	switch r.Method {
	case http.MethodGet:
		row := a.Database.QueryRow(`SELECT id, habit_id, notes, created_at FROM logs WHERE id = ?`, id)
		var le LogEntry
		if err := row.Scan(&le.ID, &le.HabitID, &le.Notes, &le.CreatedAt); err != nil {
			writeJSONError(w, http.StatusInternalServerError, "failed to fetch log")
			return
		}
		writeJSON(w, http.StatusOK, fmt.Sprintf(`{"id":%d,"habitId":%d,"notes":"%s","createdAt":"%s"}`,
			le.ID, le.HabitID, escapeForJSON(le.Notes), le.CreatedAt.Format(time.RFC3339)))
	case http.MethodPut:
		type payload struct {
			Notes *string `json:"notes"`
		}
		var p payload
		if err := jsonDecodeStrict(r, &p); err != nil {
			writeJSONError(w, http.StatusBadRequest, "invalid request body")
			return
		}
		if p.Notes == nil {
			writeJSONError(w, http.StatusBadRequest, "no fields to update")
			return
		}
		if _, err := a.Database.Exec(`UPDATE logs SET notes = ? WHERE id = ?`, *p.Notes, id); err != nil {
			log.Printf("log update error: %v", err)
			writeJSONError(w, http.StatusInternalServerError, "failed to update log")
			return
		}
		w.WriteHeader(http.StatusNoContent)
	case http.MethodDelete:
		if _, err := a.Database.Exec(`DELETE FROM logs WHERE id = ?`, id); err != nil {
			log.Printf("log delete error: %v", err)
			writeJSONError(w, http.StatusInternalServerError, "failed to delete log")
			return
		}
		w.WriteHeader(http.StatusNoContent)
	default:
		writeJSONError(w, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
	}
}

// -------------------- Utilities --------------------

func parseIDParam(path string) (int64, bool) {
	// expects "/resource/{id}" style; extract last segment as ID
	i := len(path) - 1
	for i >= 0 && path[i] == '/' {
		i--
	}
	j := i
	for j >= 0 && path[j] != '/' {
		j--
	}
	if j < 0 || j == i {
		return 0, false
	}
	var id int64
	if _, err := fmt.Sscanf(path[j+1:i+1], "%d", &id); err != nil {
		return 0, false
	}
	return id, true
}

func joinComma(parts []string) string {
	if len(parts) == 0 {
		return ""
	}
	out := parts[0]
	for i := 1; i < len(parts); i++ {
		out += ", " + parts[i]
	}
	return out
}

func parseQueryID(r *http.Request, key string) (int64, bool) {
	v := r.URL.Query().Get(key)
	if v == "" {
		return 0, false
	}
	var id int64
	if _, err := fmt.Sscanf(v, "%d", &id); err != nil {
		return 0, false
	}
	return id, true
}

// Strict JSON decoder that rejects unknown fields and limits size.
func jsonDecodeStrict(r *http.Request, dst any) error {
	const maxBody = 1 << 20 // 1MB
	defer r.Body.Close()
	lr := io.LimitReader(r.Body, maxBody+1)
	data, err := io.ReadAll(lr)
	if err != nil {
		return err
	}
	if int64(len(data)) > maxBody {
		return fmt.Errorf("request body too large")
	}
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields()
	return dec.Decode(dst)
}

// -------------------- Migrations --------------------

func runMigrations(db *sql.DB, dir string) error {
	// Ensure migrations bookkeeping table exists
	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS schema_migrations (filename TEXT PRIMARY KEY, applied_at TIMESTAMP NOT NULL)`); err != nil {
		return fmt.Errorf("create schema_migrations: %w", err)
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("read migrations dir: %w", err)
	}
	// Simple lexicographic order (001_, 002_, ...)
	for _, e := range entries {
		if e.IsDir() || filepath.Ext(e.Name()) != ".sql" {
			continue
		}
		var exists int
		row := db.QueryRow(`SELECT COUNT(1) FROM schema_migrations WHERE filename = ?`, e.Name())
		if err := row.Scan(&exists); err != nil {
			return fmt.Errorf("check migration applied: %w", err)
		}
		if exists > 0 {
			continue
		}
		path := filepath.Join(dir, e.Name())
		sqlBytes, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("read migration %s: %w", e.Name(), err)
		}
		tx, err := db.Begin()
		if err != nil {
			return fmt.Errorf("begin tx: %w", err)
		}
		if _, err := tx.Exec(string(sqlBytes)); err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("apply migration %s: %w", e.Name(), err)
		}
		if _, err := tx.Exec(`INSERT INTO schema_migrations (filename, applied_at) VALUES (?, ?)`, e.Name(), time.Now().UTC()); err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("record migration %s: %w", e.Name(), err)
		}
		if err := tx.Commit(); err != nil {
			return fmt.Errorf("commit migration %s: %w", e.Name(), err)
		}
		log.Printf("applied migration: %s", e.Name())
	}
	return nil
}
