package services

import (
	"fmt"
	"strings"

	"palworld-helper/internal/core/domain"
	"palworld-helper/internal/core/ports"
)

type adminService struct {
	repo ports.AdminRepository
}

// NewAdminService creates a new admin service
func NewAdminService(repo ports.AdminRepository) ports.AdminService {
	return &adminService{
		repo: repo,
	}
}

// GetDatabaseSchema retrieves the complete database schema
func (s *adminService) GetDatabaseSchema() ([]domain.TableInfo, error) {
	tables, err := s.repo.GetTables()
	if err != nil {
		return nil, err
	}

	var tableInfos []domain.TableInfo
	for _, tableName := range tables {
		tableInfo, err := s.repo.GetTableInfo(tableName)
		if err != nil {
			return nil, err
		}
		tableInfos = append(tableInfos, *tableInfo)
	}

	return tableInfos, nil
}

// ExecuteQuery executes a raw SQL query
func (s *adminService) ExecuteQuery(query string) ([]map[string]interface{}, error) {
	return s.repo.ExecuteQuery(query)
}

// GetTableData retrieves all data from a specific table
func (s *adminService) GetTableData(tableName string) ([]map[string]interface{}, error) {
	tables, err := s.repo.GetTables()
	if err != nil {
		return nil, fmt.Errorf("failed to get tables list: %w", err)
	}

	tableExists := false
	for _, table := range tables {
		if table == tableName {
			tableExists = true
			break
		}
	}

	if !tableExists {
		return nil, fmt.Errorf("table '%s' does not exist", tableName)
	}

	fmt.Printf("Getting data from table: %s\n", tableName)

	query := fmt.Sprintf("SELECT * FROM %s", tableName)
	data, err := s.repo.ExecuteQuery(query)
	if err != nil {
		fmt.Printf("Error executing query '%s': %v\n", query, err)
		return nil, fmt.Errorf("failed to query table %s: %w", tableName, err)
	}

	fmt.Printf("Retrieved %d rows from table %s\n", len(data), tableName)
	return data, nil
}

// CreateTable creates a new table with specified columns
func (s *adminService) CreateTable(tableName string, columns []domain.ColumnInfo) error {
	var columnDefs []string
	for _, col := range columns {
		colDef := fmt.Sprintf("%s %s", col.Name, col.Type)

		if col.PrimaryKey {
			colDef += " PRIMARY KEY"
			if strings.ToUpper(col.Type) == "INTEGER" {
				colDef += " AUTOINCREMENT"
			}
		}

		if col.NotNull && !col.PrimaryKey {
			colDef += " NOT NULL"
		}

		if col.DefaultValue != "" {
			colDef += fmt.Sprintf(" DEFAULT %s", col.DefaultValue)
		}

		columnDefs = append(columnDefs, colDef)
	}

	query := fmt.Sprintf("CREATE TABLE %s (%s)", tableName, strings.Join(columnDefs, ", "))
	return s.repo.CreateTable(query)
}

// InsertData inserts new data into a table
func (s *adminService) InsertData(tableName string, data map[string]interface{}) error {
	fmt.Printf("Attempting to insert data into table %s: %+v\n", tableName, data)

	// Get table schema to identify auto-increment columns
	tableInfo, err := s.repo.GetTableInfo(tableName)
	if err != nil {
		return fmt.Errorf("failed to get table info for %s: %w", tableName, err)
	}

	var columns []string
	var placeholders []string
	var values []interface{}

	// Filter out primary key columns that are auto-increment (INTEGER PRIMARY KEY)
	for column, value := range data {
		shouldInclude := true

		// Check if this column is an auto-increment primary key
		for _, col := range tableInfo.Columns {
			if col.Name == column && col.PrimaryKey && strings.ToUpper(col.Type) == "INTEGER" {
				// Skip auto-increment primary keys, unless explicitly provided and not empty
				if value == nil || value == "" || value == "0" {
					shouldInclude = false
					fmt.Printf("Skipping auto-increment column %s\n", column)
				}
				break
			}
		}

		if shouldInclude {
			// Nettoyer les valeurs vides
			if value == "" {
				value = nil
			}
			columns = append(columns, column)
			placeholders = append(placeholders, "?")
			values = append(values, value)
		}
	}

	if len(columns) == 0 {
		return fmt.Errorf("no valid columns to insert")
	}

	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",
		tableName,
		strings.Join(columns, ", "),
		strings.Join(placeholders, ", "))

	fmt.Printf("Executing INSERT query: %s with values: %+v\n", query, values)

	// Use a dedicated method for INSERT operations
	err = s.repo.ExecuteNonQuery(query, values...)
	if err != nil {
		fmt.Printf("Insert failed with error: %v\n", err)
		return fmt.Errorf("failed to insert data into %s: %w", tableName, err)
	}

	fmt.Printf("Insert successful\n")
	return nil
}

// UpdateData updates existing data in a table
func (s *adminService) UpdateData(tableName string, id int, data map[string]interface{}) error {
	var setParts []string
	var values []interface{}

	for column, value := range data {
		if column != "id" { // Don't update the ID column
			setParts = append(setParts, fmt.Sprintf("%s = ?", column))
			values = append(values, value)
		}
	}

	if len(setParts) == 0 {
		return fmt.Errorf("no columns to update")
	}

	values = append(values, id)
	query := fmt.Sprintf("UPDATE %s SET %s WHERE id = ?",
		tableName,
		strings.Join(setParts, ", "))

	// Use a dedicated method for UPDATE operations
	return s.repo.ExecuteNonQuery(query, values...)
}

// DeleteData deletes data from a table by ID
func (s *adminService) DeleteData(tableName string, id int) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id = ?", tableName)
	// Use a dedicated method for DELETE operations
	return s.repo.ExecuteNonQuery(query, id)
}
