package templates

const MainPageHTML = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Palworld Helper</title>
    <link rel="stylesheet" href="/static/css/main.css">
</head>
<body>
    <div class="container">
        <header>
            <h1>🎮 Palworld Helper</h1>
            <nav>
                <a href="/" class="nav-link active">Crafting Helper</a>
                <a href="/admin" class="nav-link">Database Admin</a>
            </nav>
        </header>

        <div class="search-section">
            <input type="text" id="searchBox" class="search-box" placeholder="Search for items...">
            <div class="category-filter" id="categoryFilter">
                <button class="category-btn active" data-category="all">All</button>
            </div>
        </div>

        <div class="recipes-grid" id="recipesGrid"></div>

        <div class="cart-section">
            <h2>Selected Items</h2>
            <div id="cart"></div>
            <button class="calculate-btn" onclick="calculateResources()">Calculate Total Resources</button>
        </div>

        <div class="results-section hidden" id="results">
            <h2>Total Resources Needed</h2>
            <div id="resourceTotals"></div>
        </div>
    </div>

    <script src="/static/js/main.js"></script>
</body>
</html>`
