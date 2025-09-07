package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"palworld-helper/internal/core/domain"
	"palworld-helper/internal/core/ports"
	"palworld-helper/web/templates"
)

type AdminHandler struct {
	service ports.AdminService
}

func NewAdminHandler(service ports.AdminService) *AdminHandler {
	return &AdminHandler{
		service: service,
	}
}

func (h *AdminHandler) AdminPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(templates.AdminPageHTML))
}

func (h *AdminHandler) GetSchema(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	schema, err := h.service.GetDatabaseSchema()
	if err != nil {
		http.Error(w, "Failed to get schema: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(schema)
}

func (h *AdminHandler) HandleTableOperations(w http.ResponseWriter, r *http.Request) {
	// Extract table name from URL path
	path := strings.TrimPrefix(r.URL.Path, "/admin/api/table/")
	parts := strings.Split(path, "/")

	if len(parts) < 1 || parts[0] == "" {
		http.Error(w, "Table name required", http.StatusBadRequest)
		return
	}

	tableName := parts[0]

	switch r.Method {
	case "GET":
		h.getTableData(w, r, tableName)
	case "POST":
		h.insertTableData(w, r, tableName)
	case "PUT":
		if len(parts) < 2 {
			http.Error(w, "Record ID required for update", http.StatusBadRequest)
			return
		}
		id, err := strconv.Atoi(parts[1])
		if err != nil {
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}
		h.updateTableData(w, r, tableName, id)
	case "DELETE":
		if len(parts) < 2 {
			http.Error(w, "Record ID required for delete", http.StatusBadRequest)
			return
		}
		id, err := strconv.Atoi(parts[1])
		if err != nil {
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}
		h.deleteTableData(w, r, tableName, id)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *AdminHandler) getTableData(w http.ResponseWriter, r *http.Request, tableName string) {
	// Log pour débugger
	fmt.Printf("Getting data from table: %s\n", tableName)

	data, err := h.service.GetTableData(tableName)
	if err != nil {
		// Log l'erreur complète
		fmt.Printf("Get table data error: %v\n", err)
		http.Error(w, "Failed to get table data: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Log pour débugger
	fmt.Printf("Retrieved %d rows from table %s\n", len(data), tableName)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func (h *AdminHandler) insertTableData(w http.ResponseWriter, r *http.Request, tableName string) {
	var data map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Log pour débugger
	fmt.Printf("Inserting data into table %s: %+v\n", tableName, data)

	if err := h.service.InsertData(tableName, data); err != nil {
		// Log l'erreur complète
		fmt.Printf("Insert error: %v\n", err)
		http.Error(w, "Failed to insert data: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message": "Data inserted successfully"}`))
}

func (h *AdminHandler) updateTableData(w http.ResponseWriter, r *http.Request, tableName string, id int) {
	var data map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Log pour débugger
	fmt.Printf("Updating data in table %s, id %d: %+v\n", tableName, id, data)

	if err := h.service.UpdateData(tableName, id, data); err != nil {
		// Log l'erreur complète
		fmt.Printf("Update error: %v\n", err)
		http.Error(w, "Failed to update data: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message": "Data updated successfully"}`))
}

func (h *AdminHandler) deleteTableData(w http.ResponseWriter, r *http.Request, tableName string, id int) {
	if err := h.service.DeleteData(tableName, id); err != nil {
		http.Error(w, "Failed to delete data: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte(`{"message": "Data deleted successfully"}`))
}

func (h *AdminHandler) ExecuteQuery(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Query string `json:"query"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	if req.Query == "" {
		http.Error(w, "Query is required", http.StatusBadRequest)
		return
	}

	results, err := h.service.ExecuteQuery(req.Query)
	if err != nil {
		http.Error(w, "Query execution failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"results": results,
		"count":   len(results),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *AdminHandler) CreateTable(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		TableName string              `json:"table_name"`
		Columns   []domain.ColumnInfo `json:"columns"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	if req.TableName == "" {
		http.Error(w, "Table name is required", http.StatusBadRequest)
		return
	}

	if len(req.Columns) == 0 {
		http.Error(w, "At least one column is required", http.StatusBadRequest)
		return
	}

	if err := h.service.CreateTable(req.TableName, req.Columns); err != nil {
		http.Error(w, fmt.Sprintf("Failed to create table: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"message": "Table created successfully"}`))
}
