const API = '/api';
const WS_URL = `${window.location.protocol === 'https:' ? 'wss:' : 'ws:'}//${window.location.host}/api/ws`;

const state = {
    readings: [],
    activeAlerts: [],
    chart: null,
    ws: null,
    reconnectTimer: null,
    status: {},

    get latest() {
        return this.readings[this.readings.length - 1] || null;
    }
};

function $(id) { return document.getElementById(id); }

function fmtTime(iso) {
    if (!iso) return '---';
    const d = new Date(iso);
    return d.toLocaleTimeString('es-ES', { hour: '2-digit', minute: '2-digit', second: '2-digit' });
}

function fmtDateTime(iso) {
    if (!iso) return '---';
    const d = new Date(iso);
    return d.toLocaleString('es-ES', {
        day: '2-digit', month: '2-digit', hour: '2-digit', minute: '2-digit', second: '2-digit'
    });
}

async function apiFetch(path) {
    const res = await fetch(`${API}${path}`);
    if (!res.ok) throw new Error(`API ${res.status}`);
    return res.json();
}

function setConnStatus(state_) {
    const dot = $('connDot');
    const text = $('connText');
    dot.className = 'dot';
    if (state_ === 'connected') {
        dot.classList.add('connected');
        text.textContent = 'Conectado';
    } else if (state_ === 'connecting') {
        dot.classList.add('connecting');
        text.textContent = 'Conectando...';
    } else {
        dot.classList.add('disconnected');
        text.textContent = 'Desconectado';
    }
}

function updateLastUpdate() {
    $('lastUpdate').textContent = `Última actualización: ${fmtTime(new Date().toISOString())}`;
}

async function loadInitialData() {
    try {
        const [statusData, readingsData, alertsData] = await Promise.all([
            apiFetch('/status'),
            apiFetch('/readings?limit=20'),
            apiFetch('/alerts/active')
        ]);
        state.status = statusData;
        state.readings = readingsData.slice().reverse();
        state.activeAlerts = alertsData;
        renderAll();
    } catch (err) {
        console.error('load initial data:', err);
    }
}

function renderAll() {
    renderStatusCards();
    renderSensorCards();
    updateAlertsUI();
    updateTable();
    updateChart();
    updateReadingCount();
    updateLastUpdate();
}

function renderStatusCards() {
    const latest = state.latest;
    const status = state.status;

    const beltVal = $('beltValue');
    const beltInd = $('beltIndicator');
    const on = latest ? latest.belt_running : status.belt_running;
    beltVal.textContent = on ? 'Encendida' : 'Apagada';
    beltVal.style.color = on ? 'var(--success)' : 'var(--text-muted)';
    beltInd.className = 'card-indicator ' + (on ? 'active' : 'inactive');

    const fanVal = $('fanValue');
    const fanInd = $('fanIndicator');
    const fanOn = latest ? latest.fan_on : status.fan_on;
    fanVal.textContent = fanOn ? 'Encendido' : 'Apagado';
    fanVal.style.color = fanOn ? 'var(--success)' : 'var(--text-muted)';
    fanInd.className = 'card-indicator ' + (fanOn ? 'active' : 'inactive');

    const buzzerVal = $('buzzerValue');
    const buzzerInd = $('buzzerIndicator');
    const buzzerOn = latest ? latest.buzzer_on : status.buzzer_on;
    buzzerVal.textContent = buzzerOn ? 'Encendido' : 'Apagado';
    buzzerVal.style.color = buzzerOn ? 'var(--success)' : 'var(--text-muted)';
    buzzerInd.className = 'card-indicator ' + (buzzerOn ? 'active' : 'inactive');

    const doorVal = $('doorValue');
    const angle = latest ? latest.door_angle : status.door_angle;
    doorVal.textContent = angle != null ? `${angle}°` : '---';
    doorVal.style.color = 'var(--accent)';

    const arc = $('doorArc');
    const needle = $('doorNeedle');
    if (angle != null) {
        const offset = 283 - (angle / 180) * 283;
        arc.setAttribute('stroke-dashoffset', offset);
        needle.setAttribute('transform', `rotate(${angle - 90} 50 50)`);
    }
}

