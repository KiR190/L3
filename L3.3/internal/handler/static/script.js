const API_BASE = "";

            // State
            let offset = 0;
            let limit = 10;
            let sort = "DESC";
            let q = "";

            const el = (id) => document.getElementById(id);
            const commentContainer = el("commentContainer");
            const searchInput = el("searchInput");
            const searchBtn = el("searchBtn");
            const sortSelect = el("sortSelect");
            const prevPage = el("prevPage");
            const nextPage = el("nextPage");
            const pageInfo = el("pageInfo");
            const pageSizeEl = el("pageSize");
            const postRoot = el("postRoot");
            const rootText = el("rootText");
            const rootAuthor = el("rootAuthor");

            // Helpers
            function buildURL(path, params = {}) {
                const url = new URL((API_BASE || location.origin) + path);
                Object.entries(params).forEach(([k, v]) => {
                    if (v !== undefined && v !== null && v !== "")
                        url.searchParams.set(k, v);
                });
                return url.toString();
            }

            async function apiGetComments(parent = null) {
                // API: GET /comments?parent={id}&page=&size=&sort=&q=
                const params = {
                    parent: parent,
                    offset: offset,
                    limit: limit,
                    sort: sort,
                };
                const url = buildURL("/comments", params);
                try {
                    const res = await fetch(url);
                    if (!res.ok) throw new Error("server: " + res.status);
                    return await res.json();
                } catch (e) {
                    console.error(e);
                    showError("Не удалось загрузить комментарии: " + e.message);
                    return [];
                }
            }

            async function apiSearchComments() {
                const params = {
                    query: searchInput.value.trim(),
                    limit: limit,
                    offset: offset,
                    sort: sort,
                };
                const url = buildURL("/comments/search", params);

                try {
                    const res = await fetch(url);
                    if (!res.ok) throw new Error("server: " + res.status);
                    return await res.json();
                } catch (e) {
                    console.error(e);
                    showError("Не удалось выполнить поиск: " + e.message);
                    return [];
                }
            }

            async function apiPostComment(parent_id, author, text) {
                const url = buildURL("/comments");
                const body = {
                    parent_id: parent_id === null ? null : parent_id,
                    author,
                    text,
                };
                const res = await fetch(url, {
                    method: "POST",
                    headers: { "Content-Type": "application/json" },
                    body: JSON.stringify(body),
                });
                if (!res.ok) throw new Error("server: " + res.status);
                return await res.json();
            }

            async function apiDeleteComment(id) {
                const url = buildURL("/comments/" + id);
                const res = await fetch(url, { method: "DELETE" });
                if (!res.ok) throw new Error("server: " + res.status);
                return true;
            }

            function showError(text) {
                // lightweight ephemeral error
                const node = document.createElement("div");
                node.textContent = text;
                node.style.cssText =
                    "position:fixed;right:20px;bottom:20px;background:#ff4d6d;padding:10px;border-radius:8px;color:white;z-index:9999;box-shadow:0 8px 30px rgba(0,0,0,.6)";
                document.body.appendChild(node);
                setTimeout(() => node.remove(), 4000);
            }

            function makeAvatar(name) {
                const initials = name
                    ? name
                          .split(" ")
                          .map((s) => s[0] || "")
                          .slice(0, 2)
                          .join("")
                          .toUpperCase()
                    : "U";
                const node = document.createElement("div");
                node.className = "avatar";
                node.textContent = initials;
                return node;
            }

            function formatDate(iso) {
                try {
                    const d = new Date(iso);
                    return d.toLocaleString();
                } catch (e) {
                    return iso;
                }
            }

            function renderComments(list) {
                commentContainer.innerHTML = "";
                if (!Array.isArray(list) || list.length === 0) {
                    commentContainer.innerHTML =
                        '<div class="panel small-muted">Нет комментариев</div>';
                    return;
                }
                list.forEach((c) => {
                    const node = renderCommentNode(c);
                    commentContainer.appendChild(node);
                });
            }

            function renderCommentNode(c) {
                const wrap = document.createElement("div");
                wrap.className = "panel";
                wrap.style.marginBottom = "10px";
                const row = document.createElement("div");
                row.className = "comment";
                const left = makeAvatar(c.author || "U");
                const body = document.createElement("div");
                body.className = "body";
                const meta = document.createElement("div");
                meta.className = "meta-row";
                const author = document.createElement("strong");
                author.textContent = c.author || "Anonymous";
                const when = document.createElement("div");
                when.className = "small-muted";
                when.textContent = formatDate(
                    c.created_at || c.createdAt || ""
                );
                meta.appendChild(author);
                meta.appendChild(when);

                const text = document.createElement("div");
                text.className = "text";
                text.textContent = c.text || "";

                const actions = document.createElement("div");
                actions.className = "actions";
                const replyBtn = document.createElement("button");
                replyBtn.className = "btn ghost";
                replyBtn.textContent = "Ответить";
                const delBtn = document.createElement("button");
                delBtn.className = "btn ghost";
                delBtn.textContent = "Удалить";
                actions.appendChild(replyBtn);
                actions.appendChild(delBtn);

                body.appendChild(meta);
                body.appendChild(text);
                body.appendChild(actions);
                row.appendChild(left);
                row.appendChild(body);
                wrap.appendChild(row);

                // reply area (hidden)
                const replyArea = document.createElement("div");
                replyArea.className = "reply-form";
                replyArea.style.display = "none";
                const ta = document.createElement("textarea");
                ta.placeholder = "Написать ответ...";
                ta.style.width = "100%";
                ta.style.minHeight = "60px";
                ta.style.padding = "8px";
                ta.style.borderRadius = "8px";
                const nameIn = document.createElement("input");
                nameIn.placeholder = "Имя";
                nameIn.style.marginTop = "8px";
                nameIn.style.padding = "8px";
                nameIn.style.width = "40%";
                nameIn.style.borderRadius = "8px";
                const send = document.createElement("button");
                send.className = "btn";
                send.textContent = "Отправить";
                send.style.marginLeft = "8px";
                replyArea.appendChild(ta);
                replyArea.appendChild(nameIn);
                replyArea.appendChild(send);
                body.appendChild(replyArea);

                replyBtn.addEventListener("click", () => {
                    replyArea.style.display =
                        replyArea.style.display === "none" ? "block" : "none";
                });

                send.addEventListener("click", async () => {
                    const textVal = ta.value.trim();
                    const authorVal = nameIn.value.trim() || "Anon";
                    if (!textVal) return showError("Введите текст ответа");
                    try {
                        const created = await apiPostComment(
                            c.id,
                            authorVal,
                            textVal
                        );
                        // простая стратегия: перезагрузить текущ page
                        await loadAndRender();
                    } catch (e) {
                        showError(e.message);
                    }
                });

                delBtn.addEventListener("click", async () => {
                    if (!confirm("Удалить комментарий?")) return;
                    try {
                        await apiDeleteComment(c.id);
                        await loadAndRender();
                    } catch (e) {
                        showError("Ошибка удаления: " + e.message);
                    }
                });

                // children
                if (c.children && c.children.length) {
                    const ch = document.createElement("div");
                    ch.className = "children";
                    c.children.forEach((child) =>
                        ch.appendChild(renderCommentNode(child))
                    );
                    wrap.appendChild(ch);
                }

                return wrap;
            }

            async function loadAndRender() {
                pageInfo.textContent = Math.floor(offset / limit) + 1;
                // For tree display, we fetch top-level comments and expect server to return children nested
                const data = await apiGetComments(null);

                // If server provides additional pagination metadata, use it. For now we're assuming full tree for requested page.
                renderComments(data);
            }

            // Events
            searchBtn.addEventListener("click", async () => {
                /*q = searchInput.value.trim();
                offset = 0;
                loadAndRender();*/
                offset = 0; // сброс страницы на первую
                const data = await apiSearchComments();
                renderComments(data);
            });

            searchInput.addEventListener("keydown", (e) => {
                if (e.key === "Enter") searchBtn.click();
            });

            sortSelect.addEventListener("change", () => {
                sort = sortSelect.value;
                offset = 0;
                loadAndRender();
            });
            pageSizeEl.addEventListener("change", () => {
                limit = Number(pageSizeEl.value);
                offset = 0;
                loadAndRender();
            });
            prevPage.addEventListener("click", () => {
                if (offset >= limit) {
                    offset -= limit;
                    loadAndRender();
                }
            });
            nextPage.addEventListener("click", () => {
                offset += limit;
                loadAndRender();
            });

            postRoot.addEventListener("click", async () => {
                const text = rootText.value.trim();
                const author = rootAuthor.value.trim() || "Anon";
                if (!text) return showError("Введите текст комментария");
                try {
                    await apiPostComment(null, author, text);
                    rootText.value = "";
                    rootAuthor.value = "";
                    offset = 0;
                    await loadAndRender();
                } catch (e) {
                    showError("Ошибка отправки: " + e.message);
                }
            });

            // quick enter to search
            searchInput.addEventListener("keydown", (e) => {
                if (e.key === "Enter") searchBtn.click();
            });

            // initial load
            loadAndRender();