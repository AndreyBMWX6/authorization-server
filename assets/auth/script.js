const allowBtn = document.getElementById('allow-btn');
const denyBtn = document.getElementById('deny-btn');

allowBtn.addEventListener('click', async () => {
    try {
        let redirectUri = 'http://localhost:7000/authorization/auth';

        const params = getCookie("params")
        if (params !== "") {
            redirectUri += `?${params}`
        }

        const response = await fetch(redirectUri, {
            credentials: "same-origin",
        })
        if (!response.ok) {
            console.error("get authorization code failed")
        }

        const data = await response.json()
        authCode = data["code"]
        state = data["state"]

        const tokenUri = `http://localhost:8000/token?code=${authCode}&state=${state}`
        window.location.href = tokenUri;
    } catch (error) {
        console.error('Error getting authorization code:', error);
    }
});

denyBtn.addEventListener('click', () => {
    try {
        modal.remove();
        alert('Access denied');
    } catch (error) {
        console.error('Error handling deny response:', error);
    }
});

function getCookie(name) {
    const value = `; ${document.cookie}`;
    const parts = value.split(`; ${name}=`);
    if (parts.length > 1) return parts[1].split(';').shift();
}
