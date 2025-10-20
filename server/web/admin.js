const formEl = document.getElementById("admin-form");
const statusEl = document.getElementById("admin-status");
const resultsEl = document.getElementById("results");

function formatDate(isoString) {
    const date = new Date(isoString);
    if (Number.isNaN(date.getTime())) {
        return isoString;
    }
    return date.toLocaleString();
}

function renderSubmissions(items) {
    resultsEl.innerHTML = "";
    if (!items.length) {
        resultsEl.innerHTML = "<p>No submissions found for that filter.</p>";
        return;
    }

    items.forEach((item) => {
        const container = document.createElement("article");
        container.className = "submission";

        const title = document.createElement("h2");
        title.textContent = `${item.name} · ${item.country}`;
        container.appendChild(title);

        const meta = document.createElement("div");
        meta.className = "meta";
        const parts = [];
        if (item.socialHandle) {
            parts.push(item.socialHandle);
        }
        parts.push(`Submitted ${formatDate(item.createdAt)}`);
        meta.textContent = parts.join(" · ");
        container.appendChild(meta);

        const link = document.createElement("a");
        link.href = item.videoUrl;
        link.target = "_blank";
        link.rel = "noopener";
        link.textContent = "Watch audition video";
        container.appendChild(link);

        if (item.message) {
            const msg = document.createElement("p");
            msg.className = "message";
            msg.textContent = item.message;
            container.appendChild(msg);
        }

        resultsEl.appendChild(container);
    });
}

formEl.addEventListener("submit", async (event) => {
    event.preventDefault();
    statusEl.textContent = "Loading submissions...";
    resultsEl.innerHTML = "";

    const formData = new FormData(formEl);
    const params = new URLSearchParams();

    const token = formData.get("token")?.trim();
    if (!token) {
        statusEl.textContent = "Admin token is required.";
        return;
    }
    params.set("token", token);

    const country = formData.get("country")?.trim();
    if (country) {
        params.set("country", country);
    }

    const limit = formData.get("limit")?.trim();
    if (limit) {
        params.set("limit", limit);
    }

    try {
        const response = await fetch(`/api/auditions?${params.toString()}`);
        if (!response.ok) {
            const text = await response.text();
            throw new Error(text || `Request failed with ${response.status}`);
        }
        const data = await response.json();
        renderSubmissions(data);
        statusEl.textContent = `Loaded ${data.length} submission(s).`;
    } catch (error) {
        console.error(error);
        statusEl.textContent = error.message || "Failed to load submissions.";
    }
});
