package domain

// CraftingRecipe represents a crafting recipe in the domain
type CraftingRecipe struct {
	ID          int    `json:"id" db:"id"`
	Name        string `json:"name" db:"name"`
	Category    string `json:"category" db:"category"`
	Description string `json:"description" db:"description"`
}

// Resource represents a crafting resource
type Resource struct {
	ID       int    `json:"id" db:"id"`
	Name     string `json:"name" db:"name"`
	Quantity int    `json:"quantity" db:"quantity"`
}

// RecipeResource represents the relationship between recipes and resources
type RecipeResource struct {
	ID         int `json:"id" db:"id"`
	RecipeID   int `json:"recipe_id" db:"recipe_id"`
	ResourceID int `json:"resource_id" db:"resource_id"`
	Quantity   int `json:"quantity" db:"quantity"`
}

// RecipeWithResources represents a recipe with its required resources
type RecipeWithResources struct {
	CraftingRecipe
	Resources []ResourceWithQuantity `json:"resources"`
}

// ResourceWithQuantity represents a resource with its required quantity
type ResourceWithQuantity struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Quantity int    `json:"quantity"`
}

// CraftingRequest represents a request to calculate resources
type CraftingRequest struct {
	Items []CraftingItem `json:"items"`
}

// CraftingItem represents an item in a crafting request
type CraftingItem struct {
	ID       int `json:"id"`
	Quantity int `json:"quantity"`
}

// ResourceTotal represents the total quantity needed for a resource
type ResourceTotal struct {
	Name  string `json:"name"`
	Total int    `json:"total"`
}

// TableInfo represents database table information
type TableInfo struct {
	Name    string       `json:"name"`
	Columns []ColumnInfo `json:"columns"`
}

// ColumnInfo represents database column information
type ColumnInfo struct {
	Name         string `json:"name"`
	Type         string `json:"type"`
	NotNull      bool   `json:"not_null"`
	DefaultValue string `json:"default_value"`
	PrimaryKey   bool   `json:"primary_key"`
}
