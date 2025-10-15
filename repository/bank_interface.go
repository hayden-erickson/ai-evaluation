package repository

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/hayden-erickson/ai-evaluation/models"
)

// Bank represents the database layer for business operations
type Bank struct {
	db *sql.DB
}

// NewBank creates a new Bank instance with MySQL connection
func NewBank(dsn string) (*Bank, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

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

// GetBUserByID retrieves a business user by ID
func (b *Bank) GetBUserByID(BUserID int) (*models.BUser, error) {
	query := `
		SELECT company_uuid, id, sites 
		FROM business_users 
		WHERE id = ?
	`

	var user models.BUser
	var sitesStr string

	err := b.db.QueryRow(query, BUserID).Scan(
		&user.CompanyUUID,
		&user.Id,
		&sitesStr,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("business user with ID %d not found", BUserID)
		}
		return nil, fmt.Errorf("failed to get business user: %w", err)
	}

	// Parse sites string (assuming comma-separated values)
	if sitesStr != "" {
		user.Sites = []string{sitesStr} // Simplified - in real implementation you'd parse CSV
	}

	return &user, nil
}

// V2UnitGetById retrieves a unit by ID and site ID
func (b *Bank) V2UnitGetById(unitID int, siteID int) (*models.Unit, error) {
	query := `
		SELECT site_id, rental_state 
		FROM units 
		WHERE id = ? AND site_id = ?
	`

	var unit models.Unit

	err := b.db.QueryRow(query, unitID, siteID).Scan(
		&unit.SiteID,
		&unit.RentalState,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("unit with ID %d not found in site %d", unitID, siteID)
		}
		return nil, fmt.Errorf("failed to get unit: %w", err)
	}

	return &unit, nil
}

// GetCodesForUnits retrieves access codes for given units and site
func (b *Bank) GetCodesForUnits(units []int, siteID int) ([]models.GateAccessCode, error) {
	if len(units) == 0 {
		return []models.GateAccessCode{}, nil
	}

	// Build query with IN clause
	query := `
		SELECT access_code, unit_id, user_id, site_id, state 
		FROM gate_access_codes 
		WHERE site_id = ? AND unit_id IN (`

	// Add placeholders for units
	for i := range units {
		if i > 0 {
			query += ", "
		}
		query += "?"
	}
	query += ") ORDER BY unit_id, access_code"

	// Build arguments: siteID first, then units
	args := []interface{}{siteID}
	for _, unit := range units {
		args = append(args, unit)
	}

	rows, err := b.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query access codes: %w", err)
	}
	defer rows.Close()

	var codes []models.GateAccessCode
	for rows.Next() {
		var code models.GateAccessCode
		err := rows.Scan(
			&code.AccessCode,
			&code.UnitID,
			&code.UserID,
			&code.SiteID,
			&code.State,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan access code: %w", err)
		}

		// Set default validation state
		code.IsValid = true
		codes = append(codes, code)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating access codes: %w", err)
	}

	return codes, nil
}

// UpdateAccessCodes updates access codes in the database
func (b *Bank) UpdateAccessCodes(codes []string, siteID int) error {
	if len(codes) == 0 {
		return nil
	}

	// Start a transaction for atomic updates
	tx, err := b.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Prepare the update statement
	updateQuery := `
		UPDATE gate_access_codes 
		SET state = 'active', updated_at = NOW() 
		WHERE access_code = ? AND site_id = ?
	`

	stmt, err := tx.Prepare(updateQuery)
	if err != nil {
		return fmt.Errorf("failed to prepare update statement: %w", err)
	}
	defer stmt.Close()

	// Update each access code
	for _, code := range codes {
		_, err := stmt.Exec(code, siteID)
		if err != nil {
			return fmt.Errorf("failed to update access code %s: %w", code, err)
		}
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// NewCommandCenterClient creates a new command center client for the given site
func (b *Bank) NewCommandCenterClient(siteID int, ctx context.Context) CommandCenterClient {
	return &CommandCenter{}
}
