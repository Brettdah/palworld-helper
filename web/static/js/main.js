let recipes = [];
let selectedItems = [];
let categories = [];
let activeCategory = 'all'; // Ajouter une variable pour tracker le filtre actif

// Initialize the application
document.addEventListener('DOMContentLoaded', function() {
    // Restaurer le filtre actif depuis localStorage
    activeCategory = localStorage.getItem('activeCategory') || 'all';

    loadRecipes();
    renderCart();
    setupEventListeners();
});

function setupEventListeners() {
    // Search functionality
    document.getElementById('searchBox').addEventListener('input', (e) => {
        const searchTerm = e.target.value.toLowerCase();
        const filtered = recipes.filter(recipe =>
            recipe.name.toLowerCase().includes(searchTerm) ||
            recipe.description.toLowerCase().includes(searchTerm) ||
            recipe.resources.some(resource => resource.name.toLowerCase().includes(searchTerm))
        );
        renderRecipes(filtered);
    });
}

async function loadRecipes() {
    try {
        const response = await fetch('/api/recipes');
        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }
        recipes = await response.json();

        // Extract unique categories
        categories = [...new Set(recipes.map(r => r.category))];
        renderCategoryFilter();

        // Appliquer le filtre actif après chargement
        filterByCategory(activeCategory);
    } catch (error) {
        console.error('Error loading recipes:', error);
        showError('Failed to load recipes. Please check if the server is running.');
    }
}

function renderCategoryFilter() {
    const filterContainer = document.getElementById('categoryFilter');
    const allBtn = filterContainer.querySelector('[data-category="all"]');

    // Ajouter l'événement pour le bouton "All" s'il n'en a pas déjà
    if (allBtn && !allBtn.hasAttribute('data-listener-added')) {
        allBtn.addEventListener('click', () => filterByCategory('all'));
        allBtn.setAttribute('data-listener-added', 'true');
    }

    // Clear existing category buttons except "All"
    const existingBtns = filterContainer.querySelectorAll('[data-category]:not([data-category="all"])');
    existingBtns.forEach(btn => btn.remove());

    categories.forEach(category => {
        const btn = document.createElement('button');
        btn.className = 'category-btn';
        btn.textContent = category;
        btn.dataset.category = category;
        btn.addEventListener('click', () => filterByCategory(category));
        filterContainer.appendChild(btn);
    });

    // Mettre à jour l'état actif des boutons
    updateActiveButton();
}

function renderRecipes(filteredRecipes = recipes) {
    const grid = document.getElementById('recipesGrid');
    grid.innerHTML = '';

    if (filteredRecipes.length === 0) {
        grid.innerHTML = '<p style="text-align: center; color: #e94560; grid-column: 1/-1;">No recipes found.</p>';
        return;
    }

    filteredRecipes.forEach(recipe => {
        const card = document.createElement('div');
        card.className = 'recipe-card';
        card.innerHTML = `
            <div class="recipe-title">${escapeHtml(recipe.name)}</div>
            <div class="recipe-category">${escapeHtml(recipe.category)}</div>
            <div class="recipe-description">${escapeHtml(recipe.description || '')}</div>
            <ul class="resources-list">
                ${recipe.resources.map(r => `<li>${escapeHtml(r.name)}: ${r.quantity}</li>`).join('')}
            </ul>
            <div style="display: flex; align-items: center; margin-top: 10px;">
                <input type="number" class="quantity-input" min="1" value="1" data-recipe-id="${recipe.id}">
                <button onclick="addToCart(${recipe.id})" class="btn btn-primary">Add recipe</button>
            </div>
        `;
        grid.appendChild(card);
    });
}

function filterByCategory(category) {
    // Sauvegarder le filtre actif
    activeCategory = category;
    localStorage.setItem('activeCategory', category);

    updateActiveButton();

    const filtered = category === 'all' ? recipes : recipes.filter(r => r.category === category);
    renderRecipes(filtered);
}

function updateActiveButton() {
    // Mettre à jour l'état actif des boutons
    document.querySelectorAll('.category-btn').forEach(btn => {
        btn.classList.remove('active');
        if (btn.dataset.category === activeCategory) {
            btn.classList.add('active');
        }
    });
}

