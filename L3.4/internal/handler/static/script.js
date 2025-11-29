const list = document.getElementById("list");
const typeSelect = document.getElementById("taskType");
const placeholderSVG =
    "data:image/svg+xml;base64," +
    btoa(`
    <svg width="90" height="90" viewBox="0 0 100 100" xmlns="http://www.w3.org/2000/svg">
    <rect width="100" height="100" rx="8" ry="8" fill="#e5e5e5"/>
    <circle cx="50" cy="50" r="22" stroke="#bbb" stroke-width="6" fill="none" 
        stroke-dasharray="34 20" stroke-linecap="round">
    <animateTransform attributeName="transform" type="rotate" values="0 50 50;360 50 50"
        dur="1.2s" repeatCount="indefinite"/>
    </circle>
    </svg>`);

// Переключение параметров по типу задачи
typeSelect.onchange = function () {
    document.getElementById("param-size").style.display =
        this.value === "resize" ? "block" : "none";
    document.getElementById("param-watermark").style.display =
        this.value === "watermark" ? "block" : "none";
};

// Отправка формы
document.getElementById("uploadForm").onsubmit = async (e) => {
    e.preventDefault();

    const form = new FormData(e.target);
    const res = await fetch("/upload", {
        method: "POST",
        body: form,
    });
    const data = await res.json();

    if (data.id) {
        addItem(data.id);
        pollStatus(data.id);
    }
};

// Добавление задачи в список
function addItem(id) {
    const div = document.createElement("div");
    div.className = "item";
    div.id = "item-" + id;

    div.innerHTML = `
        <a href="/image/${id}" target="_blank">
            <img class="preview" id="img-${id}" src="${placeholderSVG}" alt="preview">
        </a>
        <div style="flex-grow:1;">
            <div><b>ID:</b> ${id}</div>
            <div class="status processing" id="status-${id}">processing...</div>
        </div>
        <button onclick="deleteImage('${id}')">Удалить</button>
    `;

    list.prepend(div);
}

// Опрос статуса
async function pollStatus(id) {
    const st = document.getElementById("status-" + id);
    const img = document.getElementById("img-" + id);
    const link = document.getElementById("link-" + id);

    const interval = setInterval(async () => {
        const res = await fetch("/status/" + id);
        if (!res.ok) return;

        const data = await res.json();
        st.textContent = data.status;
        st.className = "status " + data.status;

        if (data.status === "done" && data.result) {
            const url = data.result + "?t=" + Date.now();

            img.src = url;

            link.href = url;
            link.onclick = null;

            clearInterval(interval);
        }

        if (data.status === "failed") {
            img.src = "";
            clearInterval(interval);
        }
    }, 2000);
}

// Удаление изображения
async function deleteImage(id) {
    await fetch("/image/" + id, { method: "DELETE" });
    document.getElementById("item-" + id).remove();
}
