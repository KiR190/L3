// Global State
let gridApi;
let currentFilters = {
    from: null,
    to: null,
    type: null
};

// AG-Grid Configuration
const columnDefs = [
    {
        headerName: 'Тип',
        field: 'type',
        width: 120,
        cellRenderer: params => {
            const type = params.value;
            const icon = type === 'income' ? '💰' : '💸';
            const label = type === 'income' ? 'Доход' : 'Расход';
            const color = type === 'income' ? '#10b981' : '#ef4444';
            return `<span style="color: ${color}; font-weight: 600;">${icon} ${label}</span>`;
        },
        filter: true
    },
    {
        headerName: 'Сумма',
        field: 'amount',
        width: 150,
        valueFormatter: params => {
            const amount = params.value / 100;
            return new Intl.NumberFormat('ru-RU', {
                minimumFractionDigits: 2,
                maximumFractionDigits: 2
            }).format(amount);
        },
        cellStyle: params => {
            const type = params.data.type;
            return {
                color: type === 'income' ? '#10b981' : '#ef4444',
                fontWeight: '600'
            };
        },
        sort: 'desc'
    },
    {
        headerName: 'Валюта',
        field: 'currency',
        width: 100
    },
    {
        headerName: 'Категория',
        field: 'category_id',
        width: 150,
        valueFormatter: params => params.value || '—',
        filter: true
    },
    {
        headerName: 'Описание',
        field: 'description',
        flex: 1,
        minWidth: 200,
        valueFormatter: params => params.value || '—'
    },
    {
        headerName: 'Дата',
        field: 'occurred_at',
        width: 180,
        valueFormatter: params => {
            const date = new Date(params.value);
            return new Intl.DateTimeFormat('ru-RU', {
                year: 'numeric',
                month: 'long',
                day: 'numeric',
                hour: '2-digit',
                minute: '2-digit'
            }).format(date);
        },
        sort: 'desc',
        filter: 'agDateColumnFilter'
    },
    {
        headerName: 'Действия',
        field: 'id',
        width: 120,
        cellRenderer: params => {
            return `<button class="btn-danger" onclick="deleteItem('${params.value}')">🗑️ Удалить</button>`;
        },
        sortable: false,
        filter: false
    }
];

const gridOptions = {
    columnDefs: columnDefs,
    defaultColDef: {
        sortable: true,
        filter: true,
        resizable: true
    },
    pagination: true,
    paginationPageSize: 20,
    paginationPageSizeSelector: [10, 20, 50, 100],
    animateRows: true,
    rowSelection: 'multiple',
    suppressCellFocus: true,
    domLayout: 'normal'
};

// Initialize App
document.addEventListener('DOMContentLoaded', () => {
    // Initialize AG-Grid
    const gridDiv = document.querySelector('#myGrid');
    gridApi = agGrid.createGrid(gridDiv, gridOptions);
    
    // Set current datetime for form
    const now = new Date();
    now.setMinutes(now.getMinutes() - now.getTimezoneOffset());
    document.getElementById('date').value = now.toISOString().slice(0, 16);
    
    // Load initial data
    loadItems();
    loadAnalytics();
    
    // Setup event listeners
    setupEventListeners();
});

// Event Listeners
function setupEventListeners() {
    // Form submission
    document.getElementById('add-item-form').addEventListener('submit', handleFormSubmit);
    
    // Filter buttons
    document.getElementById('apply-filters').addEventListener('click', applyFilters);
    document.getElementById('reset-filters').addEventListener('click', resetFilters);
    
    // Export CSV
    document.getElementById('export-csv').addEventListener('click', exportCSV);
}

// API Functions
async function loadItems() {
    try {
        const response = await fetch('/items/');
        if (!response.ok) throw new Error('Failed to load items');
        
        const items = await response.json();
        gridApi.setGridOption('rowData', items);
    } catch (error) {
        showToast('Ошибка загрузки данных: ' + error.message, 'error');
    }
}

async function loadAnalytics() {
    try {
        let url = '/analytics';
        const params = new URLSearchParams();
        
        if (currentFilters.from) {
            params.append('from', currentFilters.from);
        }
        if (currentFilters.to) {
            params.append('to', currentFilters.to);
        }
        
        if (params.toString()) {
            url += '?' + params.toString();
        }
        
        const response = await fetch(url);
        if (!response.ok) throw new Error('Failed to load analytics');
        
        const data = await response.json();
        updateAnalytics(data);
    } catch (error) {
        console.error('Error loading analytics:', error);
    }
}

