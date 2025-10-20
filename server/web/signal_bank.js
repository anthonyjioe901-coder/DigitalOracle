const form = document.getElementById("contribution-form");
const statusEl = document.getElementById("contribution-status");
const ledgerEl = document.getElementById("ledger");
const totalDisplay = document.getElementById("total-display");
const refreshBtn = document.getElementById("refresh-ledger");

function setStatus(message, className = "") {
    statusEl.textContent = message;
    statusEl.className = className;
}

function formatCurrency(amount) {
    const formatter = new Intl.NumberFormat(undefined, {
        style: "currency",
        currency: "USD",
        minimumFractionDigits: 2,
    });
    return formatter.format(amount);
}

function formatDate(dateString) {
    const date = new Date(dateString);
    if (Number.isNaN(date.getTime())) {
        return dateString;
    }
    return date.toLocaleString();
}

function renderLedger(entries) {
    ledgerEl.innerHTML = "";
    let total = 0;

    if (!entries.length) {
        ledgerEl.innerHTML = "<p>No contributions recorded yet.</p>";
        totalDisplay.textContent = "Total: $0.00";
        return;
    }

    entries.forEach((entry) => {
        total += Number(entry.amount || 0);

        const container = document.createElement("article");
        container.className = "entry";

        const header = document.createElement("div");
        header.className = "entry-header";

        const name = document.createElement("strong");
        name.textContent = entry.name || "Anonymous";
        header.appendChild(name);

        const amount = document.createElement("span");
        amount.textContent = formatCurrency(Number(entry.amount || 0));
        header.appendChild(amount);

        container.appendChild(header);

        const timestamp = document.createElement("div");
        timestamp.className = "meta";
        timestamp.textContent = `Logged: ${formatDate(entry.createdAt)}`;
        container.appendChild(timestamp);

        if (entry.message) {
            const message = document.createElement("p");
            message.className = "entry-message";
            message.textContent = entry.message;
            container.appendChild(message);
        }

        ledgerEl.appendChild(container);
    });

    totalDisplay.textContent = `Total: ${formatCurrency(total)}`;
}

async function loadLedger() {
    try {
        const response = await fetch("/api/signal-bank/contributions");
        if (!response.ok) {
            throw new Error(`Failed to load ledger (${response.status})`);
        }
        const data = await response.json();
        renderLedger(data);
    } catch (error) {
        console.error(error);
        setStatus(error.message || "Unable to load ledger.", "error");
    }
}

form.addEventListener("submit", async (event) => {
    event.preventDefault();
    setStatus("Submitting contribution...", "");

    const formData = new FormData(form);
    const amount = parseFloat(formData.get("amount"));

    if (Number.isNaN(amount) || amount <= 0) {
        setStatus("Amount must be greater than zero.", "error");
        return;
    }

    const payload = {
        name: formData.get("name")?.trim() || "",
        amount,
        message: formData.get("message")?.trim() || "",
    };

    try {
        const response = await fetch("/api/signal-bank/contributions", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify(payload),
        });

        if (!response.ok) {
            const text = await response.text();
            throw new Error(text || "Failed to record contribution");
        }

        form.reset();
        setStatus("Contribution logged. Thank you!", "success");
        await loadLedger();
    } catch (error) {
        console.error(error);
        setStatus(error.message || "Failed to record contribution.", "error");
    }
});

refreshBtn.addEventListener("click", () => {
    setStatus("Ledger refreshed.", "success");
    loadLedger();
});

loadLedger();
