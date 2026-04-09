package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/thetsGit/spend-wise-be/internal/models"
)

func (db *DB) InsertEmail(e models.RawEmail) (models.Email, error) {
	var result models.Email
	err := db.Pool.QueryRow(
		context.Background(),
		`INSERT INTO email (sender, recipient, subject, body, date)
		 VALUES ($1, $2, $3, $4, $5)
		 ON CONFLICT (sender, recipient, subject, date) DO NOTHING RETURNING *`,
		e.Sender, e.Recipient, e.Subject, e.Body, e.Date,
	).Scan(&result.ID,
		&result.Sender,
		&result.Recipient,
		&result.Subject,
		&result.Body,
		&result.Date,
		&result.Status,
		&result.CreatedAt)

	if err != nil {
		if err.Error() == "no rows in result set" {
			return models.Email{}, nil // Duplicated, not an error (should skip)
		}
		return models.Email{}, err // Real error
	}

	return result, err
}

func (db *DB) UpdateEmailStatus(id int, status string) (string, error) {
	var updatedStatus string
	err := db.Pool.QueryRow(
		context.Background(),
		`UPDATE email SET status = $1 WHERE id = $2 RETURNING status`,
		status, id,
	).Scan(&updatedStatus)
	return updatedStatus, err
}

func (db *DB) InsertSpending(s models.Spending) (models.Spending, error) {
	var result models.Spending
	err := db.Pool.QueryRow(
		context.Background(),
		`INSERT INTO spending (email_id, merchant, amount, currency, category, transaction_date, ai_confidence, confidence)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING *`,
		s.EmailID, s.Merchant, s.Amount, s.Currency, s.Category, s.TransactionDate, s.AIConfidence, s.Confidence,
	).Scan(
		&result.ID,
		&result.Merchant,
		&result.Amount, &result.Currency,
		&result.Category,
		&result.TransactionDate,
		&result.AIConfidence,
		&result.Confidence,
		&result.CreatedAt,
		&result.EmailID,
	)
	return result, err
}

func (db *DB) InsertSaaSDiscovery(s models.SaaSDiscovery) (models.SaaSDiscovery, error) {
	var result models.SaaSDiscovery
	err := db.Pool.QueryRow(
		context.Background(),
		`INSERT INTO saas_discovery (email_id, product_name, signal_type, billing_cycle, estimated_cost, currency, ai_confidence, confidence)
         VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
         RETURNING *`,
		s.EmailID, s.ProductName, s.SignalType, s.BillingCycle, s.EstimatedCost, s.Currency, s.AIConfidence, s.Confidence,
	).Scan(
		&result.ID,
		&result.ProductName,
		&result.SignalType,
		&result.BillingCycle,
		&result.EstimatedCost,
		&result.Currency,
		&result.AIConfidence,
		&result.Confidence,
		&result.CreatedAt,
		&result.EmailID,
	)
	return result, err
}

func (db *DB) GetSpending(filter models.SpendingFilter) ([]models.Spending, error) {
	query := `SELECT id, email_id, merchant, amount, currency, category, transaction_date, ai_confidence, confidence, created_at
		FROM spending WHERE 1=1`
	args := []any{}
	argIdx := 1

	if filter.Category != "" {
		query += fmt.Sprintf(" AND category = $%d", argIdx)
		args = append(args, filter.Category)
		argIdx++
	}
	if filter.StartDate != "" {
		query += fmt.Sprintf(" AND transaction_date >= $%d", argIdx)
		args = append(args, filter.StartDate)
		argIdx++
	}
	if filter.EndDate != "" {
		query += fmt.Sprintf(" AND transaction_date <= $%d", argIdx)
		args = append(args, filter.EndDate)
		argIdx++
	}

	query += " ORDER BY transaction_date DESC"

	rows, err := db.Pool.Query(context.Background(), query, args...)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, pgx.RowToStructByName[models.Spending])
}

func (db *DB) GetSpendingSummary() (models.SpendingSummary, error) {
	var categories []models.CategorySummary
	rows, err := db.Pool.Query(
		context.Background(),
		`SELECT category, COALESCE(SUM(amount), 0) as total_spend, COUNT(*) as count
		 FROM spending
		 GROUP BY category
		 ORDER BY total_spend DESC`,
	)
	if err != nil {
		return models.SpendingSummary{}, err
	}
	defer rows.Close()

	var totalSpend float64
	var totalCount int
	for rows.Next() {
		var c models.CategorySummary
		if err := rows.Scan(&c.Category, &c.TotalSpend, &c.TotalCount); err != nil {
			return models.SpendingSummary{}, err
		}
		totalSpend += c.TotalSpend
		totalCount += c.TotalCount
		categories = append(categories, c)
	}

	return models.SpendingSummary{
		TotalSpend: totalSpend,
		TotalCount: totalCount,
		ByCategory: categories,
	}, nil
}

func (db *DB) GetSaaSDiscoveries(filter models.SaaSDiscoveryFilter) ([]models.SaaSDiscovery, error) {

	query := `SELECT id, email_id, product_name, signal_type, billing_cycle, estimated_cost, currency, ai_confidence, confidence, created_at
	 FROM saas_discovery WHERE 1=1`

	args := []any{}
	argIdx := 1

	if filter.ProductName != "" {
		query += fmt.Sprintf(" AND product_name = $%d", argIdx)
		args = append(args, filter.ProductName)
		argIdx++
	}
	if filter.SignalType != "" {
		query += fmt.Sprintf(" AND signal_type >= $%d", argIdx)
		args = append(args, filter.SignalType)
		argIdx++
	}

	query += " ORDER BY product_name, created_at DESC"

	rows, err := db.Pool.Query(context.Background(), query, args...)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, pgx.RowToStructByName[models.SaaSDiscovery])
}

func (db *DB) GetSaaSDiscoverySummary() (models.SaaSSummary, error) {
	var summary models.SaaSSummary
	err := db.Pool.QueryRow(
		context.Background(),
		`SELECT
			COALESCE(SUM(
				CASE
					WHEN billing_cycle = 'annual' THEN estimated_cost / 12
					ELSE estimated_cost
				END
			), 0) as total_monthly_spend,
			COUNT(DISTINCT product_name) as total_tools_found
		 FROM saas_discovery`,
	).Scan(&summary.TotalMonthlySpend, &summary.TotalToolsFound)
	return summary, err
}
