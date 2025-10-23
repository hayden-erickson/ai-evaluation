package http

import (
	"encoding/json"
	"log"
	stdhttp "net/http"
	"strconv"
	"strings"

	"github.com/hayden-erickson/ai-evaluation/internal/http/middleware"
	"github.com/hayden-erickson/ai-evaluation/internal/models"
	"github.com/hayden-erickson/ai-evaluation/internal/service"
	"github.com/hayden-erickson/ai-evaluation/internal/validation"
)

// NewRouter wires handlers and returns an http.Handler.
func NewRouter(auth service.AuthService, users service.UserService, habits service.HabitService, logs service.LogService, jwtSecret, jwtIssuer string) stdhttp.Handler {
	mux := stdhttp.NewServeMux()

	// Auth endpoints (no auth required)
	mux.HandleFunc("/auth/register", func(w stdhttp.ResponseWriter, r *stdhttp.Request) {
		if r.Method != stdhttp.MethodPost { methodNotAllowed(w); return }
		var req struct{
			Email string `json:"email"`
			Password string `json:"password"`
			Name string `json:"name"`
			ProfileImageURL string `json:"profile_image_url"`
			TimeZone string `json:"time_zone"`
			Phone string `json:"phone_number"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil { writeError(w, stdhttp.StatusBadRequest, "invalid json"); return }
		if err := validation.ValidateEmail(req.Email); err != nil { writeError(w, stdhttp.StatusBadRequest, err.Error()); return }
		if err := validation.ValidatePassword(req.Password); err != nil { writeError(w, stdhttp.StatusBadRequest, err.Error()); return }
		u := &models.User{Email: req.Email, Password: req.Password, Name: req.Name, ProfileImageURL: req.ProfileImageURL, TimeZone: req.TimeZone, Phone: req.Phone}
		id, err := auth.Register(r.Context(), u)
		if err != nil { log.Printf("register error: %v", err); writeError(w, stdhttp.StatusInternalServerError, "could not register"); return }
		u.ID = id
		writeJSON(w, stdhttp.StatusCreated, u)
	})

	mux.HandleFunc("/auth/login", func(w stdhttp.ResponseWriter, r *stdhttp.Request) {
		if r.Method != stdhttp.MethodPost { methodNotAllowed(w); return }
		var req struct{ Email, Password string }
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil { writeError(w, stdhttp.StatusBadRequest, "invalid json"); return }
		if err := validation.ValidateEmail(req.Email); err != nil { writeError(w, stdhttp.StatusBadRequest, err.Error()); return }
		tok, user, err := auth.Login(r.Context(), req.Email, req.Password)
		if err != nil { log.Printf("login error: %v", err); writeError(w, stdhttp.StatusUnauthorized, "invalid credentials"); return }
		writeJSON(w, stdhttp.StatusOK, map[string]any{"token": tok, "user": user})
	})

	// Authenticated subrouter
	authenticated := func(h stdhttp.HandlerFunc) stdhttp.Handler {
		return middleware.Authenticate(jwtSecret, jwtIssuer, h)
	}

	// Users - operate on current user only
	mux.Handle("/users/me", authenticated(func(w stdhttp.ResponseWriter, r *stdhttp.Request) {
		uid := middleware.GetUserID(r)
		if uid == 0 { writeError(w, stdhttp.StatusUnauthorized, "unauthorized"); return }
		switch r.Method {
		case stdhttp.MethodGet:
			u, err := users.GetByID(r.Context(), uid)
			if err != nil { log.Printf("get user error: %v", err); writeError(w, stdhttp.StatusInternalServerError, "server error"); return }
			if u == nil { writeError(w, stdhttp.StatusNotFound, "not found"); return }
			writeJSON(w, stdhttp.StatusOK, u)
		case stdhttp.MethodPut:
			var req struct{
				ProfileImageURL string `json:"profile_image_url"`
				Name string `json:"name"`
				TimeZone string `json:"time_zone"`
				Phone string `json:"phone_number"`
			}
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil { writeError(w, stdhttp.StatusBadRequest, "invalid json"); return }
			u := &models.User{ID: uid, ProfileImageURL: req.ProfileImageURL, Name: req.Name, TimeZone: req.TimeZone, Phone: req.Phone}
			if err := users.UpdateSelf(r.Context(), u); err != nil { log.Printf("update user error: %v", err); writeError(w, stdhttp.StatusInternalServerError, "server error"); return }
			writeJSON(w, stdhttp.StatusOK, map[string]string{"status": "updated"})
		case stdhttp.MethodDelete:
			if err := users.DeleteSelf(r.Context(), uid); err != nil { log.Printf("delete user error: %v", err); writeError(w, stdhttp.StatusInternalServerError, "server error"); return }
			writeJSON(w, stdhttp.StatusOK, map[string]string{"status": "deleted"})
		default:
			methodNotAllowed(w)
		}
	}))

	// Habits
	mux.Handle("/habits", authenticated(func(w stdhttp.ResponseWriter, r *stdhttp.Request) {
		uid := middleware.GetUserID(r)
		switch r.Method {
		case stdhttp.MethodPost:
			var req struct{ Name, Description string }
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil { writeError(w, stdhttp.StatusBadRequest, "invalid json"); return }
			if err := validation.ValidateNonEmpty(req.Name, "name"); err != nil { writeError(w, stdhttp.StatusBadRequest, err.Error()); return }
			h := &models.Habit{UserID: uid, Name: req.Name, Description: req.Description}
			id, err := habits.Create(r.Context(), h)
			if err != nil { log.Printf("create habit error: %v", err); writeError(w, stdhttp.StatusInternalServerError, "server error"); return }
			h.ID = id
			writeJSON(w, stdhttp.StatusCreated, h)
		case stdhttp.MethodGet:
			list, err := habits.ListByUser(r.Context(), uid)
			if err != nil { log.Printf("list habits error: %v", err); writeError(w, stdhttp.StatusInternalServerError, "server error"); return }
			writeJSON(w, stdhttp.StatusOK, list)
		default:
			methodNotAllowed(w)
		}
	}))

	// Habit by ID and nested logs
	mux.Handle("/habits/", authenticated(func(w stdhttp.ResponseWriter, r *stdhttp.Request) {
		uid := middleware.GetUserID(r)
		// logs subresource detection
		if strings.Contains(r.URL.Path, "/logs") {
			// Match /habits/{id}/logs or /habits/{id}/logs/
			prefix := "/habits/"
			rest := strings.TrimPrefix(r.URL.Path, prefix)
			parts := strings.SplitN(rest, "/", 3)
			if len(parts) < 2 || parts[1] != "logs" { stdhttp.NotFound(w, r); return }
			habitID, err := parseInt64(parts[0])
			if err != nil { writeError(w, stdhttp.StatusBadRequest, "invalid id"); return }
			switch r.Method {
			case stdhttp.MethodPost:
				var req struct{ Notes string }
				if err := json.NewDecoder(r.Body).Decode(&req); err != nil { writeError(w, stdhttp.StatusBadRequest, "invalid json"); return }
				l := &models.LogEntry{HabitID: habitID, Notes: req.Notes}
				id, err := logs.Create(r.Context(), uid, l)
				if err != nil { log.Printf("create log error: %v", err); writeError(w, stdhttp.StatusForbidden, err.Error()); return }
				l.ID = id
				writeJSON(w, stdhttp.StatusCreated, l)
			case stdhttp.MethodGet:
				list, err := logs.ListByHabit(r.Context(), uid, habitID)
				if err != nil { log.Printf("list logs error: %v", err); writeError(w, stdhttp.StatusForbidden, err.Error()); return }
				writeJSON(w, stdhttp.StatusOK, list)
			default:
				methodNotAllowed(w)
			}
			return
		}

		// Habit by ID
		idStr := strings.TrimPrefix(r.URL.Path, "/habits/")
		id, err := parseInt64(idStr)
		if err != nil { writeError(w, stdhttp.StatusBadRequest, "invalid id"); return }
		switch r.Method {
		case stdhttp.MethodGet:
			h, err := habits.GetByID(r.Context(), uid, id)
			if err != nil { log.Printf("get habit error: %v", err); writeError(w, stdhttp.StatusInternalServerError, "server error"); return }
			if h == nil { writeError(w, stdhttp.StatusNotFound, "not found"); return }
			writeJSON(w, stdhttp.StatusOK, h)
		case stdhttp.MethodPut:
			var req struct{ Name, Description string }
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil { writeError(w, stdhttp.StatusBadRequest, "invalid json"); return }
			h := &models.Habit{ID: id, Name: req.Name, Description: req.Description}
			if err := habits.Update(r.Context(), uid, h); err != nil {
				log.Printf("update habit error: %v", err)
				status := stdhttp.StatusInternalServerError
				if strings.Contains(err.Error(), "forbidden") { status = stdhttp.StatusForbidden }
				if strings.Contains(err.Error(), "not found") { status = stdhttp.StatusNotFound }
				writeError(w, status, err.Error())
				return
			}
			writeJSON(w, stdhttp.StatusOK, map[string]string{"status": "updated"})
		case stdhttp.MethodDelete:
			if err := habits.Delete(r.Context(), uid, id); err != nil {
				log.Printf("delete habit error: %v", err)
				status := stdhttp.StatusInternalServerError
				if strings.Contains(err.Error(), "forbidden") { status = stdhttp.StatusForbidden }
				if strings.Contains(err.Error(), "not found") { status = stdhttp.StatusNotFound }
				writeError(w, status, err.Error())
				return
			}
			writeJSON(w, stdhttp.StatusOK, map[string]string{"status": "deleted"})
		default:
			methodNotAllowed(w)
		}
	}))

	// Logs by ID
	mux.Handle("/logs/", authenticated(func(w stdhttp.ResponseWriter, r *stdhttp.Request) {
		uid := middleware.GetUserID(r)
		id, err := parseInt64(strings.TrimPrefix(r.URL.Path, "/logs/"))
		if err != nil { writeError(w, stdhttp.StatusBadRequest, "invalid id"); return }
		switch r.Method {
		case stdhttp.MethodGet:
			l, err := logs.GetByID(r.Context(), uid, id)
			if err != nil { log.Printf("get log error: %v", err); writeError(w, stdhttp.StatusForbidden, err.Error()); return }
			if l == nil { writeError(w, stdhttp.StatusNotFound, "not found"); return }
			writeJSON(w, stdhttp.StatusOK, l)
		case stdhttp.MethodPut:
			var req struct{ Notes string }
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil { writeError(w, stdhttp.StatusBadRequest, "invalid json"); return }
			l := &models.LogEntry{ID: id, Notes: req.Notes}
			if err := logs.Update(r.Context(), uid, l); err != nil { log.Printf("update log error: %v", err); writeError(w, stdhttp.StatusForbidden, err.Error()); return }
			writeJSON(w, stdhttp.StatusOK, map[string]string{"status": "updated"})
		case stdhttp.MethodDelete:
			if err := logs.Delete(r.Context(), uid, id); err != nil { log.Printf("delete log error: %v", err); writeError(w, stdhttp.StatusForbidden, err.Error()); return }
			writeJSON(w, stdhttp.StatusOK, map[string]string{"status": "deleted"})
		default:
			methodNotAllowed(w)
		}
	}))

	return mux
}

func parseInt64(s string) (int64, error) { return strconv.ParseInt(strings.TrimSpace(s), 10, 64) }
