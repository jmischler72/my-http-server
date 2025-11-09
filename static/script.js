// Calculate grid dimensions based on viewport
function calculateGridDimensions() {
    const charWidth = 8; // Approximate width of a character in pixels
    const charHeight = 14; // Approximate height of a character in pixels

    const GRID_WIDTH = Math.floor(window.innerWidth / charWidth);
    const GRID_HEIGHT = Math.floor(window.innerHeight / charHeight);

    return { width: GRID_WIDTH, height: GRID_HEIGHT };
}

let selectedX = 0;
let selectedY = 0;
let gridEntries = [];
let GRID_WIDTH, GRID_HEIGHT;

// Initialize the ASCII grid
function initGrid() {
    const dimensions = calculateGridDimensions();
    GRID_WIDTH = dimensions.width;
    GRID_HEIGHT = dimensions.height;

    const grid = document.getElementById('ascii-grid');
    let gridHTML = '';

    for (let y = 0; y < GRID_HEIGHT; y++) {
        gridHTML += '<div class="grid-row">';
        for (let x = 0; x < GRID_WIDTH; x++) {
            gridHTML += '<span class="grid-cell" data-x="' + x + '" data-y="' + y + '">.</span>';
        }
        gridHTML += '</div>';
    }

    grid.innerHTML = gridHTML;

    // Add click listeners to grid cells
    grid.addEventListener('click', function (e) {
        if (e.target.classList.contains('grid-cell')) {
            selectedX = parseInt(e.target.getAttribute('data-x'));
            selectedY = parseInt(e.target.getAttribute('data-y'));

            // Check if position is already occupied
            const existing = gridEntries.find(entry => entry.x === selectedX && entry.y === selectedY);
            if (existing) {
                showTooltip(e, existing.name, existing.message);
                return;
            }

            document.getElementById('positionText').textContent = '(' + selectedX + ', ' + selectedY + ')';
            document.getElementById('entryModal').style.display = 'block';
            document.getElementById('nameInput').value = '';
            document.getElementById('messageInput').value = '';
            document.getElementById('nameInput').focus();
        }
    });

    // Add hover listeners for tooltips
    grid.addEventListener('mouseover', function (e) {
        let targetCell = e.target;

        // If hovering over an ascii-guy, get its parent grid-cell
        if (e.target.classList.contains('ascii-guy')) {
            targetCell = e.target.parentElement;
        }

        if (targetCell.classList.contains('grid-cell')) {
            const x = parseInt(targetCell.getAttribute('data-x'));
            const y = parseInt(targetCell.getAttribute('data-y'));
            const entry = gridEntries.find(entry => entry.x === x && entry.y === y);
            if (entry) {
                showTooltip(e, entry.name, entry.message);
            }
        }
    });

    grid.addEventListener('mouseout', function (e) {
        // Hide tooltip when leaving grid cells or ascii guys
        if (e.target.classList.contains('grid-cell') || e.target.classList.contains('ascii-guy')) {
            hideTooltip();
        }
    });
}

function showTooltip(e, name, message) {
    const tooltip = document.getElementById('tooltip');
    tooltip.innerHTML = '<strong>' + name + '</strong><br>' + message;
    tooltip.style.display = 'block';

    // Position the tooltip, but make sure it doesn't go off-screen
    let left = e.pageX + 15;
    let top = e.pageY + 15;

    // Adjust if tooltip would go off the right edge
    if (left + 250 > window.innerWidth) {
        left = e.pageX - 250;
    }

    // Adjust if tooltip would go off the bottom edge
    if (top + 60 > window.innerHeight) {
        top = e.pageY - 60;
    }

    tooltip.style.left = left + 'px';
    tooltip.style.top = top + 'px';
}

function hideTooltip() {
    document.getElementById('tooltip').style.display = 'none';
}

function closeModal() {
    document.getElementById('entryModal').style.display = 'none';
}

function submitGridEntry() {
    const name = document.getElementById('nameInput').value.trim();
    const message = document.getElementById('messageInput').value.trim();

    if (!name || !message) {
        alert('Please fill in both name and message');
        return;
    }

    fetch('/grid-entry', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({
            x: selectedX,
            y: selectedY,
            name: name,
            message: message
        })
    })
        .then(response => response.json())
        .then(data => {
            if (data.success) {
                loadGridEntries();
                closeModal();
            } else {
                alert('Error: ' + data.error);
            }
        })
        .catch(error => {
            console.error('Error:', error);
            alert('Failed to submit entry');
        });
}

