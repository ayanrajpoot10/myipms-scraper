// Event listener for Enter key on captcha input
document.getElementById('captcha-input').addEventListener('keypress', function(e) {
    if (e.key === 'Enter') {
        submitCaptcha();
    }
});

function showLoading(show) {
    document.getElementById('loading').style.display = show ? 'block' : 'none';
    document.getElementById('submit-btn').disabled = show;
    document.getElementById('refresh-btn').disabled = show;
    document.getElementById('captcha-input').disabled = show;
}

function showStatus(message, type = 'info') {
    const statusDiv = document.getElementById('status');
    statusDiv.innerHTML = '<div class="status ' + type + '">' + message + '</div>';
}

function refreshCaptcha() {
    showLoading(true);
    showStatus('Refreshing captcha image...', 'info');

    fetch('/refresh-captcha', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/x-www-form-urlencoded',
        }
    })
    .then(response => response.json())
    .then(data => {
        if (data.success) {
            const img = document.getElementById('captcha-image');
            img.src = '/captcha-image?t=' + Date.now();
            
            img.onload = function() {
                showLoading(false);
                showStatus('Captcha image refreshed. Please enter the text.', 'info');
                document.getElementById('captcha-input').focus();
            };
            
            img.onerror = function() {
                showLoading(false);
                showStatus('Failed to load refreshed captcha image. Please try again.', 'error');
            };
        } else {
            showLoading(false);
            showStatus('Failed to refresh captcha: ' + data.message, 'error');
        }
    })
    .catch(error => {
        showLoading(false);
        showStatus('Error refreshing captcha: ' + error.message, 'error');
    });
}

function submitCaptcha() {
    const input = document.getElementById('captcha-input');
    const captchaText = input.value.trim();
    
    if (!captchaText) {
        showStatus('Please enter the captcha text.', 'error');
        input.focus();
        return;
    }
    
    showLoading(true);
    showStatus('Validating captcha...', 'info');
    
    fetch('/submit-captcha', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/x-www-form-urlencoded',
        },
        body: 'captcha=' + encodeURIComponent(captchaText)
    })
    .then(response => response.json())
    .then(data => {
        showLoading(false);
        if (data.success) {
            showStatus('✅ Captcha validated successfully! You can close this window. Scraping will continue automatically.', 'success');
            setTimeout(() => {
                window.close();
            }, 3000);
        } else {
            showStatus('❌ ' + data.message, 'error');
            input.value = '';
            input.focus();
            setTimeout(() => {
                refreshCaptcha();
            }, 1000);
        }
    })
    .catch(error => {
        showLoading(false);
        showStatus('❌ Error submitting captcha: ' + error.message, 'error');
        input.focus();
    });
}

window.onload = function() {
    document.getElementById('captcha-input').focus();
};
