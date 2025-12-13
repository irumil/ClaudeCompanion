// Load saved settings
function loadSettings() {
  browser.storage.local.get(['port']).then((result) => {
    if (result.port) {
      document.getElementById('port').value = result.port;
    }
  });
}

// Save settings
function saveSettings(e) {
  e.preventDefault();

  const port = parseInt(document.getElementById('port').value);

  if (port < 1024 || port > 65535) {
    showStatus('Порт должен быть в диапазоне 1024-65535', 'error');
    return;
  }

  browser.storage.local.set({ port: port }).then(() => {
    showStatus('Настройки сохранены успешно!', 'success');
  }).catch((error) => {
    showStatus('Ошибка сохранения: ' + error.message, 'error');
  });
}

// Test connection to desktop app
function testConnection() {
  const port = parseInt(document.getElementById('port').value);
  const endpoint = `http://127.0.0.1:${port}/health`;

  showStatus('Проверка соединения...', 'info');

  fetch(endpoint)
    .then(response => response.json())
    .then(data => {
      showStatus(`✓ Соединение установлено! Версия: ${data.version || 'unknown'}`, 'success');
    })
    .catch(error => {
      showStatus('✗ Не удалось подключиться к десктопному приложению. Убедитесь, что оно запущено.', 'error');
      console.error('Connection test failed:', error);
    });
}

// Show status message
function showStatus(message, type) {
  const statusEl = document.getElementById('status');
  statusEl.textContent = message;
  statusEl.className = 'status ' + type;
  statusEl.style.display = 'block';

  // Auto-hide after 5 seconds for success messages
  if (type === 'success') {
    setTimeout(() => {
      statusEl.style.display = 'none';
    }, 5000);
  }
}

// Event listeners
document.addEventListener('DOMContentLoaded', loadSettings);
document.getElementById('settingsForm').addEventListener('submit', saveSettings);
document.getElementById('testBtn').addEventListener('click', testConnection);
