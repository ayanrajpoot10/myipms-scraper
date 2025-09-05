// DOM elements
const captchaInput = document.getElementById('captcha-input');
const submitBtn = document.getElementById('submit-btn');
const refreshBtn = document.getElementById('refresh-btn');
const captchaImage = document.getElementById('captcha-image');
const loadingDiv = document.getElementById('loading');
const statusDiv = document.getElementById('status');

// Event listeners
captchaInput.addEventListener('keypress', function(e) {
    if (e.key === 'Enter') {
        submitCaptcha();
    }
});

document.addEventListener('DOMContentLoaded', function() {
    captchaInput.focus();
});

// Utility functions
function showLoading(show) {
    loadingDiv.style.display = show ? 'block' : 'none';
    submitBtn.disabled = show;
    refreshBtn.disabled = show;
    captchaInput.disabled = show;
}

function showStatus(message, type = 'info') {
    statusDiv.innerHTML = `<div class="status ${type}">${message}</div>`;
}

async function refreshCaptcha() {
    try {
        showLoading(true);
        showStatus('Refreshing captcha image...', 'info');

        const response = await fetch('/refresh-captcha', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/x-www-form-urlencoded',
            }
        });

        const data = await response.json();

        if (data.success) {
            captchaImage.src = `/captcha-image?t=${Date.now()}`;

            captchaImage.onload = function() {
                showLoading(false);
                showStatus('Captcha image refreshed. Please enter the text.', 'info');
                captchaInput.focus();
                captchaInput.select();
            };

            captchaImage.onerror = function() {
                showLoading(false);
                showStatus('Failed to load refreshed captcha image. Please try again.', 'error');
            };
        } else {
            showLoading(false);
            showStatus(`Failed to refresh captcha: ${data.message}`, 'error');
        }
    } catch (error) {
        showLoading(false);
        showStatus(`Error refreshing captcha: ${error.message}`, 'error');
    }
}

async function submitCaptcha() {
    const captchaText = captchaInput.value.trim();

    if (!captchaText) {
        showStatus('Please enter the captcha text.', 'error');
        captchaInput.focus();
        return;
    }

    try {
        showLoading(true);
        showStatus('Validating captcha...', 'info');

        const response = await fetch('/submit-captcha', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/x-www-form-urlencoded',
            },
            body: `captcha=${encodeURIComponent(captchaText)}`
        });

        const data = await response.json();

        showLoading(false);

        if (data.success) {
            showStatus('✅ Captcha validated successfully! You can close this window. Scraping will continue automatically.', 'success');
            setTimeout(() => {
                window.close();
            }, 3000);
        } else {
            showStatus(`❌ ${data.message}`, 'error');
            captchaInput.value = '';
            captchaInput.focus();
            setTimeout(() => {
                refreshCaptcha();
            }, 1000);
        }
    } catch (error) {
        showLoading(false);
        showStatus(`❌ Error submitting captcha: ${error.message}`, 'error');
        captchaInput.focus();
    }
}
