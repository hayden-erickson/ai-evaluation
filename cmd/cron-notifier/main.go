package main

import (
    "context"
    "database/sql"
    "errors"
    "fmt"
    "log"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"

    _ "github.com/go-sql-driver/mysql"
    _ "time/tzdata"
)

type AppConfig struct {
    MySQLDSN          string
    MySQLHost         string
    MySQLPort         string
    MySQLUser         string
    MySQLPass         string
    MySQLDB           string
    TwilioAccountSID  string
    TwilioAuthToken   string
    TwilioFromNumber  string
    DryRun            bool
}

type User struct {
    ID         int64
    Name       string
    Phone      string
    TimeZoneID string
}

func getenv(key, def string) string {
    v := os.Getenv(key)
    if v == "" {
        return def
    }
    return v
}

func mustGetenv(key string) (string, error) {
    v := os.Getenv(key)
    if v == "" {
        return "", fmt.Errorf("missing required env %s", key)
    }
    return v, nil
}

func loadConfig() (AppConfig, error) {
    cfg := AppConfig{}
    cfg.MySQLDSN = getenv("MYSQL_DSN", "")
    cfg.MySQLHost = getenv("MYSQL_HOST", "")
    cfg.MySQLPort = getenv("MYSQL_PORT", "3306")
    cfg.MySQLUser = getenv("MYSQL_USER", "")
    cfg.MySQLPass = getenv("MYSQL_PASSWORD", "")
    cfg.MySQLDB = getenv("MYSQL_DB", "")
    cfg.TwilioAccountSID = getenv("TWILIO_ACCOUNT_SID", "")
    cfg.TwilioAuthToken = getenv("TWILIO_AUTH_TOKEN", "")
    cfg.TwilioFromNumber = getenv("TWILIO_FROM_NUMBER", "")
    cfg.DryRun = getenv("DRY_RUN", "false") == "true"

    if cfg.MySQLDSN == "" {
        if cfg.MySQLHost == "" || cfg.MySQLUser == "" || cfg.MySQLDB == "" {
            return cfg, errors.New("either MYSQL_DSN or MYSQL_HOST, MYSQL_USER, MYSQL_DB must be set")
        }
        cfg.MySQLDSN = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&charset=utf8mb4,utf8",
            cfg.MySQLUser, cfg.MySQLPass, cfg.MySQLHost, cfg.MySQLPort, cfg.MySQLDB,
        )
    }

    if cfg.TwilioAccountSID == "" || cfg.TwilioAuthToken == "" || cfg.TwilioFromNumber == "" {
        return cfg, errors.New("missing Twilio configuration: TWILIO_ACCOUNT_SID, TWILIO_AUTH_TOKEN, TWILIO_FROM_NUMBER")
    }
    return cfg, nil
}

func startHealthServer(ctx context.Context, addr string) *http.Server {
    mux := http.NewServeMux()
    mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        _, _ = w.Write([]byte("ok"))
    })
    srv := &http.Server{Addr: addr, Handler: mux}
    go func() {
        if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
            log.Printf("health server error: %v", err)
        }
    }()
    go func() {
        <-ctx.Done()
        shutdownCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
        defer cancel()
        _ = srv.Shutdown(shutdownCtx)
    }()
    return srv
}

