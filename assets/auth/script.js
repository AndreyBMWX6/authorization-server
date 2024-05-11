const allowBtn = document.getElementById('allow-btn');
const denyBtn = document.getElementById('deny-btn');

allowBtn.addEventListener('click', () => {
    try {
        const redirectUri = 'http://localhost:7000/authorization/auth';
        // const authorizationCode = crypto.getRandomValues(new Uint32Array(1))[0].toString(36).slice(2);
        const redirectUrl = new URL(redirectUri);
        const url = document.URL
        redirectUrl.searchParams.set('url', url);
        redirectUrl.searchParams.set('code', authorizationCode);
        window.location.href = redirectUrl.href;
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
