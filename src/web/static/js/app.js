document.addEventListener('DOMContentLoaded', () => {
    fetchStats();
    fetchContent();

    setInterval(() => {
        fetchStats();
        fetchContent();
    }, 15000);
});

let charts = {};

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
                label: 'Sources by Criticality',
                data: values,
                backgroundColor: 'rgba(6, 182, 212, 0.5)',
                borderColor: 'rgba(6, 182, 212, 1)',
                borderWidth: 1
            }]
        },
        options: {
            responsive: true,
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
