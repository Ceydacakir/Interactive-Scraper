document.addEventListener('DOMContentLoaded', () => {
    fetchStats();
    fetchContent();

    setInterval(() => {
        fetchStats();
        fetchContent();
    }, 15000);
});

// Close modal
window.onclick = function (event) {
    const modal = document.getElementById('contentModal');
    if (event.target == modal) {
        closeModal();
    }
}

let charts = {};

window.switchPanel = function (event, viewId) {
    console.log("Switching panel to:", viewId);
    if (event) event.preventDefault();

    // Hide all views
    document.querySelectorAll('[id^="view-"]').forEach(el => el.classList.add('hidden'));

    // Show selected view
    document.getElementById(`view-${viewId}`).classList.remove('hidden');

    // Update nav links
    document.querySelectorAll('.nav-link').forEach(el => {
        el.classList.remove('bg-cyan-900/20', 'text-cyan-400', 'border', 'border-cyan-900/50');
        el.classList.add('text-slate-400', 'hover:bg-slate-800');
    });

    const activeLink = document.getElementById(`link-${viewId}`);
    activeLink.classList.remove('text-slate-400', 'hover:bg-slate-800');
    activeLink.classList.add('bg-cyan-900/20', 'text-cyan-400', 'border', 'border-cyan-900/50');

    if (viewId === 'sources') {
        fetchSources();
    }
}

async function fetchSources() {
    try {
        const res = await fetch('/api/sources');
        const data = await res.json();

        const tbody = document.getElementById('sources-table-body');
        tbody.innerHTML = '';

        if (!data || data.length === 0) {
            tbody.innerHTML = '<tr><td colspan="4" class="p-3 text-center text-slate-500">No sources found</td></tr>';
            return;
        }

        data.forEach(source => {
            const tr = document.createElement('tr');
            tr.className = 'hover:bg-slate-800/50 transition';
            tr.innerHTML = `
                <td class="p-3 font-mono text-xs text-slate-500">#${source.id}</td>
                <td class="p-3 text-white font-bold">${source.name}</td>
                <td class="p-3 text-cyan-400 text-xs font-mono">${source.url}</td>
                <td class="p-3">
                    <span class="px-2 py-1 rounded text-xs font-bold ${source.criticality_score > 7 ? 'bg-red-900/50 text-red-400' : 'bg-green-900/50 text-green-400'}">
                        Score: ${source.criticality_score}
                    </span>
                </td>
                <td class="p-3">
                    <button onclick="openEditModal(${source.id}, '${source.name}', ${source.criticality_score})" class="text-xs border border-slate-600 hover:bg-slate-700 text-slate-300 px-2 py-1 rounded transition">Edit</button>
                </td>
            `;
            tbody.appendChild(tr);
        });
    } catch (e) {
        console.error("Failed to fetch sources:", e);
    }
}

async function fetchStats() {
    const res = await fetch('/api/stats');
    const data = await res.json();

    document.getElementById('stat-total').innerText = data.total_content;
    document.getElementById('stat-sources').innerText = data.total_sources;

    const list = document.getElementById('criticality-list');
    list.innerHTML = '';
    if (data.criticality) {
        data.criticality.forEach(item => {
            const div = document.createElement('div');
            div.className = 'flex items-center justify-between p-2 rounded bg-slate-800/50';
            div.innerHTML = `
                <span class="text-xs text-slate-400">Puan ${item.Score}</span>
                <div class="h-2 flex-1 mx-2 bg-slate-700 rounded overflow-hidden">
                    <div class="h-full bg-cyan-500" style="width: ${item.Count * 10}%"></div>
                </div>
                <span class="text-xs font-bold text-white">${item.Count}</span>
            `;
            list.appendChild(div);
        });

        renderChart(data.criticality);
    }

    if (data.categories) {
        renderCategoryChart(data.categories);
    }
}

async function fetchContent() {
    const res = await fetch('/api/content');
    const data = await res.json();

    const tbody = document.getElementById('feed-table-body');
    tbody.innerHTML = '';

    data.forEach(item => {
        const tr = document.createElement('tr');
        tr.className = 'hover:bg-slate-800/50 transition';

        const date = new Date(item.created_at).toLocaleTimeString();
        const critColor = item.source.criticality_score > 7 ? 'text-red-500' : (item.source.criticality_score > 4 ? 'text-yellow-500' : 'text-green-500');

        tr.innerHTML = `
            <td class="p-3 whitespace-nowrap text-slate-500 font-mono text-xs">${date}</td>
            <td class="p-3 text-cyan-400">${item.source.name || 'Bilinmiyor'}</td>
            <td class="p-3">
                <span class="px-2 py-0.5 rounded text-[10px] uppercase font-bold bg-slate-700 text-slate-300 border border-slate-600">
                    ${item.category || 'GENERIC'}
                </span>
            </td>
            <td class="p-3 text-white truncate max-w-xs">${item.title}</td>
            <td class="p-3 font-bold ${critColor}">${item.source.criticality_score || '-'}</td>
            <td class="p-3">
                <button onclick="viewDetails(${item.id})" class="text-xs bg-slate-700 hover:bg-slate-600 text-white px-2 py-1 rounded">Bak</button>
            </td>
        `;
        tr.dataset.content = JSON.stringify(item);
        tbody.appendChild(tr);
    });
}

