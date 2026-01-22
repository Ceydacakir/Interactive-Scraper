document.addEventListener('DOMContentLoaded', () => {
    update();
    setInterval(update, 15000);
});

function update() {
    getStats();
    getContent();
}

window.onclick = function (event) {
    if (event.target == document.getElementById('contentModal')) closeModal();
}

let charts = {};

function switchPanel(e, id) {
    if (e) e.preventDefault();
    document.querySelectorAll('[id^="view-"]').forEach(el => el.classList.add('hidden'));
    document.getElementById(`view-${id}`).classList.remove('hidden');

    if (id === 'sources') getSources();
}

async function getSources() {
    const res = await fetch('/api/sources');
    const data = await res.json();
    const tbody = document.getElementById('sources-table-body');
    tbody.innerHTML = '';

    if (!data || data.length === 0) {
        tbody.innerHTML = '<tr><td colspan="5" class="p-3 text-center">No sources</td></tr>';
        return;
    }

    data.forEach(s => {
        let color = s.criticality_score > 7 ? 'text-red-400' : 'text-green-400';
        tbody.innerHTML += `
            <tr class="hover:bg-slate-800">
                <td class="p-3 text-slate-500">#${s.id}</td>
                <td class="p-3 text-white font-bold">${s.name}</td>
                <td class="p-3 text-cyan-400 text-xs">${s.url}</td>
                <td class="p-3 ${color}">Score: ${s.criticality_score}</td>
                <td class="p-3 flex gap-2">
                    <button onclick="editModal(${s.id}, '${s.name}', ${s.criticality_score})" class="border border-slate-600 text-slate-300 px-2 py-1 rounded text-xs">Edit</button>
                    <button onclick="delSource(${s.id})" class="border border-red-900 text-red-400 px-2 py-1 rounded text-xs">Delete</button>
                </td>
            </tr>`;
    });
}

async function getStats() {
    const res = await fetch('/api/stats');
    const data = await res.json();

    document.getElementById('stat-total').innerText = data.total_content;
    document.getElementById('stat-sources').innerText = data.total_sources;

    const list = document.getElementById('criticality-list');
    list.innerHTML = '';

    if (data.criticality) {
        data.criticality.forEach(i => {
            list.innerHTML += `
                <div class="flex items-center justify-between p-2">
                    <span class="text-xs text-slate-400">Score ${i.Score}</span>
                    <div class="h-2 flex-1 mx-2 bg-slate-700 rounded"><div class="h-full bg-cyan-500" style="width: ${i.Count * 10}%"></div></div>
                    <span class="text-xs font-bold text-white">${i.Count}</span>
                </div>`;
        });
        drawChart(data.criticality);
    }
    if (data.categories) drawPie(data.categories);
}

async function getContent() {
    const res = await fetch('/api/content');
    const data = await res.json();
    const tbody = document.getElementById('feed-table-body');
    tbody.innerHTML = '';

    data.forEach(i => {
        let date = new Date(i.created_at).toLocaleTimeString();
        tbody.innerHTML += `
            <tr class="hover:bg-slate-800">
                <td class="p-3 text-slate-500 text-xs">${date}</td>
                <td class="p-3 text-cyan-400">${i.source.name}</td>
                <td class="p-3"><span class="bg-slate-700 text-slate-300 px-2 rounded text-[10px]">${i.category}</span></td>
                <td class="p-3 text-white truncate max-w-xs">${i.title}</td>
                <td class="p-3 font-bold">${i.source.criticality_score}</td>
                <td class="p-3"><button onclick='show(${JSON.stringify(i)})' class="bg-slate-700 text-white px-2 py-1 rounded text-xs">View</button></td>
            </tr>`;
    });
}

function drawChart(data) {
    const ctx = document.getElementById('mainChart').getContext('2d');
    if (charts.main) charts.main.destroy();

    charts.main = new Chart(ctx, {
        type: 'bar',
        data: {
            labels: data.map(d => `Score ${d.Score}`),
            datasets: [{ label: 'Sources', data: data.map(d => d.Count), backgroundColor: '#06b6d4' }]
        },
        options: { responsive: true, confirmProperties: false, scales: { x: { display: false } }, plugins: { legend: { display: false } } }
    });
}

function drawPie(data) {
    const ctx = document.getElementById('categoryChart').getContext('2d');
    if (charts.pie) charts.pie.destroy();

    charts.pie = new Chart(ctx, {
        type: 'doughnut',
        data: {
            labels: data.map(d => d.Category),
            datasets: [{ data: data.map(d => d.Count), backgroundColor: ['#ef4444', '#f97316', '#eab308', '#22c55e', '#06b6d4', '#3b82f6', '#a855f7'], borderWidth: 0 }]
        },
        options: { responsive: true, plugins: { legend: { position: 'right', labels: { color: '#94a3b8', boxWidth: 10 } } } }
    });
}

function show(item) {
    document.getElementById('modal-title').innerText = item.title;
    document.getElementById('modal-source').innerText = item.source.name + " (" + item.source.url + ")";
    document.getElementById('modal-content').innerText = item.raw_content;
    document.getElementById('contentModal').classList.remove('hidden');
    document.getElementById('contentModal').classList.add('flex');
}

function closeModal() {
    document.getElementById('contentModal').classList.add('hidden');
    document.getElementById('contentModal').classList.remove('flex');
}

function editModal(id, name, score) {
    document.getElementById('edit-source-id').value = id;
    document.getElementById('edit-source-name').value = name;
    document.getElementById('edit-source-score').value = score;
    document.getElementById('editSourceModal').classList.remove('hidden');
    document.getElementById('editSourceModal').classList.add('flex');
}

function closeEditModal() {
    document.getElementById('editSourceModal').classList.add('hidden');
    document.getElementById('editSourceModal').classList.remove('flex');
}

async function saveSourceCriticality() {
    const id = document.getElementById('edit-source-id').value;
    const score = parseInt(document.getElementById('edit-source-score').value);
    await fetch(`/api/sources/${id}/criticality`, { method: 'POST', headers: { 'Content-Type': 'application/json' }, body: JSON.stringify({ score }) });
    closeEditModal();
    getSources();
    getStats();
}

function openAddSourceModal() {
    document.getElementById('addSourceModal').classList.remove('hidden');
    document.getElementById('addSourceModal').classList.add('flex');
}

function closeAddSourceModal() {
    document.getElementById('addSourceModal').classList.add('hidden');
    document.getElementById('addSourceModal').classList.remove('flex');
}

async function createSource() {
    const name = document.getElementById('add-source-name').value;
    const url = document.getElementById('add-source-url').value;
    await fetch('/api/sources', { method: 'POST', headers: { 'Content-Type': 'application/json' }, body: JSON.stringify({ name, url }) });
    closeAddSourceModal();
    getSources();
    alert("Added!");
}

async function delSource(id) {
    if (!confirm("Delete?")) return;
    await fetch(`/api/sources/${id}`, { method: 'DELETE' });
    getSources();
    getStats();
}
