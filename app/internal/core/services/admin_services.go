package services

import (
	"fmt"
	"strings"

	"palworld-helper/app/internal/core/domain"
	"palworld-helper/app/internal/core/ports"
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
	query := fmt.Sprintf("SELECT * FROM %s", tableName)
	return s.repo.ExecuteQuery(query)
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
	var columns []string
	var placeholders []string
	var values []interface{}

	for column, value := range data {
		columns = append(columns, column)
		placeholders = append(placeholders, "?")
		values = append(values, value)
	}

	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",
		tableName,
		strings.Join(columns, ", "),
		strings.Join(placeholders, ", "))

	_, err := s.repo.ExecuteQuery(query)
	return err
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

	values = append(values, id)
	query := fmt.Sprintf("UPDATE %s SET %s WHERE id = ?",
		tableName,
		strings.Join(setParts, ", "))

	_, err := s.repo.ExecuteQuery(query)
	return err
}

// DeleteData deletes data from a table by ID
func (s *adminService) DeleteData(tableName string, id int) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id = ?", tableName)
	_, err := s.repo.ExecuteQuery(query)
	return err
}
