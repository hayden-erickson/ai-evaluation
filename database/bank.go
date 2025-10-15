package database

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/hayden-erickson/ai-evaluation/commandcenter"
	"github.com/hayden-erickson/ai-evaluation/constants"
	"github.com/hayden-erickson/ai-evaluation/models"
)

// Bank represents the database connection and operations
type Bank struct {
	db *sql.DB
}

// Config holds database configuration
type Config struct {
	Host     string
	Port     int
	Username string
	Password string
	Database string
}

// NewBank creates a new Bank instance with MySQL connection
func NewBank(config Config) (*Bank, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&charset=utf8mb4&collation=utf8mb4_unicode_ci",
		config.Username, config.Password, config.Host, config.Port, config.Database)
	
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}
	
	// Configure connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(300) // 5 minutes
	
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}
	
	return &Bank{db: db}, nil
}

// NewBankFromDSN creates a new Bank instance from a DSN string
func NewBankFromDSN(dsn string) (*Bank, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}
	
	// Configure connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(300) // 5 minutes
	
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}
	
	return &Bank{db: db}, nil
}

// Close closes the database connection
func (b *Bank) Close() error {
	if b.db != nil {
		return b.db.Close()
	}
	return nil
}

// Ping checks if the database connection is alive
func (b *Bank) Ping() error {
	if b.db == nil {
		return fmt.Errorf("database connection is nil")
	}
	return b.db.Ping()
}

// Stats returns database connection statistics
func (b *Bank) Stats() sql.DBStats {
	if b.db == nil {
		return sql.DBStats{}
	}
	return b.db.Stats()
}

// BeginTx starts a transaction with context
func (b *Bank) BeginTx(ctx context.Context) (*sql.Tx, error) {
	if b.db == nil {
		return nil, fmt.Errorf("database connection is nil")
	}
	return b.db.BeginTx(ctx, nil)
}

// GetBUserByID retrieves a business user by ID
func (b *Bank) GetBUserByID(BUserID int) (*models.BUser, error) {
	query := `
		SELECT u.id, u.company_uuid, GROUP_CONCAT(us.site_id) as sites
		FROM business_users u
		LEFT JOIN user_sites us ON u.id = us.user_id
		WHERE u.id = ? AND u.deleted_at IS NULL
		GROUP BY u.id, u.company_uuid
	`
	
	var user models.BUser
	var sitesStr sql.NullString
	
	err := b.db.QueryRow(query, BUserID).Scan(&user.Id, &user.CompanyUUID, &sitesStr)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no_ob_found")
		}
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}
	
	// Parse sites string into slice
	if sitesStr.Valid && sitesStr.String != "" {
		user.Sites = strings.Split(sitesStr.String, ",")
	}
	
	return &user, nil
}

// V2UnitGetById retrieves a unit by ID and site ID
func (b *Bank) V2UnitGetById(unitID int, siteID int) (*models.Unit, error) {
	query := `
		SELECT site_id, rental_state
		FROM units
		WHERE id = ? AND site_id = ? AND deleted_at IS NULL
	`
	
	var unit models.Unit
	err := b.db.QueryRow(query, unitID, siteID).Scan(&unit.SiteID, &unit.RentalState)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("unit not found")
		}
		return nil, fmt.Errorf("failed to get unit by ID: %w", err)
	}
	
	return &unit, nil
}

// GetCodesForUnits retrieves access codes for the given units and site
func (b *Bank) GetCodesForUnits(units []int, siteID int) ([]models.GateAccessCode, error) {
	if len(units) == 0 {
		return []models.GateAccessCode{}, nil
	}
	
	// Create placeholders for IN clause
	placeholders := strings.Repeat("?,", len(units))
	placeholders = placeholders[:len(placeholders)-1] // Remove trailing comma
	
	query := fmt.Sprintf(`
		SELECT access_code, unit_id, user_id, site_id, state
		FROM gate_access_codes
		WHERE unit_id IN (%s) AND site_id = ? AND deleted_at IS NULL
		ORDER BY unit_id, created_at DESC
	`, placeholders)
	
	// Build args slice
	args := make([]interface{}, len(units)+1)
	for i, unit := range units {
		args[i] = unit
	}
	args[len(units)] = siteID
	
	rows, err := b.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get access codes for units: %w", err)
	}
	defer rows.Close()
	
	var codes []models.GateAccessCode
	for rows.Next() {
		var code models.GateAccessCode
		err := rows.Scan(&code.AccessCode, &code.UnitID, &code.UserID, &code.SiteID, &code.State)
		if err != nil {
			return nil, fmt.Errorf("failed to scan access code: %w", err)
		}
		codes = append(codes, code)
	}
	
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating access codes: %w", err)
	}
	
	return codes, nil
}

// UpdateAccessCodes updates access codes in the database
func (b *Bank) UpdateAccessCodes(gacs models.GateAccessCodes, siteID int) error {
	if len(gacs) == 0 {
		return nil
	}
	
	tx, err := b.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()
	
	for _, gac := range gacs {
		// Check if access code already exists
		var existingID int
		checkQuery := `
			SELECT id FROM gate_access_codes 
			WHERE access_code = ? AND unit_id = ? AND site_id = ? AND deleted_at IS NULL
		`
		err := tx.QueryRow(checkQuery, gac.AccessCode, gac.UnitID, gac.SiteID).Scan(&existingID)
		
		if err == sql.ErrNoRows {
			// Insert new access code
			insertQuery := `
				INSERT INTO gate_access_codes (access_code, unit_id, user_id, site_id, state, created_at, updated_at)
				VALUES (?, ?, ?, ?, ?, NOW(), NOW())
			`
			_, err = tx.Exec(insertQuery, gac.AccessCode, gac.UnitID, gac.UserID, gac.SiteID, gac.State)
			if err != nil {
				return fmt.Errorf("failed to insert access code: %w", err)
			}
		} else if err != nil {
			return fmt.Errorf("failed to check existing access code: %w", err)
		} else {
			// Update existing access code
			updateQuery := `
				UPDATE gate_access_codes 
				SET state = ?, user_id = ?, updated_at = NOW()
				WHERE id = ?
			`
			_, err = tx.Exec(updateQuery, gac.State, gac.UserID, existingID)
			if err != nil {
				return fmt.Errorf("failed to update access code: %w", err)
			}
		}
	}
	
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	
	return nil
}

// NewCommandCenterClient creates a new command center client for the given site
func (b *Bank) NewCommandCenterClient(siteID int, ctx context.Context) commandcenter.AccessCodeManager {
	factory := commandcenter.NewFactory()
	return factory.(*commandcenter.Factory).NewClientWithContext(siteID, ctx)
}
