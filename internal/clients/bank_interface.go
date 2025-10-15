package clients

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/hayden-erickson/ai-evaluation/internal/config"
	"github.com/hayden-erickson/ai-evaluation/internal/models"
)

// BankInterface defines the interface for bank database operations
type BankInterface interface {
	GetBUserByID(BUserID int) (*models.BUser, error)
	V2UnitGetById(unitID int, siteID int) (*models.Unit, error)
	GetCodesForUnits(units []int, siteID int) ([]models.GateAccessCode, error)
	UpdateAccessCodes(codes []string, siteID int) error
	NewCommandCenterClient(siteID int, ctx context.Context) CommandCenterInterface
	Close() error
}

// Bank is a concrete implementation of BankInterface
type Bank struct {
	db *sql.DB
}

// NewBank creates a new Bank instance with MySQL connection from configuration
func NewBank(cfg *config.Config) (*Bank, error) {
	// Validate required configuration
	if cfg.Database.Password == "" {
		return nil, fmt.Errorf("database password is required")
	}

	if cfg.Database.Name == "" {
		return nil, fmt.Errorf("database name is required")
	}

	// Create connection string
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.Database.Username, cfg.Database.Password, cfg.Database.Host, cfg.Database.Port, cfg.Database.Name)

	// Open database connection
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	// Test the connection
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

// Bank method implementations
func (b *Bank) GetBUserByID(BUserID int) (*models.BUser, error) {
	query := `
		SELECT u.id, u.company_uuid, GROUP_CONCAT(us.site_id) as sites
		FROM users u
		LEFT JOIN user_sites us ON u.id = us.user_id
		WHERE u.id = ?
		GROUP BY u.id, u.company_uuid
	`

	var user models.BUser
	var sitesStr sql.NullString

	err := b.db.QueryRow(query, BUserID).Scan(&user.Id, &user.CompanyUUID, &sitesStr)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no_ob_found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Parse sites string into slice
	if sitesStr.Valid && sitesStr.String != "" {
		user.Sites = strings.Split(sitesStr.String, ",")
	}

	return &user, nil
}

func (b *Bank) V2UnitGetById(unitID int, siteID int) (*models.Unit, error) {
	query := `
		SELECT site_id, rental_state
		FROM units
		WHERE id = ? AND site_id = ?
	`

	var unit models.Unit
	err := b.db.QueryRow(query, unitID, siteID).Scan(&unit.SiteID, &unit.RentalState)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("unit not found")
		}
		return nil, fmt.Errorf("failed to get unit: %w", err)
	}

	return &unit, nil
}

func (b *Bank) GetCodesForUnits(units []int, siteID int) ([]models.GateAccessCode, error) {
	if len(units) == 0 {
		return []models.GateAccessCode{}, nil
	}

	// Create placeholders for IN clause
	placeholders := make([]string, len(units))
	args := make([]interface{}, len(units)+1)

	for i, unitID := range units {
		placeholders[i] = "?"
		args[i] = unitID
	}
	args[len(units)] = siteID

	query := fmt.Sprintf(`
		SELECT access_code, unit_id, user_id, site_id, state
		FROM gate_access_codes
		WHERE unit_id IN (%s) AND site_id = ?
	`, strings.Join(placeholders, ","))

	rows, err := b.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get access codes: %w", err)
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

func (b *Bank) UpdateAccessCodes(codes []string, siteID int) error {
	if len(codes) == 0 {
		return nil
	}

	tx, err := b.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// For simplicity, this assumes codes are JSON strings that need to be parsed
	// In a real implementation, you'd parse the JSON and update accordingly
	query := `
		UPDATE gate_access_codes 
		SET updated_at = NOW()
		WHERE site_id = ? AND access_code IN (%s)
	`

	placeholders := make([]string, len(codes))
	args := make([]interface{}, len(codes)+1)
	args[0] = siteID

	for i, code := range codes {
		placeholders[i] = "?"
		args[i+1] = code
	}

	finalQuery := fmt.Sprintf(query, strings.Join(placeholders, ","))
	_, err = tx.Exec(finalQuery, args...)
	if err != nil {
		return fmt.Errorf("failed to update access codes: %w", err)
	}

	return tx.Commit()
}

func (b *Bank) NewCommandCenterClient(siteID int, ctx context.Context) CommandCenterInterface {
	return &CommandCenterClient{}
}