function renderSensorCards() {
    const latest = state.latest;
    if (!latest) return;

    const gasVal = $('gasValue');
    gasVal.textContent = latest.gas_value;
    gasVal.style.color = latest.gas_value > 300 ? 'var(--danger)' : 'var(--text-primary)';
    const gasCard = gasVal.closest('.sensor-card');
    gasCard.dataset.alert = latest.gas_value > 300;

    const humVal = $('humidityValue');
    humVal.textContent = latest.humidity_value;
    humVal.style.color = latest.humidity_value > 2500 ? 'var(--danger)' : 'var(--text-primary)';
    const humCard = humVal.closest('.sensor-card');
    humCard.dataset.alert = latest.humidity_value > 2500;

    $('distanceValue').textContent = latest.distance_cm.toFixed(1);
    $('objectCount').textContent = latest.object_count;
}

function updateAlertsUI() {
    const container = $('alertsList');
    const badge = $('alertsBadge');

    badge.textContent = state.activeAlerts.length;

    if (state.activeAlerts.length === 0) {
        container.innerHTML = '<div class="alert-empty">Sin alertas activas</div>';
        return;
    }

    container.innerHTML = state.activeAlerts.map(a => {
        const isGas = a.type === 'GAS';
        const cls = isGas ? 'alert-gas' : 'alert-humidity';
        const icon = isGas ? '💨' : '💧';
        const label = isGas ? 'Gas elevado' : 'Humedad elevada';
        return `
            <div class="alert-item ${cls}" data-id="${a.id}">
                <div class="alert-icon">${icon}</div>
                <div class="alert-info">
                    <div class="alert-type">${label}</div>
                    <div class="alert-detail">${a.trigger_value} &gt; ${a.threshold}</div>
                    <div class="alert-time">${fmtDateTime(a.timestamp)}</div>
                </div>
                <button class="alert-dismiss" data-id="${a.id}">Resolver</button>
            </div>
        `;
    }).join('');

    container.querySelectorAll('.alert-dismiss').forEach(btn => {
        btn.addEventListener('click', async () => {
            const id = btn.dataset.id;
            try {
                await fetch(`${API}/alerts/${id}/resolve`, { method: 'POST' });
                state.activeAlerts = state.activeAlerts.filter(a => a.id !== id);
                updateAlertsUI();
            } catch (err) {
                console.error('resolve alert:', err);
            }
        });
    });
}

function updateTable() {
    const tbody = $('readingsBody');
    if (state.readings.length === 0) {
        tbody.innerHTML = '<tr><td colspan="9" style="text-align:center;color:var(--text-muted);padding:32px;">Sin lecturas</td></tr>';
        return;
    }

    tbody.innerHTML = state.readings.slice().reverse().slice(0, 50).map(r => `
        <tr>
            <td>${fmtDateTime(r.timestamp)}</td>
            <td style="color:${r.gas_value > 300 ? 'var(--danger)' : 'inherit'}">${r.gas_value}</td>
            <td style="color:${r.humidity_value > 2500 ? 'var(--danger)' : 'inherit'}">${r.humidity_value}</td>
            <td>${r.distance_cm.toFixed(1)}</td>
            <td>${r.object_count}</td>
            <td><span class="tag ${r.belt_running ? 'tag-on' : 'tag-off'}">${r.belt_running ? 'Sí' : 'No'}</span></td>
            <td><span class="tag ${r.fan_on ? 'tag-on' : 'tag-off'}">${r.fan_on ? 'Sí' : 'No'}</span></td>
            <td><span class="tag ${r.buzzer_on ? 'tag-on' : 'tag-off'}">${r.buzzer_on ? 'Sí' : 'No'}</span></td>
            <td>${r.door_angle}°</td>
        </tr>
    `).join('');
}

function updateReadingCount() {
    $('readingCount').textContent = `${state.readings.length} lecturas`;
}

function fmtChartLabel(iso) {
    const d = new Date(iso);
    return d.toLocaleTimeString('es-ES', { hour: '2-digit', minute: '2-digit', second: '2-digit' });
}

