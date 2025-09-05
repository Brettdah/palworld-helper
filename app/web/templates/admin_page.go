package templates

const AdminPageHTML = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Palworld Helper - Database Admin</title>
    <link rel="stylesheet" href="/static/css/main.css">
    <link rel="stylesheet" href="/static/css/admin.css">
</head>
<body>
    <div class="container">
        <header>
            <h1>🛠️ Database Administration</h1>
            <nav>
                <a href="/" class="nav-link">Crafting Helper</a>
                <a href="/admin" class="nav-link active">Database Admin</a>
            </nav>
        </header>

        <div class="admin-tabs">
            <button class="tab-btn active" onclick="showTab('schema')">Database Schema</button>
            <button class="tab-btn" onclick="showTab('query')">SQL Query</button>
            <button class="tab-btn" onclick="showTab('tables')">Manage Tables</button>
            <button class="tab-btn" onclick="showTab('create')">Create Table</button>
        </div>

        <!-- Schema Tab -->
        <div id="schema-tab" class="tab-content active">
            <h2>Database Schema</h2>
            <button onclick="loadSchema()" class="btn btn-primary">Refresh Schema</button>
            <div id="schemaContainer"></div>
        </div>

        <!-- Query Tab -->
        <div id="query-tab" class="tab-content">
            <h2>Execute SQL Query</h2>
            <div class="query-section">
                <textarea id="queryInput" placeholder="Enter your SQL query here..." rows="6"></textarea>
                <button onclick="executeQuery()" class="btn btn-primary">Execute Query</button>
                <div class="query-examples">
                    <h3>Example Queries:</h3>
                    <button onclick="setQuery('SELECT * FROM crafting_recipes')" class="btn btn-secondary">View All Recipes</button>
                    <button onclick="setQuery('SELECT * FROM resources')" class="btn btn-secondary">View All Resources</button>
                    <button onclick="setQuery('SELECT cr.name, r.name, rr.quantity FROM crafting_recipes cr JOIN recipe_resources rr ON cr.id = rr.recipe_id JOIN resources r ON rr.resource_id = r.id')" class="btn btn-secondary">View Recipe Details</button>
                </div>
            </div>
            <div id="queryResults"></div>
        </div>

        <!-- Tables Tab -->
        <div id="tables-tab" class="tab-content">
            <h2>Manage Table Data</h2>
            <div class="table-selector">
                <select id="tableSelect" onchange="loadTableData()">
                    <option value="">Select a table...</option>
                </select>
                <button onclick="addNewRecord()" class="btn btn-success" id="addRecordBtn" disabled>Add New Record</button>
            </div>
            <div id="tableDataContainer"></div>
        </div>

        <!-- Create Table Tab -->
        <div id="create-tab" class="tab-content">
            <h2>Create New Table</h2>
            <div class="create-table-form">
                <div class="form-group">
                    <label for="newTableName">Table Name:</label>
                    <input type="text" id="newTableName" placeholder="Enter table name">
                </div>
                <div class="columns-section">
                    <h3>Columns:</h3>
                    <div id="columnsContainer">
                        <div class="column-row">
                            <input type="text" placeholder="Column name" class="column-name">
                            <select class="column-type">
                                <option value="INTEGER">INTEGER</option>
                                <option value="TEXT">TEXT</option>
                                <option value="REAL">REAL</option>
                                <option value="BLOB">BLOB</option>
                            </select>
                            <label><input type="checkbox" class="column-primary"> Primary Key</label>
                            <label><input type="checkbox" class="column-notnull"> Not Null</label>
                            <input type="text" placeholder="Default value" class="column-default">
                            <button type="button" onclick="removeColumn(this)" class="btn btn-danger btn-small">Remove</button>
                        </div>
                    </div>
                    <button type="button" onclick="addColumn()" class="btn btn-secondary">Add Column</button>
                </div>
                <button onclick="createTable()" class="btn btn-primary">Create Table</button>
            </div>
        </div>

        <!-- Modal for editing records -->
        <div id="editModal" class="modal">
            <div class="modal-content">
                <span class="close" onclick="closeModal()">&times;</span>
                <h3 id="modalTitle">Edit Record</h3>
                <form id="editForm"></form>
                <div class="modal-actions">
                    <button onclick="saveRecord()" class="btn btn-primary">Save</button>
                    <button onclick="closeModal()" class="btn btn-secondary">Cancel</button>
                </div>
            </div>
        </div>
    </div>

    <script src="/static/js/admin.js"></script>
</body>
</html>`
