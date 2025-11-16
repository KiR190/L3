const apiBase = "";

async function createShort() {
    const url = document.getElementById("origUrl").value.trim();
    const custom = document
        .getElementById("customCode")
        .value.trim();
    if (!url) return alert("Введите URL");
    const btn = document.getElementById("createBtn");
    btn.disabled = true;
    btn.textContent = "...";
    try {
        const res = await fetch(apiBase + "/shorten", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ url, custom }),
        });
        if (!res.ok) {
            const txt = await res.text();
            throw new Error(txt || res.statusText);
        }
        const data = await res.json();
        showCreated(data);
        await fetchLatest();
    } catch (e) {
        console.error(e);
        alert("Ошибка создания: " + e.message);
    } finally {
        btn.disabled = false;
        btn.textContent = "Сократить";
    }
}

function showCreated(data) {
    const out = document.getElementById("createResult");
    const short =
        data.ShortCode ||
        data.short_code ||
        data.ShortCode ||
        data.short ||
        data.ShortURL ||
        data.Short ||
        data.shortCode ||
        data.short;
    // best guess for returned object
    const href = short
        ? window.location.origin + "/s/" + short
        : data.Original ||
          data.original ||
          data.original_url ||
          data.OriginalURL ||
          "";
    out.innerHTML = `<div style="display:flex;gap:10px;align-items:center"><div class="short-badge">${
        short || "—"
    }</div><div class="muted small">${href}</div><div style="margin-left:auto"><button class="btn-ghost" onclick="copyText('${href.replace(
        /'/g,
        "\\'"
    )}')">Копировать</button></div></div>`;
}

function copyText(s) {
    navigator.clipboard
        .writeText(s)
        .then(() => alert("Скопировано"));
}

async function fetchLatest() {
    const list = document.getElementById("latestList");
    list.innerHTML = '<div class="muted small">Загрузка...</div>';
    try {
        const res = await fetch(apiBase + "/analytics/latest");
        if (!res.ok) throw new Error(await res.text());
        const items = await res.json();
        renderLatest(items);
    } catch (e) {
        console.error(e);
        list.innerHTML =
            '<div class="muted small">Не удалось загрузить</div>';
    }
}

function renderLatest(items) {
    const list = document.getElementById("latestList");
    if (!Array.isArray(items)) {
        // maybe returned as {items: [...]}
        items = items.items || items.data || [];
    }
    if (items.length === 0) {
        list.innerHTML =
            '<div class="muted small">Нет ссылок</div>';
        return;
    }
    list.innerHTML = "";
    items.forEach((it) => {
        const code =
            it.ShortCode ||
            it.short_code ||
            it.short ||
            it.Short ||
            it.shortCode ||
            "";
        const original =
            it.Original ||
            it.original ||
            it.OriginalURL ||
            it.original_url ||
            "";
        const created =
            it.CreatedAt ||
            it.createdAt ||
            it.created_at ||
            it.created ||
            "";
        const node = document.createElement("div");
        node.className = "link-item";
        node.innerHTML = `<div class="link-left"><div class="short-badge">${code}</div><div><div style="font-weight:600">${original}</div><div class="muted small">${
            created ? dayjs(created).format("YYYY-MM-DD HH:mm") : ""
        }</div></div></div>
      <div class="actions">
        <button class="btn-ghost" onclick="openInNew('${
            window.location.origin + "/s/" + code
        }')">Открыть</button>
        <button class="btn-ghost" onclick="loadAnalytics('${code}')">Аналитика</button>
      </div>`;
        list.appendChild(node);
    });
}

function openInNew(u) {
    window.open(u, "_blank");
}

let timeseriesChart = null,
    uaChart = null;

async function loadAnalytics(code) {
    if (!code) return;
    document.getElementById(
        "analyticsTitle"
    ).textContent = `Аналитика — ${code}`;
    try {
        const res = await fetch(
            apiBase + "/analytics/" + encodeURIComponent(code)
        );
        if (!res.ok) throw new Error(await res.text());
        const payload = await res.json();
        // payload may be {analytics: stats} or [] of events
        let events = [];
        if (Array.isArray(payload)) events = payload;
        else if (Array.isArray(payload.analytics))
            events = payload.analytics;
        else if (Array.isArray(payload.events))
            events = payload.events;
        else if (Array.isArray(payload.data)) events = payload.data;
        else {
            // try to interpret as aggregated stats -> convert to events if possible
            console.warn("Unknown analytics format, showing raw");
            document.getElementById("eventsTable").innerHTML =
                '<tr><td colspan="2">Неожиданный формат аналитики</td></tr>';
            return;
        }

        // normalize events: look for fields: Timestamp, timestamp, time, UserAgent, user_agent, userAgent
        events = events
            .map((e) => ({
                timestamp:
                    e.Timestamp ||
                    e.timestamp ||
                    e.time ||
                    e.created_at ||
                    e.TimestampUTC ||
                    e.ts,
                ua:
                    e.UserAgent ||
                    e.user_agent ||
                    e.userAgent ||
                    e.ua ||
                    "",
            }))
            .sort(
                (a, b) =>
                    new Date(a.timestamp) - new Date(b.timestamp)
            );

        renderEventsTable(events);
        renderCharts(events);
    } catch (e) {
        console.error(e);
        alert("Не удалось загрузить аналитику: " + e.message);
    }
}

