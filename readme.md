# Palworld Helper

> ⚠️👷 this is a work in progress ⚠️

A dark-themed web application to help you calculate and track crafting resources needed for items in Palworld. Built with Go and designed to run in Docker without requiring Go installation on your development machine.

## Features

- 🌙 **Dark Theme**: Easy on the eyes gaming aesthetic
- 🔍 **Search & Filter**: Find items by name, description, or required resources
- 📦 **Shopping Cart**: Add multiple items with quantities
- 🧮 **Resource Calculator**: Calculate total resources needed for all selected items
- 📱 **Responsive Design**: Works on desktop and mobile devices
- 🐳 **Docker Ready**: No need to install Go locally
- TODO **technology Tree**: add technology tree with check box to enable specific recipes that need to be unlocked first

## Quick Start

### Prerequisites

- Docker and Docker Compose installed on your system

### Running the Application

1. **Clone or download** all the files to a directory
2. **Navigate** to the project directory in your terminal
3. **Run** the application using Docker Compose:

```bash
# Build and start the application
docker-compose up -d

# Or use the Makefile for easier commands
make run
```

4. **Access** the application at `http://localhost:8080`

## Available Commands

If you have `make` installed, you can use these convenient commands:

```bash
make help     # Show all available commands
make build    # Build the Docker container
make run      # Build and run the application
make stop     # Stop the running container
make logs     # Show application logs
make clean    # Remove containers and images
make rebuild  # Clean rebuild and run
make shell    # Open shell in running container
```

### Manual Docker Commands

If you prefer using Docker commands directly:

```bash
# Build the image
docker-compose build

# Run the application
docker-compose up -d

# View logs
docker-compose logs -f

# Stop the application
docker-compose down

# Remove everything
docker-compose down -v --rmi all
```

## Application Structure

```tree
palworld-crafting/
├── main.go              # Main Go application
├── go.mod              # Go module file
├── Dockerfile          # Docker build instructions
├── docker-compose.yml  # Docker Compose configuration
├── Makefile           # Build automation (optional)
└── README.md          # This file
```

## How to Use

1. **Browse Items**: Scroll through the available crafting recipes
2. **Search**: Use the search box to find specific items
3. **Filter**: Click category buttons to filter by item type
4. **Add recipe**: Select quantity and add items to your crafting list
5. **Calculate**: Click "Calculate Total Resources" to see what you need to gather

## Adding New Items

To add new crafting recipes, edit the `craftingRecipes` slice in `main.go`:

```go
{
    ID:          11,
    Name:        "Your New Item",
    Category:    "Category",
    Description: "Item description",
    Resources: []Resource{
        {Name: "Resource Name", Quantity: 5},
        {Name: "Another Resource", Quantity: 3},
    },
},
```

After making changes, rebuild the container:

```bash
make rebuild
# or
docker-compose down && docker-compose up --build -d
```

## Customization

### Adding More Categories

The application automatically detects categories from the recipe data. Just add recipes with new category names.

### Modifying the Theme

The CSS is embedded in the HTML template within `main.go`. Look for the `<style>` section to customize colors and appearance.

### Changing the Port

Modify the port in both `docker-compose.yml` and the Go code if needed:

- `docker-compose.yml`: Change `"8080:8080"` to `"YOUR_PORT:8080"`
- `main.go`: Change `":8080"` in `http.ListenAndServe(":8080", nil)`

## Development

The application is completely self-contained in a single Go file with embedded HTML, CSS, and JavaScript. This makes it easy to deploy and modify without dealing with separate frontend build processes.

### Architecture

- **Backend**: Go HTTP server with JSON API endpoints
- **Frontend**: Vanilla HTML/CSS/JavaScript (no frameworks)
- **Storage**: In-memory (no database required)
- **Container**: Multi-stage Docker build for minimal image size

## Contributing

Feel free to fork this project and add more Palworld items, improve the UI, or add new features like:

- Item icons
- Crafting station requirements
- Technology tree prerequisites
- Import/export crafting lists
- Multiple save profiles

## License

This project is provided as-is for personal use. Palworld is a trademark of Pocket Pair, Inc.
