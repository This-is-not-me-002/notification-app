package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

type Notification struct {
	ID        string    `json:"id"`
	Message   string    `json:"message"`
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
}

type NotificationStore struct {
	mu            sync.RWMutex
	notifications []Notification
}

var store = &NotificationStore{
	notifications: make([]Notification, 0),
}

func main() {
	http.HandleFunc("/", serveHTML)
	http.HandleFunc("/manifest.json", serveManifest)
	http.HandleFunc("/service-worker.js", serveServiceWorker)
	http.HandleFunc("/api/send-notification", sendNotificationHandler)
	http.HandleFunc("/api/notifications", getNotificationsHandler)
	
	port := "8080"
	fmt.Printf("\n========================================\n")
	fmt.Printf("🚀 Notification Server Started!\n")
	fmt.Printf("========================================\n")
	fmt.Printf("📱 Open in browser: http://localhost:%s\n", port)
	fmt.Printf("📱 Or use your IP: http://YOUR_IP:%s\n", port)
	fmt.Printf("\n💡 Setup Instructions:\n")
	fmt.Printf("1. Open the URL in your mobile browser\n")
	fmt.Printf("2. Click 'Install App' or 'Add to Home Screen'\n")
	fmt.Printf("3. Allow notifications when prompted\n")
	fmt.Printf("4. Notifications will appear in your notification bar!\n")
	fmt.Printf("========================================\n\n")
	
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

func serveHTML(w http.ResponseWriter, r *http.Request) {
	html := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <meta name="theme-color" content="#667eea">
    <title>Push Notification System</title>
    <link rel="manifest" href="/manifest.json">
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }
        
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, sans-serif;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            min-height: 100vh;
            display: flex;
            justify-content: center;
            align-items: center;
            padding: 20px;
        }
        
        .container {
            background: white;
            border-radius: 20px;
            padding: 40px;
            box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
            max-width: 600px;
            width: 100%;
        }
        
        h1 {
            color: #333;
            margin-bottom: 10px;
            font-size: 28px;
        }
        
        .subtitle {
            color: #666;
            margin-bottom: 30px;
            font-size: 14px;
        }
        
        .alert {
            padding: 15px;
            border-radius: 10px;
            margin-bottom: 20px;
            font-size: 14px;
        }
        
        .alert.info {
            background: #d1ecf1;
            color: #0c5460;
            border: 1px solid #bee5eb;
        }
        
        .alert.success {
            background: #d4edda;
            color: #155724;
            border: 1px solid #c3e6cb;
        }
        
        .alert.warning {
            background: #fff3cd;
            color: #856404;
            border: 1px solid #ffeaa7;
        }
        
        .alert.error {
            background: #f8d7da;
            color: #721c24;
            border: 1px solid #f5c6cb;
        }
        
        .button-container {
            margin-bottom: 30px;
        }
        
        button {
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: white;
            border: none;
            padding: 15px 40px;
            border-radius: 50px;
            font-size: 16px;
            font-weight: 600;
            cursor: pointer;
            transition: all 0.3s ease;
            box-shadow: 0 4px 15px rgba(102, 126, 234, 0.4);
            width: 100%;
            margin-bottom: 10px;
        }
        
        button:hover {
            transform: translateY(-2px);
            box-shadow: 0 6px 20px rgba(102, 126, 234, 0.6);
        }
        
        button:active {
            transform: translateY(0);
        }
        
        button:disabled {
            background: #ccc;
            cursor: not-allowed;
            box-shadow: none;
        }
        
        .status-section {
            margin-top: 30px;
        }
        
        .status-section h2 {
            color: #333;
            font-size: 20px;
            margin-bottom: 15px;
        }
        
        .notifications-list {
            max-height: 400px;
            overflow-y: auto;
        }
        
        .notification-item {
            background: #f8f9fa;
            border-left: 4px solid #667eea;
            padding: 15px;
            margin-bottom: 10px;
            border-radius: 8px;
            animation: slideIn 0.3s ease;
        }
        
        @keyframes slideIn {
            from {
                opacity: 0;
                transform: translateY(-10px);
            }
            to {
                opacity: 1;
                transform: translateY(0);
            }
        }
        
        .notification-item.success {
            border-left-color: #28a745;
        }
        
        .notification-message {
            font-weight: 600;
            color: #333;
            margin-bottom: 5px;
        }
        
        .notification-meta {
            font-size: 12px;
            color: #666;
            display: flex;
            justify-content: space-between;
            align-items: center;
        }
        
        .status-badge {
            display: inline-block;
            padding: 4px 12px;
            border-radius: 12px;
            font-size: 11px;
            font-weight: 600;
            text-transform: uppercase;
            background: #d4edda;
            color: #155724;
        }
        
        .empty-state {
            text-align: center;
            color: #999;
            padding: 40px;
            font-style: italic;
        }
        
        .spinner {
            display: inline-block;
            width: 16px;
            height: 16px;
            border: 2px solid rgba(255, 255, 255, 0.3);
            border-radius: 50%;
            border-top-color: white;
            animation: spin 0.6s linear infinite;
            margin-right: 8px;
            vertical-align: middle;
        }
        
        @keyframes spin {
            to { transform: rotate(360deg); }
        }
        
        .install-btn {
            background: linear-gradient(135deg, #28a745 0%, #20c997 100%);
            display: none;
        }
        
        .permission-status {
            padding: 10px;
            border-radius: 8px;
            text-align: center;
            margin-bottom: 15px;
            font-size: 13px;
            font-weight: 600;
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>🔔 Push Notification System</h1>
        <p class="subtitle">Real device notifications in your notification bar!</p>
        
        <div id="status-alerts"></div>
        
        <div id="permission-status" class="permission-status"></div>
        
        <div class="button-container">
            <button id="installBtn" class="install-btn" onclick="installPWA()">
                📱 Install App
            </button>
            <button id="enableBtn" onclick="enableNotifications()">
                🔔 Enable Notifications
            </button>
            <button id="sendBtn" onclick="sendNotification()" disabled>
                📨 Send Test Notification
            </button>
        </div>
        
        <div class="status-section">
            <h2>Notification History</h2>
            <div id="notifications" class="notifications-list">
                <div class="empty-state">No notifications sent yet</div>
            </div>
        </div>
    </div>

    <script>
        let deferredPrompt;
        let swRegistration = null;
        
        window.addEventListener('beforeinstallprompt', (e) => {
            e.preventDefault();
            deferredPrompt = e;
            document.getElementById('installBtn').style.display = 'block';
        });
        
        function installPWA() {
            if (deferredPrompt) {
                deferredPrompt.prompt();
                deferredPrompt.userChoice.then((choiceResult) => {
                    if (choiceResult.outcome === 'accepted') {
                        showAlert('App installed! Now enable notifications.', 'success');
                    }
                    deferredPrompt = null;
                    document.getElementById('installBtn').style.display = 'none';
                });
            }
        }
        
        async function registerServiceWorker() {
            if ('serviceWorker' in navigator) {
                try {
                    swRegistration = await navigator.serviceWorker.register('/service-worker.js');
                    console.log('Service Worker registered:', swRegistration);
                    updatePermissionStatus();
                    return true;
                } catch (error) {
                    console.error('Service Worker registration failed:', error);
                    showAlert('Failed to register service worker: ' + error.message, 'error');
                    return false;
                }
            } else {
                showAlert('Service Workers not supported in this browser', 'error');
                return false;
            }
        }
        
        async function enableNotifications() {
            const btn = document.getElementById('enableBtn');
            btn.disabled = true;
            
            if (!('Notification' in window)) {
                showAlert('This browser does not support notifications', 'error');
                btn.disabled = false;
                return;
            }
            
            if (!swRegistration) {
                const registered = await registerServiceWorker();
                if (!registered) {
                    btn.disabled = false;
                    return;
                }
            }
            
            try {
                const permission = await Notification.requestPermission();
                
                if (permission === 'granted') {
                    showAlert('✅ Notifications enabled! You will receive alerts in your notification bar.', 'success');
                    document.getElementById('sendBtn').disabled = false;
                    updatePermissionStatus();
                } else if (permission === 'denied') {
                    showAlert('❌ Notification permission denied. Please enable in browser settings.', 'error');
                } else {
                    showAlert('⚠️ Notification permission dismissed', 'warning');
                }
            } catch (error) {
                showAlert('Error requesting permission: ' + error.message, 'error');
            }
            
            btn.disabled = false;
        }
        
        function updatePermissionStatus() {
            const statusDiv = document.getElementById('permission-status');
            const sendBtn = document.getElementById('sendBtn');
            
            if (!('Notification' in window)) {
                statusDiv.innerHTML = '❌ Notifications not supported';
                statusDiv.style.background = '#f8d7da';
                statusDiv.style.color = '#721c24';
                return;
            }
            
            const permission = Notification.permission;
            
            if (permission === 'granted') {
                statusDiv.innerHTML = '✅ Notifications Enabled';
                statusDiv.style.background = '#d4edda';
                statusDiv.style.color = '#155724';
                sendBtn.disabled = false;
            } else if (permission === 'denied') {
                statusDiv.innerHTML = '❌ Notifications Blocked';
                statusDiv.style.background = '#f8d7da';
                statusDiv.style.color = '#721c24';
                sendBtn.disabled = true;
            } else {
                statusDiv.innerHTML = '⚠️ Notifications Not Enabled';
                statusDiv.style.background = '#fff3cd';
                statusDiv.style.color = '#856404';
                sendBtn.disabled = true;
            }
        }
        
        async function sendNotification() {
            const btn = document.getElementById('sendBtn');
            btn.disabled = true;
            btn.innerHTML = '<span class="spinner"></span>Sending...';
            
            try {
                const response = await fetch('/api/send-notification', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    }
                });
                
                const data = await response.json();
                
                if (response.ok) {
                    if (swRegistration && Notification.permission === 'granted') {
                        swRegistration.showNotification('Test Notification', {
                            body: 'This is a test notification from your app!',
                            icon: 'https://cdn-icons-png.flaticon.com/512/2645/2645897.png',
                            badge: 'https://cdn-icons-png.flaticon.com/512/2645/2645897.png',
                            vibrate: [200, 100, 200],
                            tag: 'test-notification',
                            requireInteraction: false
                        });
                    }
                    
                    showAlert('✅ Notification sent to your device!', 'success');
                    loadNotifications();
                } else {
                    showAlert('Failed to send notification: ' + data.error, 'error');
                }
            } catch (error) {
                showAlert('Error: ' + error.message, 'error');
            } finally {
                btn.disabled = false;
                btn.innerHTML = '📨 Send Test Notification';
            }
        }
        
        async function loadNotifications() {
            try {
                const response = await fetch('/api/notifications');
                const data = await response.json();
                
                const container = document.getElementById('notifications');
                
                if (data.notifications && data.notifications.length > 0) {
                    const notifHTML = data.notifications.map(function(notif) {
                        const timestamp = new Date(notif.timestamp).toLocaleString();
                        return '<div class="notification-item success">' +
                            '<div class="notification-message">' + notif.message + '</div>' +
                            '<div class="notification-meta">' +
                            '<span>' + timestamp + '</span>' +
                            '<span class="status-badge">' + notif.status + '</span>' +
                            '</div></div>';
                    }).reverse().join('');
                    container.innerHTML = notifHTML;
                } else {
                    container.innerHTML = '<div class="empty-state">No notifications sent yet</div>';
                }
            } catch (error) {
                console.error('Error loading notifications:', error);
            }
        }
        
        function showAlert(message, type) {
            const alertsDiv = document.getElementById('status-alerts');
            const alert = document.createElement('div');
            alert.className = 'alert ' + type;
            alert.textContent = message;
            alertsDiv.appendChild(alert);
            
            setTimeout(() => {
                alert.style.animation = 'slideIn 0.3s ease reverse';
                setTimeout(() => alert.remove(), 300);
            }, 5000);
        }
        
        registerServiceWorker();
        updatePermissionStatus();
        loadNotifications();
        setInterval(loadNotifications, 5000);
    </script>
