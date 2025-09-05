let currentTable = '';
let currentEditingRecord = null;

// Initialize admin interface
document.addEventListener('DOMContentLoaded', function() {
    showTab('schema');
    loadSchema();
    loadTableList();
});

function showTab(tabName) {
    // Hide all tabs
    document.querySelectorAll('.tab-content').forEach(tab => {
        tab.classList.remove('active');
    });
    document.querySelectorAll('.tab-btn').forEach(btn => {
        btn.classList.remove('active');
    });

    // Show selected tab
    document.getElementById(`${tabName}-tab`).classList.add('active');
    document.querySelector(`[onclick="showTab('${tabName}')"]`).classList.add('active');
}

async function loadSchema() {
    try {
        const response = await fetch('/admin/api/schema');
        if (!response.ok) throw new Error(`HTTP error! status: ${response.status}`);

        const schema = await response.json();
        renderSchema(schema);
    } catch (error) {
        console.error('Error loading schema:', error);
        showNotification('Failed to load database schema', 'error');
    }
}

function renderSchema(schema) {
    const container = document.getElementById('schemaContainer');

    if (schema.length === 0) {
        container.innerHTML = '<p style="color: #e94560; text-align: center;">No tables found in database.</p>';
        return;
    }

    container.innerHTML = schema.map(table => `
        <div class="table-schema">
            <div class="table-name">${escapeHtml(table.name)} (${table.columns.length} columns)</div>
            <div class="columns-list">
                ${table.columns.map(col => `
                    <div class="column-item">
                        <span class="column-name">${escapeHtml(col.name)}</span>
                        <span class="column-type">${escapeHtml(col.type)}</span>
                        <div class="column-attributes">
                            ${col.primary_key ? '<span class="column-attribute">PK</span>' : ''}
                            ${col.not_null ? '<span class="column-attribute">NOT NULL</span>' : ''}
                            ${col.default_value ? `<span class="column-attribute">DEFAULT: ${escapeHtml(col.default_value)}</span>` : ''}
                        </div>
                    </div>
                `).join('')}
            </div>
        </div>
    `).join('');
}

async function loadTableList() {
    try {
        const response = await fetch('/admin/api/schema');
        if (!response.ok) throw new Error(`HTTP error! status: ${response.status}`);

        const schema = await response.json();
        const tableSelect = document.getElementById('tableSelect');

        tableSelect.innerHTML = '<option value="">Select a table...</option>';
        schema.forEach(table => {
            const option = document.createElement('option');
            option.value = table.name;
            option.textContent = table.name;
            tableSelect.appendChild(option);
        });
    } catch (error) {
        console.error('Error loading table list:', error);
        showNotification('Failed to load table list', 'error');
    }
}

async function loadTableData() {
    const tableSelect = document.getElementById('tableSelect');
    const tableName = tableSelect.value;
    const addButton = document.getElementById('addRecordBtn');

    if (!tableName) {
        document.getElementById('tableDataContainer').innerHTML = '';
        addButton.disabled = true;
        return;
    }

    currentTable = tableName;
    addButton.disabled = false;

    try {
        const response = await fetch(`/admin/api/table/${tableName}`);
        if (!response.ok) throw new Error(`HTTP error! status: ${response.status}`);

        const data = await response.json();
        renderTableData(data);
    } catch (error) {
        console.error('Error loading table data:', error);
        showNotification(`Failed to load data for table: ${tableName}`, 'error');
    }
}

function renderTableData(data) {
    const container = document.getElementById('tableDataContainer');

    if (data.length === 0) {
        container.innerHTML = '<p style="color: #e94560; text-align: center;">No data found in this table.</p>';
        return;
    }

    const columns = Object.keys(data[0]);

    container.innerHTML = `
        <div class="data-table">
            <table>
                <thead>
                    <tr>
                        ${columns.map(col => `<th>${escapeHtml(col)}</th>`).join('')}
                        <th>Actions</th>
                    </tr>
                </thead>
                <tbody>
                    ${data.map(row => `
                        <tr>
                            ${columns.map(col => `<td>${escapeHtml(String(row[col] || ''))}</td>`).join('')}
                            <td>
                                <div class="action-buttons">
                                    <button onclick="editRecord(${row.id})" class="btn btn-primary btn-small">Edit</button>
                                    <button onclick="deleteRecord(${row.id})" class="btn btn-danger btn-small">Delete</button>
                                </div>
                            </td>
                        </tr>
                    `).join('')}
                </tbody>
            </table>
        </div>
    `;
}

async function executeQuery() {
    const query = document.getElementById('queryInput').value.trim();

    if (!query) {
        showNotification('Please enter a SQL query', 'error');
        return;
    }

    try {
        const response = await fetch('/admin/api/query', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ query })
        });

        if (!response.ok) {
            const errorText = await response.text();
            throw new Error(errorText);
        }

        const result = await response.json();
        renderQueryResults(result);
        showNotification('Query executed successfully', 'success');
    } catch (error) {
        console.error('Error executing query:', error);
        showNotification(`Query failed: ${error.message}`, 'error');
    }
}

