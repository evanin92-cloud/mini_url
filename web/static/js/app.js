document.addEventListener('DOMContentLoaded', function() {
    const form = document.getElementById('shorten-form');
    const result = document.getElementById('result');
    const shortUrl = document.getElementById('short-url');
    const copyBtn = document.getElementById('copy-btn');

    if (form) {
        form.addEventListener('submit', async function(e) {
            e.preventDefault();

            const urlInput = document.getElementById('original-url');
            const url = urlInput.value.trim();

            if (!url) {
                alert('Введите URL-адрес');
                return;
            }

            try {
                const response = await fetch('/api/v1/shorten', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify({ url: url })
                });

                if (!response.ok) {
                    throw new Error('Ошибка сервера');
                }

                const data = await response.json();
                shortUrl.textContent = data.short_url;
                result.classList.remove('hidden');
                urlInput.value = '';
            } catch (error) {
                alert('Ошибка при создании ссылки: ' + error.message);
            }
        });
    }

    if (copyBtn) {
        copyBtn.addEventListener('click', function() {
            navigator.clipboard.writeText(shortUrl.textContent).then(() => {
                copyBtn.textContent = '✅ Скопировано!';
                setTimeout(() => {
                    copyBtn.textContent = 'Копировать';
                }, 2000);
            });
        });
    }
});