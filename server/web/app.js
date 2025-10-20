const form = document.getElementById("audition-form");
const statusBox = document.getElementById("status");

function setStatus(message, type) {
    statusBox.textContent = message;
    statusBox.className = type || "";
}

form.addEventListener("submit", async (event) => {
    event.preventDefault();
    setStatus("Submitting...", "");

    const formData = new FormData(form);
    const payload = Object.fromEntries(formData.entries());

    try {
        const response = await fetch("/api/auditions", {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify(payload)
        });

        if (!response.ok) {
            const text = await response.text();
            throw new Error(text || "Submission failed");
        }

        const data = await response.json();
        form.reset();
        setStatus(`Audition received! Your reference ID is ${data.id}.`, "success");
    } catch (error) {
        console.error(error);
        setStatus(error.message || "Something went wrong", "error");
    }
});