function renderEventsTable(events) {
    const tbody = document.getElementById("eventsTable");
    tbody.innerHTML = "";
    const last = events.slice(-200).reverse();
    if (last.length === 0) {
        tbody.innerHTML =
            '<tr><td colspan="2" class="muted small">Нет событий</td></tr>';
        return;
    }
    last.forEach((ev) => {
        const tr = document.createElement("tr");
        const t = ev.timestamp
            ? dayjs(ev.timestamp).format("YYYY-MM-DD HH:mm:ss")
            : "-";
        tr.innerHTML = `<td>${t}</td><td style="max-width:340px;overflow:hidden;text-overflow:ellipsis;white-space:nowrap">${escapeHtml(
            ev.ua
        )}</td>`;
        tbody.appendChild(tr);
    });
}

function escapeHtml(s) {
    if (!s) return "";
    return s
        .replace(/&/g, "&amp;")
        .replace(/</g, "&lt;")
        .replace(/>/g, "&gt;");
}

function aggregateBy(events, granularity) {
    const map = new Map();
    events.forEach((ev) => {
        const d = dayjs(ev.timestamp);
        const key =
            granularity === "month"
                ? d.format("YYYY-MM")
                : d.format("YYYY-MM-DD");
        map.set(key, (map.get(key) || 0) + 1);
    });
    // sort keys
    const keys = Array.from(map.keys()).sort();
    return { labels: keys, values: keys.map((k) => map.get(k)) };
}

function topUserAgents(events, topN) {
    const map = new Map();
    events.forEach((ev) => {
        const k = (ev.ua || "unknown").split(")")[0];
        map.set(k, (map.get(k) || 0) + 1);
    });
    const arr = Array.from(map.entries())
        .sort((a, b) => b[1] - a[1])
        .slice(0, topN);
    return {
        labels: arr.map((a) => a[0]),
        values: arr.map((a) => a[1]),
    };
}

function renderCharts(events) {
    const gran = document.getElementById("granularity").value;
    const topN = parseInt(
        document.getElementById("topN").value || 5,
        10
    );
    const ts = aggregateBy(events, gran);
    const ua = topUserAgents(events, topN);

    if (timeseriesChart) timeseriesChart.destroy();
    const ctx = document
        .getElementById("timeseriesChart")
        .getContext("2d");
    timeseriesChart = new Chart(ctx, {
        type: "line",
        data: {
            labels: ts.labels,
            datasets: [
                {
                    label: "Клики",
                    data: ts.values,
                    tension: 0.3,
                    fill: true,
                },
            ],
        },
        options: {
            plugins: { legend: { display: false } },
            scales: {
                x: { grid: { display: false } },
                y: { beginAtZero: true },
            },
        },
    });

    if (uaChart) uaChart.destroy();
    const ctx2 = document
        .getElementById("uaChart")
        .getContext("2d");
    uaChart = new Chart(ctx2, {
        type: "bar",
        data: {
            labels: ua.labels,
            datasets: [{ label: "Клики", data: ua.values }],
        },
        options: {
            plugins: { legend: { display: false } },
            scales: {
                x: {
                    ticks: {
                        autoSkip: true,
                        maxRotation: 45,
                        minRotation: 0,
                    },
                },
            },
        },
    });
}

// UI wiring
document
    .getElementById("createBtn")
    .addEventListener("click", createShort);
document
    .getElementById("granularity")
    .addEventListener("change", () => {
        /* re-render if needed */ const t = document
            .getElementById("analyticsTitle")
            .textContent.split("—")[1];
        if (t) loadAnalytics(t.trim());
    });
document.getElementById("topN").addEventListener("change", () => {
    const t = document
        .getElementById("analyticsTitle")
        .textContent.split("—")[1];
    if (t) loadAnalytics(t.trim());
});

// helper to trigger loadAnalytics from latest items
window.loadAnalytics = loadAnalytics;

// initial fetch
fetchLatest();