function renderQueryResults(result) {
    const container = document.getElementById('queryResults');

    container.innerHTML = `
        <div class="query-result">
            <div class="result-info">
                Query returned ${result.count} row(s)
            </div>
            ${result.results.length > 0 ? renderQueryTable(result.results) : '<p>No results to display.</p>'}
        </div>
    `;
}

function renderQueryTable(data) {
    if (data.length === 0) return '<p>No data returned.</p>';

    const columns = Object.keys(data[0]);

    return `
        <div class="data-table">
            <table>
                <thead>
                    <tr>
                        ${columns.map(col => `<th>${escapeHtml(col)}</th>`).join('')}
                    </tr>
                </thead>
                <tbody>
                    ${data.map(row => `
                        <tr>
                            ${columns.map(col => `<td>${escapeHtml(String(row[col] || ''))}</td>`).join('')}
                        </tr>
                    `).join('')}
                </tbody>
            </table>
        </div>
    `;
}

function setQuery(query) {
    document.getElementById('queryInput').value = query;
}

function addColumn() {
    const container = document.getElementById('columnsContainer');
    const columnRow = document.createElement('div');
    columnRow.className = 'column-row';
    columnRow.innerHTML = `
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
    `;
    container.appendChild(columnRow);
}

function removeColumn(button) {
    button.parentElement.remove();
}

// FONCTIONS MANQUANTES - AJOUTÉES CI-DESSOUS

/**
 * Crée une nouvelle table
 */
async function createTable() {
    const tableName = document.getElementById('newTableName').value.trim();

    if (!tableName) {
        showNotification('Please enter a table name', 'error');
        return;
    }

    const columnRows = document.querySelectorAll('#columnsContainer .column-row');
    if (columnRows.length === 0) {
        showNotification('Please add at least one column', 'error');
        return;
    }

    const columns = [];
    let hasError = false;

    columnRows.forEach(row => {
        const name = row.querySelector('.column-name').value.trim();
        const type = row.querySelector('.column-type').value;
        const isPrimary = row.querySelector('.column-primary').checked;
        const isNotNull = row.querySelector('.column-notnull').checked;
        const defaultValue = row.querySelector('.column-default').value.trim();

        if (!name) {
            showNotification('All columns must have a name', 'error');
            hasError = true;
            return;
        }

        columns.push({
            name: name,
            type: type,
            primary_key: isPrimary,
            not_null: isNotNull,
            default_value: defaultValue || null
        });
    });

    if (hasError) return;

    try {
        const response = await fetch('/admin/api/create-table', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({
                table_name: tableName,
                columns: columns
            })
        });

        if (!response.ok) {
            const errorText = await response.text();
            throw new Error(errorText);
        }

        showNotification(`Table "${tableName}" created successfully`, 'success');

        // Reset form
        document.getElementById('newTableName').value = '';
        document.getElementById('columnsContainer').innerHTML = `
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
        `;

        // Refresh schema and table list
        loadSchema();
        loadTableList();

    } catch (error) {
        console.error('Error creating table:', error);
        showNotification(`Failed to create table: ${error.message}`, 'error');
    }
}

/**
 * Ajouter un nouvel enregistrement
 */
async function addNewRecord() {
    if (!currentTable) {
        showNotification('No table selected', 'error');
        return;
    }

    try {
        // Get table schema first to create form
        const response = await fetch('/admin/api/schema');
        if (!response.ok) throw new Error(`HTTP error! status: ${response.status}`);

        const schema = await response.json();
        const tableSchema = schema.find(t => t.name === currentTable);

        if (!tableSchema) {
            showNotification('Table schema not found', 'error');
            return;
        }

        openEditModal(null, tableSchema.columns);

    } catch (error) {
        console.error('Error preparing new record:', error);
        showNotification('Failed to prepare new record form', 'error');
    }
}

/**
 * Éditer un enregistrement
 */
async function editRecord(recordId) {
    try {
        // Get table schema
        const schemaResponse = await fetch('/admin/api/schema');
        if (!schemaResponse.ok) throw new Error(`HTTP error! status: ${schemaResponse.status}`);

        const schema = await schemaResponse.json();
        const tableSchema = schema.find(t => t.name === currentTable);

        if (!tableSchema) {
            showNotification('Table schema not found', 'error');
            return;
        }

        // Get record data
        const dataResponse = await fetch(`/admin/api/table/${currentTable}`);
        if (!dataResponse.ok) throw new Error(`HTTP error! status: ${dataResponse.status}`);

        const data = await dataResponse.json();
        const record = data.find(r => r.id === recordId);

        if (!record) {
            showNotification('Record not found', 'error');
            return;
        }

        openEditModal(record, tableSchema.columns);

    } catch (error) {
        console.error('Error loading record for edit:', error);
        showNotification('Failed to load record for editing', 'error');
    }
}