func main() {
    log.SetOutput(os.Stdout)
    log.SetFlags(log.LstdFlags | log.LUTC)

    cfg, err := loadConfig()
    if err != nil {
        log.Fatalf("config error: %v", err)
    }

    // graceful shutdown context
    ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
    defer cancel()

    // health server for probes
    _ = startHealthServer(ctx, ":8080")

    // DB connection
    db, err := sql.Open("mysql", cfg.MySQLDSN)
    if err != nil {
        log.Fatalf("db open error: %v", err)
    }
    defer db.Close()
    db.SetConnMaxLifetime(5 * time.Minute)
    db.SetMaxOpenConns(10)
    db.SetMaxIdleConns(5)

    if err := db.PingContext(ctx); err != nil {
        log.Fatalf("db ping error: %v", err)
    }

    // Fetch users
    users, err := fetchUsers(ctx, db)
    if err != nil {
        log.Fatalf("fetch users error: %v", err)
    }
    log.Printf("processing %d users", len(users))

    notifier := newTwilioNotifier(cfg.TwilioAccountSID, cfg.TwilioAuthToken, cfg.TwilioFromNumber, cfg.DryRun)

    sentCount := 0
    for _, u := range users {
        select {
        case <-ctx.Done():
            log.Printf("context cancelled, stopping loop")
            break
        default:
        }

        day1StartUTC, day1EndUTC, day2StartUTC, day2EndUTC := computeTwoDayWindowUTC(u.TimeZoneID)

        day1Count, err := countLogsForUserBetween(ctx, db, u.ID, day1StartUTC, day1EndUTC)
        if err != nil {
            log.Printf("user %d count day1 error: %v", u.ID, err)
            continue
        }
        day2Count, err := countLogsForUserBetween(ctx, db, u.ID, day2StartUTC, day2EndUTC)
        if err != nil {
            log.Printf("user %d count day2 error: %v", u.ID, err)
            continue
        }

        if day1Count == 0 && day2Count == 0 {
            msg := fmt.Sprintf("Hi %s, we noticed no habit logs in the last two days. A small step today can restart your streak!", u.Name)
            if err := notifier.SendSMS(ctx, u.Phone, msg); err != nil {
                log.Printf("send sms error (user %d): %v", u.ID, err)
                continue
            }
            sentCount++
            log.Printf("notif sent for user %d: no logs for two days", u.ID)
            continue
        }

        if day1Count == 1 && day2Count == 0 {
            msg := fmt.Sprintf("Hi %s, great job logging the day before yesterday. No log yesterday—want to add one today?", u.Name)
            if err := notifier.SendSMS(ctx, u.Phone, msg); err != nil {
                log.Printf("send sms error (user %d): %v", u.ID, err)
                continue
            }
            sentCount++
            log.Printf("notif sent for user %d: 1 log day1, 0 logs day2", u.ID)
            continue
        }

        log.Printf("no notification for user %d (day1=%d, day2=%d)", u.ID, day1Count, day2Count)
    }

    log.Printf("done. notifications sent: %d", sentCount)
}

func fetchUsers(ctx context.Context, db *sql.DB) ([]User, error) {
    rows, err := db.QueryContext(ctx, `
        SELECT id, name, phone_number, time_zone
        FROM users
        WHERE phone_number IS NOT NULL AND phone_number <> ''
    `)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    var users []User
    for rows.Next() {
        var u User
        if err := rows.Scan(&u.ID, &u.Name, &u.Phone, &u.TimeZoneID); err != nil {
            return nil, err
        }
        users = append(users, u)
    }
    if err := rows.Err(); err != nil {
        return nil, err
    }
    return users, nil
}

func countLogsForUserBetween(ctx context.Context, db *sql.DB, userID int64, startUTC, endUTC time.Time) (int, error) {
    var count int
    err := db.QueryRowContext(ctx, `
        SELECT COUNT(*)
        FROM logs l
        JOIN habits h ON h.id = l.habit_id
        WHERE h.user_id = ? AND l.created_at >= ? AND l.created_at < ?
    `, userID, startUTC.UTC(), endUTC.UTC()).Scan(&count)
    return count, err
}

func computeTwoDayWindowUTC(timeZoneID string) (time.Time, time.Time, time.Time, time.Time) {
    loc, err := time.LoadLocation(timeZoneID)
    if err != nil {
        loc = time.UTC
    }
    now := time.Now().In(loc)
    // We consider the two most recent FULL days: yesterday and the day before yesterday
    yesterday := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc).Add(-24 * time.Hour)
    dayBeforeYesterday := yesterday.Add(-24 * time.Hour)

    day1StartLocal := dayBeforeYesterday
    day1EndLocal := day1StartLocal.Add(24 * time.Hour)
    day2StartLocal := yesterday
    day2EndLocal := day2StartLocal.Add(24 * time.Hour)

    return day1StartLocal.UTC(), day1EndLocal.UTC(), day2StartLocal.UTC(), day2EndLocal.UTC()
}


