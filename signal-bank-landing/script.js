// Determine API base URL
let API_BASE;
if (window.location.hostname === "localhost" || window.location.hostname === "127.0.0.1") {
  API_BASE = `http://localhost:${window.location.port || 8081}/api`;
} else {
  API_BASE = `${window.location.origin}/api`;
}

console.log("API Base URL:", API_BASE);

// Load stats on page load
document.addEventListener("DOMContentLoaded", async () => {
  await loadStats();
  await loadLedger();
  await loadRequests();
  setupFormHandlers();
});

// Load and display stats
async function loadStats() {
  try {
    const response = await fetch(`${API_BASE}/stats`);
    const stats = await response.json();
    updateStatsDisplay(stats);
  } catch (error) {
    console.log("Stats service not available yet");
  }
}

function updateStatsDisplay(stats) {
  const elements = document.querySelectorAll("[data-stat]");
  elements.forEach((el) => {
    const stat = el.getAttribute("data-stat");
    if (stat === "balance") {
      el.textContent = `$${stats.totalBalance.toLocaleString()}`;
    } else if (stat === "distributed") {
      el.textContent = `${stats.distributedPercent}%`;
    } else if (stat === "stories") {
      el.textContent = stats.storiesFunded;
    } else if (stat === "contributors") {
      el.textContent = stats.totalContributors.toLocaleString();
    }
  });
}

// Load and display real contributions ledger
async function loadLedger() {
  try {
    const response = await fetch(`${API_BASE}/contribute`);
    const contributions = await response.json();
    
    const ledgerDiv = document.getElementById("live-ledger");
    if (!ledgerDiv) return;

    if (!contributions || contributions.length === 0) {
      ledgerDiv.innerHTML = '<div class="ledger-loading">No contributions yet. Be the first!</div>';
      return;
    }

    // Sort by most recent first
    contributions.sort((a, b) => new Date(b.timestamp) - new Date(a.timestamp));

    let ledgerHTML = '';
    contributions.slice(0, 20).forEach((contrib) => {
      const date = new Date(contrib.timestamp).toLocaleDateString('en-US', { 
        month: 'short', 
        day: 'numeric' 
      });
      const email = contrib.email.split('@')[0]; // Show just the username part
      const message = contrib.message ? ` · "${contrib.message}"` : '';
      
      ledgerHTML += `
        <div class="ledger-row">
          <span>${date} · Contributor: @${email}</span>
          <span>$${contrib.amount.toLocaleString()} added to vault${message}</span>
        </div>
      `;
    });

    ledgerDiv.innerHTML = ledgerHTML;
  } catch (error) {
    console.error("Error loading ledger:", error);
    const ledgerDiv = document.getElementById("live-ledger");
    if (ledgerDiv) {
      ledgerDiv.innerHTML = '<div class="ledger-loading">Loading live data...</div>';
    }
  }
}

// Load and display real help requests
async function loadRequests() {
  try {
    const response = await fetch(`${API_BASE}/requests`);
    const requests = await response.json();
    
    const requestsDiv = document.getElementById("requests-display");
    if (!requestsDiv) return;

    if (!requests || requests.length === 0) {
      requestsDiv.innerHTML = '<div class="loading">No active requests yet. Be the first to share your story!</div>';
      return;
    }

    // Sort by most recent first
    requests.sort((a, b) => new Date(b.timestamp) - new Date(a.timestamp));

    let requestsHTML = '';
    requests.slice(0, 9).forEach((req) => {
      const date = new Date(req.timestamp).toLocaleDateString('en-US', { 
        month: 'short', 
        day: 'numeric' 
      });
      const videoLink = req.videoUrl ? `<a href="${req.videoUrl}" target="_blank" rel="noopener">Watch story →</a>` : '';
      
      requestsHTML += `
        <div class="request-card">
          <h4>${req.name}</h4>
          <p class="story">${req.story.substring(0, 120)}${req.story.length > 120 ? '...' : ''}</p>
          <p class="amount">$${req.amount.toLocaleString()} needed</p>
          <div class="meta">
            <span>Posted ${date}</span>
            <span id="votes-${req.id}">👍 0</span>
          </div>
          ${videoLink}
          <button class="vote-btn" onclick="vote('${req.id}')">Vote to support</button>
        </div>
      `;
    });

    requestsDiv.innerHTML = requestsHTML;
  } catch (error) {
    console.error("Error loading requests:", error);
    const requestsDiv = document.getElementById("requests-display");
    if (requestsDiv) {
      requestsDiv.innerHTML = '<div class="loading">Loading requests...</div>';
    }
  }
}

// Form handlers
function setupFormHandlers() {
  // Contribution form
  const contributeForm = document.getElementById("contribute-form");
  if (contributeForm) {
    contributeForm.addEventListener("submit", handleContribution);
  }

  // Request form
  const requestForm = document.getElementById("request-form");
  if (requestForm) {
    requestForm.addEventListener("submit", handleRequestSubmit);
  }

  // Subscribe form
  const subscribeForm = document.querySelector(".signup");
  if (subscribeForm) {
    subscribeForm.addEventListener("submit", handleSubscribe);
  }
}

async function handleContribution(e) {
  e.preventDefault();

  const form = e.target;
  const formData = new FormData(form);

  const contribution = {
    email: formData.get("email"),
    amount: parseFloat(formData.get("amount")),
    message: formData.get("message") || "",
  };

  if (contribution.amount <= 0) {
    alert("Please enter a valid amount");
    return;
  }

  try {
    const response = await fetch(`${API_BASE}/contribute`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(contribution),
    });

    if (response.ok) {
      alert("Thank you for your contribution! 🎉");
      form.reset();
      await loadStats();
    } else {
      alert("Something went wrong. Please try again.");
    }
  } catch (error) {
    console.error("Error:", error);
    alert("Error submitting contribution");
  }
}

async function handleRequestSubmit(e) {
  e.preventDefault();

  const form = e.target;
  const formData = new FormData(form);

  const request = {
    name: formData.get("name"),
    email: formData.get("email"),
    story: formData.get("story"),
    videoUrl: formData.get("videoUrl"),
    amount: parseFloat(formData.get("amount")) || 0,
  };

  if (!request.name || !request.email || !request.story) {
    alert("Please fill in all required fields");
    return;
  }

  try {
    const response = await fetch(`${API_BASE}/requests`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(request),
    });

    if (response.ok) {
      alert("Your request has been submitted! Our team will review it shortly.");
      form.reset();
      await loadStats();
    } else {
      alert("Something went wrong. Please try again.");
    }
  } catch (error) {
    console.error("Error:", error);
    alert("Error submitting request");
  }
}

async function handleSubscribe(e) {
  e.preventDefault();

  const email = e.target.querySelector('input[type="email"]').value;

  if (!email) {
    alert("Please enter an email");
    return;
  }

  try {
    const response = await fetch(`${API_BASE}/subscribe`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ email }),
    });

    if (response.ok) {
      alert("Successfully subscribed! ✨");
      e.target.reset();
    } else {
      alert("Subscription failed. Please try again.");
    }
  } catch (error) {
    console.error("Error:", error);
    alert("Error subscribing");
  }
}

// Vote handler
async function vote(requestId) {
  try {
    const response = await fetch(`${API_BASE}/vote`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ requestId }),
    });

    if (response.ok) {
      alert("Vote recorded! Thank you.");
      await loadStats();
    }
  } catch (error) {
    console.error("Error voting:", error);
  }
}

// Refresh stats and ledger periodically
setInterval(async () => {
  await loadStats();
  await loadLedger();
  await loadRequests();
}, 30000); // Refresh every 30 seconds