function addToCart(recipeId) {
    const quantityInput = document.querySelector(`[data-recipe-id="${recipeId}"]`);
    const quantity = parseInt(quantityInput.value) || 1;
    const recipe = recipes.find(r => r.id === recipeId);

    if (!recipe) {
        showError('Recipe not found');
        return;
    }

    const existingItem = selectedItems.find(item => item.id === recipeId);
    if (existingItem) {
        existingItem.quantity += quantity;
    } else {
        selectedItems.push({ id: recipeId, quantity, name: recipe.name });
    }

    // Reset quantity input
    quantityInput.value = 1;

    renderCart();
    showSuccess(`Added ${recipe.name} x${quantity} to cart`);
}

function removeFromCart(recipeId) {
    const recipe = recipes.find(r => r.id === recipeId);
    selectedItems = selectedItems.filter(item => item.id !== recipeId);
    renderCart();

    if (recipe) {
        showSuccess(`Removed ${recipe.name} from cart`);
    }
}

function renderCart() {
    const cart = document.getElementById('cart');
    if (selectedItems.length === 0) {
        cart.innerHTML = '<p style="color: #0f3460; text-align: center;">No items selected</p>';
        return;
    }

    cart.innerHTML = selectedItems.map(item => `
        <div class="cart-item">
            <span>${escapeHtml(item.name)} x${item.quantity}</span>
            <button class="remove-btn" onclick="removeFromCart(${item.id})">Remove</button>
        </div>
    `).join('');
}

async function calculateResources() {
    if (selectedItems.length === 0) {
        showError('Please select some items first!');
        return;
    }

    try {
        const response = await fetch('/api/calculate', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ items: selectedItems })
        });

        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }

        const resourceTotals = await response.json();
        renderResults(resourceTotals);
    } catch (error) {
        console.error('Error calculating resources:', error);
        showError('Failed to calculate resources. Please try again.');
    }
}

function renderResults(resourceTotals) {
    const resultsSection = document.getElementById('results');
    const totalsContainer = document.getElementById('resourceTotals');

    resultsSection.classList.remove('hidden');

    if (resourceTotals.length === 0) {
        totalsContainer.innerHTML = '<p style="color: #e94560; text-align: center;">No resources needed.</p>';
        return;
    }

    totalsContainer.innerHTML = resourceTotals.map(resource => `
        <div class="resource-total">
            <span>${escapeHtml(resource.name)}</span>
            <span>${resource.total}</span>
        </div>
    `).join('');

    // Scroll to results
    resultsSection.scrollIntoView({ behavior: 'smooth' });
}

// Utility functions
function escapeHtml(text) {
    const map = {
        '&': '&amp;',
        '<': '&lt;',
        '>': '&gt;',
        '"': '&quot;',
        "'": '&#039;'
    };
    return text ? text.replace(/[&<>"']/g, function(m) { return map[m]; }) : '';
}

function showError(message) {
    showNotification(message, 'error');
}

function showSuccess(message) {
    showNotification(message, 'success');
}

function showNotification(message, type) {
    // Create notification element
    const notification = document.createElement('div');
    notification.className = `notification notification-${type}`;
    notification.textContent = message;
    notification.style.cssText = `
        position: fixed;
        top: 20px;
        right: 20px;
        padding: 15px 20px;
        border-radius: 8px;
        color: white;
        font-weight: bold;
        z-index: 1000;
        opacity: 0;
        transform: translateY(-20px);
        transition: all 0.3s ease;
        max-width: 300px;
        word-wrap: break-word;
    `;

    if (type === 'error') {
        notification.style.background = '#dc3545';
    } else if (type === 'success') {
        notification.style.background = '#28a745';
    }

    document.body.appendChild(notification);

    // Animate in
    setTimeout(() => {
        notification.style.opacity = '1';
        notification.style.transform = 'translateY(0)';
    }, 100);

    // Remove after 3 seconds
    setTimeout(() => {
        notification.style.opacity = '0';
        notification.style.transform = 'translateY(-20px)';
        setTimeout(() => {
            if (notification.parentNode) {
                notification.parentNode.removeChild(notification);
            }
        }, 300);
    }, 3000);
}