function buildChart() {
    const ctx = document.getElementById('readingsChart').getContext('2d');
    const labels = state.readings.map(r => fmtChartLabel(r.timestamp));

    state.chart = new Chart(ctx, {
        type: 'line',
        data: {
            labels,
            datasets: [
                {
                    label: 'Gas (ppm)',
                    data: state.readings.map(r => r.gas_value),
                    borderColor: '#ff5252',
                    backgroundColor: 'rgba(255, 82, 82, 0.1)',
                    fill: true,
                    tension: 0.3,
                    pointRadius: 2,
                    pointHoverRadius: 6,
                    yAxisID: 'y',
                },
                {
                    label: 'Humedad (RH)',
                    data: state.readings.map(r => r.humidity_value),
                    borderColor: '#40c4ff',
                    backgroundColor: 'rgba(64, 196, 255, 0.1)',
                    fill: true,
                    tension: 0.3,
                    pointRadius: 2,
                    pointHoverRadius: 6,
                    yAxisID: 'y',
                },
                {
                    label: 'Distancia (cm)',
                    data: state.readings.map(r => r.distance_cm),
                    borderColor: '#ffd740',
                    backgroundColor: 'rgba(255, 215, 64, 0.1)',
                    fill: true,
                    tension: 0.3,
                    pointRadius: 2,
                    pointHoverRadius: 6,
                    yAxisID: 'y1',
                    hidden: true,
                }
            ]
        },
        options: {
            responsive: true,
            maintainAspectRatio: false,
            animation: { duration: 200 },
            interaction: {
                intersect: false,
                mode: 'index',
            },
            scales: {
                x: {
                    grid: { color: 'rgba(255,255,255,0.05)' },
                    ticks: { color: '#8888aa', maxRotation: 0, maxTicksLimit: 10, font: { size: 10 } }
                },
                y: {
                    beginAtZero: true,
                    grid: { color: 'rgba(255,255,255,0.05)' },
                    ticks: { color: '#ff5252' },
                    title: { display: true, text: 'Gas / Humedad', color: '#8888aa' }
                },
                y1: {
                    beginAtZero: true,
                    position: 'right',
                    grid: { display: false },
                    ticks: { color: '#ffd740' },
                    title: { display: true, text: 'Distancia (cm)', color: '#8888aa' }
                }
            },
            plugins: {
                legend: {
                    labels: { color: '#e0e0f0', usePointStyle: true, padding: 16 }
                }
            }
        }
    });
}

function updateChart() {
    if (!state.chart) {
        buildChart();
        return;
    }

    const labels = state.readings.map(r => fmtChartLabel(r.timestamp));
    state.chart.data.labels = labels;
    state.chart.data.datasets[0].data = state.readings.map(r => r.gas_value);
    state.chart.data.datasets[1].data = state.readings.map(r => r.humidity_value);
    state.chart.data.datasets[2].data = state.readings.map(r => r.distance_cm);
    state.chart.update('none');
}

document.querySelectorAll('.chart-controls input').forEach(cb => {
    cb.addEventListener('change', () => {
        if (!state.chart) return;
        const idx = { gas: 0, humidity: 1, distance: 2 }[cb.dataset.series];
        const meta = state.chart.getDatasetMeta(idx);
        meta.hidden = !cb.checked;
        state.chart.update();
    });
});

function handleNewReading(reading) {
    state.readings.push(reading);
    if (state.readings.length > 200) {
        state.readings = state.readings.slice(-150);
    }
    renderStatusCards();
    renderSensorCards();
    updateAlertsUI();
    updateTable();
    updateChart();
    updateReadingCount();
    updateLastUpdate();
}

function handleNewAlert(alertId) {
    apiFetch('/alerts/active').then(alerts => {
        state.activeAlerts = alerts;
        updateAlertsUI();
    }).catch(() => {});
}

function handleAlertResolved(alertData) {
    state.activeAlerts = state.activeAlerts.filter(a => a.id !== alertData.id);
    updateAlertsUI();
}

function connectWS() {
    setConnStatus('connecting');
    state.ws = new WebSocket(WS_URL);

    state.ws.onopen = () => setConnStatus('connected');

    state.ws.onmessage = (ev) => {
        try {
            const msg = JSON.parse(ev.data);
            switch (msg.type) {
                case 'reading':
                    handleNewReading(msg.data);
                    break;
                case 'alert':
                    handleNewAlert(msg.data);
                    break;
                case 'alert_resolved':
                    handleAlertResolved(msg.data);
                    break;
            }
        } catch (err) {
            console.error('ws message:', err);
        }
    };

    state.ws.onclose = () => {
        setConnStatus('disconnected');
        state.reconnectTimer = setTimeout(connectWS, 2000);
    };

    state.ws.onerror = () => {
        state.ws.close();
    };
}

async function init() {
    setConnStatus('connecting');
    await loadInitialData();
    connectWS();
}

document.addEventListener('DOMContentLoaded', init);