function updateAnalytics(data) {
    document.getElementById('stat-sum').textContent = formatCurrency(data.sum || 0);
    document.getElementById('stat-avg').textContent = formatCurrency(data.avg || 0);
    document.getElementById('stat-count').textContent = data.count || 0;
    document.getElementById('stat-median').textContent = formatCurrency(data.median || 0);
}

async function handleFormSubmit(e) {
    e.preventDefault();
    
    const formData = new FormData(e.target);
    const data = {
        type: formData.get('type'),
        amount: parseFloat(formData.get('amount')),
        currency: formData.get('currency') || 'RUB',
        date: new Date(formData.get('date')).toISOString(),
        category_id: formData.get('category') || null,
        description: formData.get('description') || null
    };
    
    try {
        const response = await fetch('/items/', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(data)
        });
        
        if (!response.ok) {
            const error = await response.json();
            throw new Error(error.error || 'Failed to create item');
        }
        
        showToast('✓ Запись успешно добавлена!', 'success');
        e.target.reset();
        
        // Reset date to current
        const now = new Date();
        now.setMinutes(now.getMinutes() - now.getTimezoneOffset());
        document.getElementById('date').value = now.toISOString().slice(0, 16);
        
        // Reload data
        await loadItems();
        await loadAnalytics();
    } catch (error) {
        showToast('Ошибка: ' + error.message, 'error');
    }
}

async function deleteItem(id) {
    if (!confirm('Вы уверены, что хотите удалить эту запись?')) {
        return;
    }
    
    try {
        const response = await fetch(`/items/${id}`, {
            method: 'DELETE'
        });
        
        if (!response.ok) throw new Error('Failed to delete item');
        
        showToast('✓ Запись удалена', 'success');
        await loadItems();
        await loadAnalytics();
    } catch (error) {
        showToast('Ошибка удаления: ' + error.message, 'error');
    }
}

// Filter Functions
function applyFilters() {
    const from = document.getElementById('filter-from').value;
    const to = document.getElementById('filter-to').value;
    const type = document.getElementById('filter-type').value;
    
    currentFilters = {
        from: from ? new Date(from).toISOString() : null,
        to: to ? new Date(to + 'T23:59:59').toISOString() : null,
        type: type || null
    };
    
    // Apply client-side filtering
    gridApi.setGridOption('quickFilterText', '');
    
    const filterModel = {};
    
    if (currentFilters.type) {
        filterModel.type = {
            type: 'equals',
            filter: currentFilters.type
        };
    }
    
    gridApi.setFilterModel(filterModel);
    
    // Reload analytics with filters
    loadAnalytics();
    
    showToast('✓ Фильтры применены', 'success');
}

function resetFilters() {
    document.getElementById('filter-from').value = '';
    document.getElementById('filter-to').value = '';
    document.getElementById('filter-type').value = '';
    
    currentFilters = {
        from: null,
        to: null,
        type: null
    };
    
    gridApi.setFilterModel(null);
    loadAnalytics();
    
    showToast('✓ Фильтры сброшены', 'success');
}

// Export Functions 
function exportCSV() {
    let url = '/export/csv';
    const params = new URLSearchParams();
    
    if (currentFilters.from) {
        params.append('from', currentFilters.from);
    }
    if (currentFilters.to) {
        params.append('to', currentFilters.to);
    }
    if (currentFilters.type) {
        params.append('type', currentFilters.type);
    }
    
    if (params.toString()) {
        url += '?' + params.toString();
    }
    
    // Create temporary link and trigger download
    const link = document.createElement('a');
    link.href = url;
    link.download = `sales-tracker-${new Date().toISOString().split('T')[0]}.csv`;
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);
    
    showToast('✓ CSV файл загружен', 'success');
}

// Utility Functions
function formatCurrency(value) {
    return new Intl.NumberFormat('ru-RU', {
        style: 'currency',
        currency: 'RUB',
        minimumFractionDigits: 2
    }).format(value);
}

function showToast(message, type = 'success') {
    const toast = document.getElementById('toast');
    toast.textContent = message;
    toast.className = 'toast show';
    
    if (type === 'error') {
        toast.classList.add('error');
    }
    
    setTimeout(() => {
        toast.classList.remove('show');
        setTimeout(() => {
            toast.classList.remove('error');
        }, 300);
    }, 3000);
}

// Make deleteItem available globally
window.deleteItem = deleteItem;