</body>
</html>`
	
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(html))
}

func serveManifest(w http.ResponseWriter, r *http.Request) {
	manifest := `{
  "name": "Push Notification System",
  "short_name": "Notifications",
  "description": "Real-time push notification system",
  "start_url": "/",
  "display": "standalone",
  "background_color": "#667eea",
  "theme_color": "#667eea",
  "orientation": "portrait",
  "icons": [
    {
      "src": "https://cdn-icons-png.flaticon.com/512/2645/2645897.png",
      "sizes": "512x512",
      "type": "image/png",
      "purpose": "any maskable"
    }
  ]
}`
	
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(manifest))
}

func serveServiceWorker(w http.ResponseWriter, r *http.Request) {
	sw := `self.addEventListener('install', (event) => {
  console.log('Service Worker installing...');
  self.skipWaiting();
});

self.addEventListener('activate', (event) => {
  console.log('Service Worker activating...');
  event.waitUntil(clients.claim());
});

self.addEventListener('notificationclick', (event) => {
  event.notification.close();
  event.waitUntil(
    clients.openWindow('/')
  );
});`
	
	w.Header().Set("Content-Type", "application/javascript")
	w.Write([]byte(sw))
}

func sendNotificationHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	notification := Notification{
		ID:        fmt.Sprintf("notif-%d", time.Now().Unix()),
		Message:   "Test Notification",
		Status:    "sent",
		Timestamp: time.Now(),
	}
	
	store.mu.Lock()
	store.notifications = append(store.notifications, notification)
	store.mu.Unlock()
	
	log.Printf("✅ Notification sent: %s at %s\n", notification.Message, notification.Timestamp.Format(time.RFC3339))
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":      true,
		"notification": notification,
		"message":      "Notification sent successfully",
	})
}

func getNotificationsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	store.mu.RLock()
	notifications := store.notifications
	store.mu.RUnlock()
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"notifications": notifications,
		"count":         len(notifications),
	})
}
