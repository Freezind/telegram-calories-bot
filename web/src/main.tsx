import React from 'react';
import ReactDOM from 'react-dom/client';
import App from './App';
import './index.css';

// Initialize Telegram WebApp
if (window.Telegram?.WebApp) {
  const initData = window.Telegram.WebApp.initData;
  console.log('Telegram WebApp SDK loaded');
  console.log('initData available:', !!initData);
  console.log('initData length:', initData?.length || 0);

  // Safe preview: log first 50 chars only (contains query_id and partial user data)
  if (initData) {
    console.log('initData preview:', initData.substring(0, 50) + '...');
  } else {
    console.warn('⚠️ initData is EMPTY - app will not work. Make sure to open via Telegram WebApp button.');
  }

  window.Telegram.WebApp.ready();
  window.Telegram.WebApp.expand();
} else {
  console.error('❌ Telegram WebApp SDK not found!');
}

ReactDOM.createRoot(document.getElementById('root')!).render(
  <React.StrictMode>
    <App />
  </React.StrictMode>
);