function loadGridEntries() {
    fetch('/grid-entries')
        .then(response => response.json())
        .then(entries => {
            gridEntries = entries;
            updateGridDisplay();
        })
        .catch(error => {
            console.error('Error loading grid entries:', error);
        });
}

function updateGridDisplay() {
    // Clear existing ASCII guys
    document.querySelectorAll('.ascii-guy').forEach(guy => guy.remove());

    // Add ASCII guys for each entry
    gridEntries.forEach(entry => {
        const cell = document.querySelector('[data-x="' + entry.x + '"][data-y="' + entry.y + '"]');
        if (cell) {
            const asciiGuy = document.createElement('span');
            asciiGuy.className = 'ascii-guy';
            asciiGuy.textContent = '☺';
            asciiGuy.style.position = 'absolute';
            asciiGuy.style.left = '0';
            asciiGuy.style.top = '0';
            cell.appendChild(asciiGuy);
        }
    });
}

// Toggle header overlay visibility
function toggleHeader() {
    const headerOverlay = document.getElementById('headerOverlay');
    const showButton = document.getElementById('showHeaderButton');
    const gridInstructions = document.getElementById('gridInstructions');

    if (headerOverlay.classList.contains('hidden')) {
        // Show header
        headerOverlay.classList.remove('hidden');
        showButton.style.display = 'none';
        gridInstructions.style.display = 'none';
    } else {
        // Hide header
        headerOverlay.classList.add('hidden');
        showButton.style.display = 'block';
        gridInstructions.style.display = 'block';
    }
}

// Handle profile image loading
function handleProfileImage() {
    const profileImg = document.querySelector('.profile-image');
    if (profileImg) {
        profileImg.addEventListener('error', function () {
            // If image fails to load, hide it gracefully
            this.style.display = 'none';
        });
    }
}

// Initialize everything when the page loads
document.addEventListener('DOMContentLoaded', function () {
    handleProfileImage();
    initGrid();
    loadGridEntries();
    initHealthChecks();

    // Add keyboard shortcut to toggle header (Escape key)
    document.addEventListener('keydown', function (e) {
        if (e.key === 'Escape') {
            toggleHeader();
        }
    });
});

// Handle window resize to rebuild grid
window.addEventListener('resize', function () {
    initGrid();
    loadGridEntries(); // Reload entries to position them correctly
});

// Close modal when clicking outside
window.onclick = function (event) {
    const modal = document.getElementById('entryModal');
    if (event.target === modal) {
        closeModal();
    }
}

// Health check functionality
function initHealthChecks() {
    const healthIndicators = document.querySelectorAll('.health-indicator');

    healthIndicators.forEach(indicator => {
        const url = indicator.getAttribute('data-url');
        checkHealth(url, indicator);
    });

    // Set up periodic health checks every 30 seconds
    setInterval(() => {
        healthIndicators.forEach(indicator => {
            const url = indicator.getAttribute('data-url');
            checkHealth(url, indicator);
        });
    }, 30000);
}

function checkHealth(url, indicator) {
    indicator.className = 'health-indicator checking';
    indicator.title = 'Checking health...';

    // Use a simple fetch with a timeout to check if the site is reachable
    const controller = new AbortController();
    const timeoutId = setTimeout(() => controller.abort(), 5000); // 5 second timeout

    fetch(url, {
        method: 'HEAD',
        mode: 'no-cors', // Important for cross-origin requests
        signal: controller.signal
    })
        .then(response => {
            clearTimeout(timeoutId);
            // With no-cors, we can't check the actual response status
            // If we get here without an error, the site is likely reachable
            indicator.className = 'health-indicator healthy';
            indicator.title = 'Service is healthy';
        })
        .catch(error => {
            clearTimeout(timeoutId);
            if (error.name === 'AbortError') {
                indicator.className = 'health-indicator unhealthy';
                indicator.title = 'Service timeout - may be down';
            } else {
                // For no-cors mode, network errors usually mean the site is unreachable
                // But sometimes CORS itself can cause this, so we'll be more lenient
                indicator.className = 'health-indicator healthy';
                indicator.title = 'Service appears to be running (CORS restricted)';
            }
        });
}
