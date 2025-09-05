package ports

import "palworld-helper/internal/core/domain"

// CraftingRepository defines the interface for crafting data operations
type CraftingRepository interface {
	GetAllRecipes() ([]domain.RecipeWithResources, error)
	GetRecipeByID(id int) (*domain.RecipeWithResources, error)
	CreateRecipe(recipe *domain.CraftingRecipe) error
	UpdateRecipe(recipe *domain.CraftingRecipe) error
	DeleteRecipe(id int) error

	GetAllResources() ([]domain.Resource, error)
	CreateResource(resource *domain.Resource) error
	UpdateResource(resource *domain.Resource) error
	DeleteResource(id int) error

	CreateRecipeResource(recipeResource *domain.RecipeResource) error
	DeleteRecipeResource(recipeID, resourceID int) error
	GetRecipeResources(recipeID int) ([]domain.ResourceWithQuantity, error)
}

// AdminRepository defines the interface for admin operations
type AdminRepository interface {
	GetTables() ([]string, error)
	GetTableInfo(tableName string) (*domain.TableInfo, error)
	ExecuteQuery(query string) ([]map[string]interface{}, error)
	ExecuteNonQuery(query string, args ...interface{}) error
	CreateTable(query string) error
	DropTable(tableName string) error
}

// CraftingService defines the interface for crafting business logic
type CraftingService interface {
	GetAllRecipes() ([]domain.RecipeWithResources, error)
	CalculateResources(request domain.CraftingRequest) ([]domain.ResourceTotal, error)
	GetCategories() ([]string, error)
}

// AdminService defines the interface for admin business logic
type AdminService interface {
	GetDatabaseSchema() ([]domain.TableInfo, error)
	ExecuteQuery(query string) ([]map[string]interface{}, error)
	GetTableData(tableName string) ([]map[string]interface{}, error)
	CreateTable(tableName string, columns []domain.ColumnInfo) error
	InsertData(tableName string, data map[string]interface{}) error
	UpdateData(tableName string, id int, data map[string]interface{}) error
	DeleteData(tableName string, id int) error
}
