package database

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	"palworld-helper/internal/core/domain"

	_ "modernc.org/sqlite"
)

type SQLiteDB struct {
	db *sql.DB
}

// NewSQLiteDB creates a new SQLite database connection and initializes the schema
func NewSQLiteDB(dbPath string) (*SQLiteDB, error) {
	// Create data directory if it doesn't exist
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create directory: %v", err)
	}

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	sqlite := &SQLiteDB{db: db}

	// Initialize schema and sample data
	if err := sqlite.initSchema(); err != nil {
		return nil, fmt.Errorf("failed to initialize schema: %v", err)
	}

	return sqlite, nil
}

// Close closes the database connection
func (s *SQLiteDB) Close() error {
	return s.db.Close()
}

// initSchema creates the initial database schema and populates with sample data
func (s *SQLiteDB) initSchema() error {
	// Create tables
	queries := []string{
		`CREATE TABLE IF NOT EXISTS resources (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL UNIQUE
		)`,
		`CREATE TABLE IF NOT EXISTS crafting_recipes (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			category_id TEXT NOT NULL,
			FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE CASCADE
		)`,
		`CREATE TABLE IF NOT EXISTS recipe_resources (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			recipe_id INTEGER NOT NULL,
			resource_id INTEGER NOT NULL,
			quantity INTEGER NOT NULL,
			FOREIGN KEY (recipe_id) REFERENCES crafting_recipes(id) ON DELETE CASCADE,
			FOREIGN KEY (resource_id) REFERENCES resources(id) ON DELETE CASCADE,
			UNIQUE(recipe_id, resource_id)
		)`,
		`CREATE TABLE IF NOT EXISTS categories (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL UNIQUE
		)`,
		`CREATE TABLE IF NOT EXISTS technologies (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL UNIQUE,
			level INTEGER NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS technology_recipes (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			technology_id TEXT NOT NULL,
			recipe_id INTEGER NOT NULL,
			FOREIGN KEY (recipe_id) REFERENCES crafting_recipes(id) ON DELETE CASCADE,
			FOREIGN KEY (technology_id) REFERENCES technologies(id) ON DELETE CASCADE,
			UNIQUE(technology_id, recipe_id)
		)`,
		`CREATE TABLE IF NOT EXISTS inventory (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			resource_id INTEGER NOT NULL,
			quantity INTEGER NOT NULL,
			FOREIGN KEY (resource_id) REFERENCES resources(id) ON DELETE CASCADE,
		)`,
	}

	for _, query := range queries {
		if _, err := s.db.Exec(query); err != nil {
			return fmt.Errorf("failed to execute query: %v", err)
		}
	}

	// Check if we need to populate sample data
	var count int
	err := s.db.QueryRow("SELECT COUNT(*) FROM crafting_recipes").Scan(&count)
	if err != nil {
		return err
	}

	if count == 0 {
		return s.populateSampleData()
	}

	return nil
}

