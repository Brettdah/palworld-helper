package services

import (
	"sort"

	"palworld-helper/internal/core/domain"
	"palworld-helper/internal/core/ports"
)

type craftingService struct {
	repo ports.CraftingRepository
}

// NewCraftingService creates a new crafting service
func NewCraftingService(repo ports.CraftingRepository) ports.CraftingService {
	return &craftingService{
		repo: repo,
	}
}

// GetAllRecipes retrieves all crafting recipes
func (s *craftingService) GetAllRecipes() ([]domain.RecipeWithResources, error) {
	return s.repo.GetAllRecipes()
}

// CalculateResources calculates the total resources needed for crafting items
func (s *craftingService) CalculateResources(request domain.CraftingRequest) ([]domain.ResourceTotal, error) {
	resourceTotals := make(map[string]int)

	for _, item := range request.Items {
		recipe, err := s.repo.GetRecipeByID(item.ID)
		if err != nil {
			return nil, err
		}

		if recipe != nil {
			for _, resource := range recipe.Resources {
				resourceTotals[resource.Name] += resource.Quantity * item.Quantity
			}
		}
	}

	// Convert map to slice for consistent output
	var results []domain.ResourceTotal
	for name, total := range resourceTotals {
		results = append(results, domain.ResourceTotal{Name: name, Total: total})
	}

	// Sort by resource name for consistent output
	sort.Slice(results, func(i, j int) bool {
		return results[i].Name < results[j].Name
	})

	return results, nil
}

// GetCategories retrieves all unique categories
func (s *craftingService) GetCategories() ([]string, error) {
	recipes, err := s.repo.GetAllRecipes()
	if err != nil {
		return nil, err
	}

	categoryMap := make(map[string]bool)
	for _, recipe := range recipes {
		categoryMap[recipe.Category] = true
	}

	categories := make([]string, 0, len(categoryMap))
	for category := range categoryMap {
		categories = append(categories, category)
	}

	sort.Strings(categories)
	return categories, nil
}