/**
 * Supprimer un enregistrement
 */
async function deleteRecord(recordId) {
    if (!confirm('Are you sure you want to delete this record?')) {
        return;
    }

    try {
        const response = await fetch(`/admin/api/table/${currentTable}/${recordId}`, {
            method: 'DELETE'
        });

        if (!response.ok) {
            const errorText = await response.text();
            throw new Error(errorText);
        }

        showNotification('Record deleted successfully', 'success');
        loadTableData(); // Refresh table data

    } catch (error) {
        console.error('Error deleting record:', error);
        showNotification(`Failed to delete record: ${error.message}`, 'error');
    }
}

/**
 * Ouvrir le modal d'édition
 */
function openEditModal(record, columns) {
    currentEditingRecord = record;

    const modal = document.getElementById('editModal');
    const modalTitle = document.getElementById('modalTitle');
    const editForm = document.getElementById('editForm');

    modalTitle.textContent = record ? 'Edit Record' : 'Add New Record';

    // Create form fields
    editForm.innerHTML = columns.map(col => {
        const value = record ? (record[col.name] || '') : '';
        const disabled = record && col.primary_key ? 'disabled' : '';

        return `
            <div class="form-group">
                <label for="edit_${col.name}">${col.name} (${col.type})${col.not_null ? ' *' : ''}</label>
                <input
                    type="text"
                    id="edit_${col.name}"
                    name="${col.name}"
                    value="${escapeHtml(String(value))}"
                    ${disabled}
                    ${col.not_null && !record ? 'required' : ''}
                >
            </div>
        `;
    }).join('');

    modal.style.display = 'block';
}

/**
 * Fermer le modal
 */
function closeModal() {
    const modal = document.getElementById('editModal');
    modal.style.display = 'none';
    currentEditingRecord = null;
}

/**
 * Sauvegarder un enregistrement
 */
async function saveRecord() {
    const editForm = document.getElementById('editForm');
    const formData = new FormData(editForm);
    const data = {};

    for (let [key, value] of formData.entries()) {
        data[key] = value;
    }

    try {
        let response;

        if (currentEditingRecord) {
            // Update existing record
            response = await fetch(`/admin/api/table/${currentTable}/${currentEditingRecord.id}`, {
                method: 'PUT',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify(data)
            });
        } else {
            // Create new record
            response = await fetch(`/admin/api/table/${currentTable}`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify(data)
            });
        }

        if (!response.ok) {
            const errorText = await response.text();
            throw new Error(errorText);
        }

        const action = currentEditingRecord ? 'updated' : 'created';
        showNotification(`Record ${action} successfully`, 'success');

        closeModal();
        loadTableData(); // Refresh table data

    } catch (error) {
        console.error('Error saving record:', error);
        showNotification(`Failed to save record: ${error.message}`, 'error');
    }
}

/**
 * Échapper le HTML pour éviter les injections XSS
 */
function escapeHtml(text) {
    const map = {
        '&': '&amp;',
        '<': '&lt;',
        '>': '&gt;',
        '"': '&quot;',
        "'": '&#039;'
    };

    return text.replace(/[&<>"']/g, function(m) { return map[m]; });
}

/**
 * Afficher une notification
 */
function showNotification(message, type = 'info') {
    // Remove existing notifications
    const existing = document.querySelector('.notification');
    if (existing) {
        existing.remove();
    }

    // Create new notification
    const notification = document.createElement('div');
    notification.className = `notification notification-${type}`;
    notification.textContent = message;

    // Style the notification
    notification.style.cssText = `
        position: fixed;
        top: 20px;
        right: 20px;
        padding: 15px 20px;
        border-radius: 5px;
        color: white;
        font-weight: bold;
        z-index: 1000;
        max-width: 300px;
        box-shadow: 0 2px 10px rgba(0,0,0,0.2);
        transform: translateX(100%);
        transition: transform 0.3s ease;
    `;

    // Set background color based on type
    switch(type) {
        case 'success':
            notification.style.backgroundColor = '#28a745';
            break;
        case 'error':
            notification.style.backgroundColor = '#dc3545';
            break;
        case 'warning':
            notification.style.backgroundColor = '#ffc107';
            notification.style.color = '#212529';
            break;
        default:
            notification.style.backgroundColor = '#17a2b8';
    }

    document.body.appendChild(notification);

    // Animate in
    setTimeout(() => {
        notification.style.transform = 'translateX(0)';
    }, 100);

    // Auto remove after 5 seconds
    setTimeout(() => {
        notification.style.transform = 'translateX(100%)';
        setTimeout(() => {
            if (notification.parentNode) {
                notification.remove();
            }
        }, 300);
    }, 5000);
}

// Fermer le modal si on clique en dehors
window.onclick = function(event) {
    const modal = document.getElementById('editModal');
    if (event.target === modal) {
        closeModal();
    }
}