// populateSampleData inserts sample crafting data
func (s *SQLiteDB) populateSampleData() error {
	// Insert resources
	resources := []string{
		"Wood", "Stone", "Cloth", "Paldium Fragment", "Metal Ore", "Coal", "Fiber",
	}

	for _, resource := range resources {
		_, err := s.db.Exec("INSERT OR IGNORE INTO resources (name) VALUES (?)", resource)
		if err != nil {
			return err
		}
	}

	// Insert recipes
	recipes := []struct {
		name, category, description string
		resources                   map[string]int
	}{
		{"Wooden Club", "Weapons", "A simple wooden weapon for early combat", map[string]int{"Wood": 5, "Stone": 2}},
		{"Stone Pickaxe", "Tools", "Essential tool for mining stone and ore", map[string]int{"Wood": 5, "Stone": 5}},
		{"Stone Axe", "Tools", "Efficient tool for cutting trees", map[string]int{"Wood": 5, "Stone": 5}},
		{"Campfire", "Structures", "Cook food and provide warmth", map[string]int{"Wood": 10, "Stone": 5}},
		{"Wooden Chest", "Storage", "Basic storage container", map[string]int{"Wood": 15, "Stone": 5}},
		{"Cloth Outfit", "Armor", "Basic protection from elements", map[string]int{"Cloth": 10}},
		{"Pal Sphere", "Pal Items", "Capture wild Pals", map[string]int{"Paldium Fragment": 3, "Wood": 3, "Stone": 3}},
		{"Workbench", "Structures", "Craft advanced items", map[string]int{"Wood": 20, "Stone": 10}},
		{"Wooden Foundation", "Building", "Foundation for wooden structures", map[string]int{"Wood": 8}},
		{"Wooden Wall", "Building", "Wall for wooden structures", map[string]int{"Wood": 6}},
	}

	for _, recipe := range recipes {
		// Insert recipe
		result, err := s.db.Exec(
			"INSERT INTO crafting_recipes (name, category, description) VALUES (?, ?, ?)",
			recipe.name, recipe.category, recipe.description,
		)
		if err != nil {
			return err
		}

		recipeID, _ := result.LastInsertId()

		// Insert recipe resources
		for resourceName, quantity := range recipe.resources {
			// Get resource ID
			var resourceID int
			err := s.db.QueryRow("SELECT id FROM resources WHERE name = ?", resourceName).Scan(&resourceID)
			if err != nil {
				return err
			}

			// Insert recipe resource relationship
			_, err = s.db.Exec(
				"INSERT INTO recipe_resources (recipe_id, resource_id, quantity) VALUES (?, ?, ?)",
				recipeID, resourceID, quantity,
			)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// GetAllRecipes retrieves all recipes with their resources
func (s *SQLiteDB) GetAllRecipes() ([]domain.RecipeWithResources, error) {
	query := `
		SELECT cr.id, cr.name, cr.category, cr.description
		FROM crafting_recipes cr
		ORDER BY cr.name
	`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var recipes []domain.RecipeWithResources
	for rows.Next() {
		var recipe domain.RecipeWithResources
		err := rows.Scan(&recipe.ID, &recipe.Name, &recipe.Category, &recipe.Description)
		if err != nil {
			return nil, err
		}

		// Get resources for this recipe
		resources, err := s.GetRecipeResources(recipe.ID)
		if err != nil {
			return nil, err
		}
		recipe.Resources = resources

		recipes = append(recipes, recipe)
	}

	return recipes, nil
}

// GetRecipeByID retrieves a specific recipe by ID
func (s *SQLiteDB) GetRecipeByID(id int) (*domain.RecipeWithResources, error) {
	query := `
		SELECT cr.id, cr.name, cr.category, cr.description
		FROM crafting_recipes cr
		WHERE cr.id = ?
	`

	var recipe domain.RecipeWithResources
	err := s.db.QueryRow(query, id).Scan(&recipe.ID, &recipe.Name, &recipe.Category, &recipe.Description)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	// Get resources for this recipe
	resources, err := s.GetRecipeResources(recipe.ID)
	if err != nil {
		return nil, err
	}
	recipe.Resources = resources

	return &recipe, nil
}

// GetRecipeResources retrieves resources for a specific recipe
func (s *SQLiteDB) GetRecipeResources(recipeID int) ([]domain.ResourceWithQuantity, error) {
	query := `
		SELECT r.id, r.name, rr.quantity
		FROM resources r
		JOIN recipe_resources rr ON r.id = rr.resource_id
		WHERE rr.recipe_id = ?
		ORDER BY r.name
	`

	rows, err := s.db.Query(query, recipeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var resources []domain.ResourceWithQuantity
	for rows.Next() {
		var resource domain.ResourceWithQuantity
		err := rows.Scan(&resource.ID, &resource.Name, &resource.Quantity)
		if err != nil {
			return nil, err
		}
		resources = append(resources, resource)
	}

	return resources, nil
}

// Implement remaining methods for CraftingRepository interface...
func (s *SQLiteDB) CreateRecipe(recipe *domain.CraftingRecipe) error {
	result, err := s.db.Exec(
		"INSERT INTO crafting_recipes (name, category, description) VALUES (?, ?, ?)",
		recipe.Name, recipe.Category, recipe.Description,
	)
	if err != nil {
		return err
	}

	id, _ := result.LastInsertId()
	recipe.ID = int(id)
	return nil
}

func (s *SQLiteDB) UpdateRecipe(recipe *domain.CraftingRecipe) error {
	_, err := s.db.Exec(
		"UPDATE crafting_recipes SET name = ?, category = ?, description = ? WHERE id = ?",
		recipe.Name, recipe.Category, recipe.Description, recipe.ID,
	)
	return err
}

func (s *SQLiteDB) DeleteRecipe(id int) error {
	_, err := s.db.Exec("DELETE FROM crafting_recipes WHERE id = ?", id)
	return err
}

func (s *SQLiteDB) GetAllResources() ([]domain.Resource, error) {
	rows, err := s.db.Query("SELECT id, name FROM resources ORDER BY name")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var resources []domain.Resource
	for rows.Next() {
		var resource domain.Resource
		err := rows.Scan(&resource.ID, &resource.Name)
		if err != nil {
			return nil, err
		}
		resources = append(resources, resource)
	}

	return resources, nil
}

func (s *SQLiteDB) CreateResource(resource *domain.Resource) error {
	result, err := s.db.Exec("INSERT INTO resources (name) VALUES (?)", resource.Name)
	if err != nil {
		return err
	}

	id, _ := result.LastInsertId()
	resource.ID = int(id)
	return nil
}

func (s *SQLiteDB) UpdateResource(resource *domain.Resource) error {
	_, err := s.db.Exec("UPDATE resources SET name = ? WHERE id = ?", resource.Name, resource.ID)
	return err
}

func (s *SQLiteDB) DeleteResource(id int) error {
	_, err := s.db.Exec("DELETE FROM resources WHERE id = ?", id)
	return err
}

func (s *SQLiteDB) CreateRecipeResource(recipeResource *domain.RecipeResource) error {
	result, err := s.db.Exec(
		"INSERT INTO recipe_resources (recipe_id, resource_id, quantity) VALUES (?, ?, ?)",
		recipeResource.RecipeID, recipeResource.ResourceID, recipeResource.Quantity,
	)
	if err != nil {
		return err
	}

	id, _ := result.LastInsertId()
	recipeResource.ID = int(id)
	return nil
}

func (s *SQLiteDB) DeleteRecipeResource(recipeID, resourceID int) error {
	_, err := s.db.Exec(
		"DELETE FROM recipe_resources WHERE recipe_id = ? AND resource_id = ?",
		recipeID, resourceID,
	)
	return err
}

// Admin interface methods
func (s *SQLiteDB) GetTables() ([]string, error) {
	rows, err := s.db.Query("SELECT name FROM sqlite_master WHERE type='table' AND name NOT LIKE 'sqlite_%'")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var table string
		if err := rows.Scan(&table); err != nil {
			return nil, err
		}
		tables = append(tables, table)
	}

	return tables, nil
}

func (s *SQLiteDB) GetTableInfo(tableName string) (*domain.TableInfo, error) {
	query := fmt.Sprintf("PRAGMA table_info(%s)", tableName)
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var columns []domain.ColumnInfo
	for rows.Next() {
		var cid int
		var name, dataType, defaultValue sql.NullString
		var notNull, pk bool

		err := rows.Scan(&cid, &name, &dataType, &notNull, &defaultValue, &pk)
		if err != nil {
			return nil, err
		}

		column := domain.ColumnInfo{
			Name:         name.String,
			Type:         dataType.String,
			NotNull:      notNull,
			DefaultValue: defaultValue.String,
			PrimaryKey:   pk,
		}
		columns = append(columns, column)
	}

	return &domain.TableInfo{
		Name:    tableName,
		Columns: columns,
	}, nil
}

func (s *SQLiteDB) ExecuteQuery(query string) ([]map[string]interface{}, error) {
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	var results []map[string]interface{}
	for rows.Next() {
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))

		for i := range values {
			valuePtrs[i] = &values[i]
		}

		if err := rows.Scan(valuePtrs...); err != nil {
			return nil, err
		}

		row := make(map[string]interface{})
		for i, col := range columns {
			val := values[i]
			if b, ok := val.([]byte); ok {
				row[col] = string(b)
			} else {
				row[col] = val
			}
		}

		results = append(results, row)
	}

	return results, nil
}

func (s *SQLiteDB) ExecuteNonQuery(query string, args ...interface{}) error {
	result, err := s.db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to execute statement: %w", err)
	}

	if rowsAffected, err := result.RowsAffected(); err == nil {
		fmt.Printf("Query executed successfully, rows affected: %d\n", rowsAffected)
	}
	return nil
}

func (s *SQLiteDB) CreateTable(query string) error {
	_, err := s.db.Exec(query)
	return err
}

func (s *SQLiteDB) DropTable(tableName string) error {
	query := fmt.Sprintf("DROP TABLE IF EXISTS %s", tableName)
	_, err := s.db.Exec(query)
	return err
}