function renderChart(critData) {
    const ctx = document.getElementById('mainChart').getContext('2d');

    const labels = critData.map(d => `Puan ${d.Score}`);
    const values = critData.map(d => d.Count);

    if (charts.main) {
        charts.main.data.labels = labels;
        charts.main.data.datasets[0].data = values;
        charts.main.update();
        return;
    }

    charts.main = new Chart(ctx, {
        type: 'bar',
        data: {
            labels: labels,
            datasets: [{
                label: 'Sources',
                data: values,
                backgroundColor: 'rgba(6, 182, 212, 0.5)',
                borderColor: 'rgba(6, 182, 212, 1)',
                borderWidth: 1
            }]
        },
        options: {
            responsive: true,
            maintainAspectRatio: false,
            scales: {
                y: { beginAtZero: true, grid: { color: 'rgba(255,255,255,0.1)' } },
                x: { grid: { display: false } }
            },
            plugins: {
                legend: { display: false }
            }
        }
    });
}

function renderCategoryChart(catData) {
    const ctx = document.getElementById('categoryChart').getContext('2d');

    const labels = catData.map(d => d.Category || 'Generic');
    const values = catData.map(d => d.Count);

    if (charts.category) {
        charts.category.data.labels = labels;
        charts.category.data.datasets[0].data = values;
        charts.category.update();
        return;
    }

    charts.category = new Chart(ctx, {
        type: 'doughnut',
        data: {
            labels: labels,
            datasets: [{
                data: values,
                backgroundColor: [
                    'rgba(239, 68, 68, 0.6)',  // Red
                    'rgba(249, 115, 22, 0.6)', // Orange
                    'rgba(234, 179, 8, 0.6)',  // Yellow
                    'rgba(34, 197, 94, 0.6)',  // Green
                    'rgba(6, 182, 212, 0.6)',  // Cyan
                    'rgba(59, 130, 246, 0.6)', // Blue
                    'rgba(168, 85, 247, 0.6)'  // Purple
                ],
                borderColor: '#1e293b',
                borderWidth: 2
            }]
        },
        options: {
            responsive: true,
            maintainAspectRatio: false,
            plugins: {
                legend: { position: 'right', labels: { color: '#94a3b8', boxWidth: 10, font: { size: 10 } } }
            }
        }
    });
}

function viewDetails(id) {
    const buttons = document.querySelectorAll('button[onclick^="viewDetails"]');
    for (let btn of buttons) {
        if (btn.getAttribute('onclick') === `viewDetails(${id})`) {
            const row = btn.closest('tr');
            const data = JSON.parse(row.dataset.content);
            openModal(data);
            break;
        }
    }
}

function openModal(item) {
    document.getElementById('modal-title').innerText = item.title;
    document.getElementById('modal-source').innerText = `Source: ${item.source.name} (${item.source.url})`;
    document.getElementById('modal-date').innerText = `Published: ${new Date(item.publish_date).toLocaleString()}`;
    document.getElementById('modal-content').innerText = item.raw_content;

    const modal = document.getElementById('contentModal');
    modal.classList.remove('hidden');
    modal.classList.add('flex');
}

function closeModal() {
    const modal = document.getElementById('contentModal');
    modal.classList.add('hidden');
    modal.classList.remove('flex');
}

// Edit Source Modal Functions
function openEditModal(id, name, score) {
    document.getElementById('edit-source-id').value = id;
    document.getElementById('edit-source-name').value = name;
    document.getElementById('edit-source-score').value = score;

    const modal = document.getElementById('editSourceModal');
    modal.classList.remove('hidden');
    modal.classList.add('flex');
}

function closeEditModal() {
    const modal = document.getElementById('editSourceModal');
    modal.classList.add('hidden');
    modal.classList.remove('flex');
}

async function saveSourceCriticality() {
    const id = document.getElementById('edit-source-id').value;
    const score = parseInt(document.getElementById('edit-source-score').value);

    if (isNaN(score) || score < 1 || score > 10) {
        alert("Please enter a valid score between 1 and 10.");
        return;
    }

    try {
        const res = await fetch(`/api/sources/${id}/criticality`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ score: score })
        });

        if (res.ok) {
            closeEditModal();
            fetchSources(); // Refresh list to show new score
            fetchStats();   // Refresh stats to show new distribution
        } else {
            alert("Failed to update criticality score.");
        }
    } catch (e) {
        console.error("Error updating source:", e);
        alert("An error occurred.");
    }
}

// Add New Source Modal Functions
function openAddSourceModal() {
    const modal = document.getElementById('addSourceModal');
    modal.classList.remove('hidden');
    modal.classList.add('flex');
}

function closeAddSourceModal() {
    const modal = document.getElementById('addSourceModal');
    modal.classList.add('hidden');
    modal.classList.remove('flex');

    // Clear inputs
    document.getElementById('add-source-name').value = '';
    document.getElementById('add-source-url').value = '';
}

async function createSource() {
    const name = document.getElementById('add-source-name').value;
    const url = document.getElementById('add-source-url').value;

    if (!name || !url) {
        alert("Please provide both name and URL.");
        return;
    }

    try {
        const res = await fetch('/api/sources', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ name: name, url: url })
        });

        if (res.ok) {
            closeAddSourceModal();
            fetchSources(); // Refresh list
            alert("Source added successfully! The scraper will pick it up in the next cycle.");
        } else {
            const data = await res.json();
            alert("Failed to add source: " + (data.error || "Unknown error"));
        }
    } catch (e) {
        console.error("Error creating source:", e);
        alert("An error occurred.");
    }
}
