package postgres

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/google/uuid"

	config "github.com/kiranpt03/factorysphere/devicesphere/engines/rulesengine/pkg/config"
	log "github.com/kiranpt03/factorysphere/devicesphere/engines/rulesengine/pkg/utils/loggers"
	_ "github.com/lib/pq"
)

// PostgreSQLRepository represents a PostgreSQL repository.
type PostgreSQLRepository struct {
	db *sql.DB
}

// NewPostgreSQLRepository creates a new PostgreSQL repository.
func NewPostgreSQLRepository(config *config.Config) (*PostgreSQLRepository, error) {
	// Create a connection string
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.Postgres.Host, config.Postgres.Port, config.Postgres.User, config.Postgres.Password, config.Postgres.DBName)

	// Connect to the database
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	// Ping the database to verify the connection
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	log.Debug("Connected to the database!")

	return &PostgreSQLRepository{db: db}, nil
}

// GetAll retrieves all records from the specified table.
func (r *PostgreSQLRepository) GetAll(tableName string) ([]map[string]interface{}, error) {
	log.Info("Get all devices from database")
	var records []map[string]interface{}

	query := fmt.Sprintf("SELECT * FROM %s", tableName)
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		columns := make([]interface{}, len(cols))
		columnPointers := make([]interface{}, len(cols))
		for i := range columns {
			columnPointers[i] = &columns[i]
		}

		if err := rows.Scan(columnPointers...); err != nil {
			return nil, err
		}

		record := make(map[string]interface{})
		for i, col := range cols {
			record[col] = columns[i]
		}

		records = append(records, record)
	}

	return records, nil
}

// Get retrieves a record by ID from the specified table.
func (r *PostgreSQLRepository) Get(tableName string, id string) (map[string]interface{}, error) {
	var record map[string]interface{}

	query := fmt.Sprintf("SELECT * FROM %s WHERE id = $1", tableName)
	rows, err := r.db.Query(query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	if rows.Next() {
		columns := make([]interface{}, len(cols))
		columnPointers := make([]interface{}, len(cols))
		for i := range columns {
			columnPointers[i] = &columns[i]
		}

		if err := rows.Scan(columnPointers...); err != nil {
			return nil, err
		}

		record = make(map[string]interface{})
		for i, col := range cols {
			record[col] = columns[i]
		}
	}

	return record, nil
}

// Create creates a new record in the specified table.
func (r *PostgreSQLRepository) Create(tableName string, data map[string]interface{}) (string, error) {
	// Generate a random ID
	id := uuid.New().String()

	// Add the ID to the data map
	data["id"] = id

	columns := make([]string, 0, len(data))
	values := make([]interface{}, 0, len(data))
	params := make([]string, 0, len(data))

	for k, v := range data {
		columns = append(columns, k)
		values = append(values, v)
		params = append(params, fmt.Sprintf("$%d", len(params)+1))
	}

	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s) RETURNING id",
		tableName,
		strings.Join(columns, ", "),
		strings.Join(params, ", "))

	var createdID string
	err := r.db.QueryRow(query, values...).Scan(&createdID)
	fmt.Println(err)

	return createdID, err
}

// Update updates an existing record in the specified table.
func (r *PostgreSQLRepository) Update(tableName string, id string, data map[string]interface{}) (int64, error) {
	var affectedRows int64

	columns := make([]string, 0, len(data))
	values := make([]interface{}, 0, len(data))
	for k, v := range data {
		columns = append(columns, fmt.Sprintf("%s = $%d", k, len(values)+1))
		values = append(values, v)
	}

	query := fmt.Sprintf("UPDATE %s SET %s WHERE id = $%d RETURNING id",
		tableName,
		strings.Join(columns, ", "),
		len(values)+1)

	values = append(values, id)

	err := r.db.QueryRow(query, values...).Scan(&affectedRows)

	return affectedRows, err
}

// Delete deletes a record by the specified column and value from the specified table.
func (r *PostgreSQLRepository) Delete(tableName string, column string, value interface{}) error {

	query := fmt.Sprintf("DELETE FROM %s WHERE %s = $1", tableName, column)
	_, err := r.db.Exec(query, value)

	return err
}

// CheckExist checks if a record with the specified column and value exists in the table.
func (r *PostgreSQLRepository) CheckExist(tableName string, column string, value interface{}) (bool, error) {
	var exists bool

	query := fmt.Sprintf("SELECT EXISTS (SELECT 1 FROM %s WHERE %s = $1)", tableName, column)
	err := r.db.QueryRow(query, value).Scan(&exists)

	return exists, err
}

// ExecuteQuery executes a raw SQL query and returns the results.
func (r *PostgreSQLRepository) ExecuteQuery(query string, args ...interface{}) ([]map[string]interface{}, error) {
	// log.Debug("Executing query: %s", query)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	var results []map[string]interface{}
	for rows.Next() {
		columns := make([]interface{}, len(cols))
		columnPointers := make([]interface{}, len(cols))
		for i := range columns {
			columnPointers[i] = &columns[i]
		}

		if err := rows.Scan(columnPointers...); err != nil {
			return nil, err
		}

		record := make(map[string]interface{})
		for i, col := range cols {
			record[col] = columns[i]
		}

		results = append(results, record)
	}

	return results, nil
}
