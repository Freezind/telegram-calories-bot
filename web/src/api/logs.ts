// Telegram WebApp types (from @twa-dev/sdk)
declare global {
  interface Window {
    Telegram?: {
      WebApp: {
        initData: string;
        ready: () => void;
        expand: () => void;
      };
    };
  }
}

export interface Log {
  id: string;
  userId: number;
  foodItems: string[];
  calories: number;
  confidence: 'high' | 'medium' | 'low';
  timestamp: string;
  createdAt: string;
  updatedAt: string;
}

export interface LogCreate {
  foodItems: string[];
  calories: number;
  confidence: 'high' | 'medium' | 'low';
  timestamp?: string;
}

export interface LogUpdate {
  foodItems?: string[];
  calories?: number;
  confidence?: 'high' | 'medium' | 'low';
  timestamp?: string;
}

// Get API base URL from environment variable
const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || '';

// Get initData from Telegram WebApp
function getInitData(): string {
  return window.Telegram?.WebApp?.initData || '';
}

// Fetch all logs for the current user
export async function fetchLogs(): Promise<Log[]> {
  try {
    const response = await fetch(`${API_BASE_URL}/api/logs`, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
        'X-Telegram-Init-Data': getInitData(),
      },
    });

    if (!response.ok) {
      const errorText = await response.text();
      throw new Error(`Failed to fetch logs (${response.status}): ${errorText || response.statusText}`);
    }

    return response.json();
  } catch (error) {
    if (error instanceof Error) {
      // Check if it's a network error
      if (error.message.includes('Failed to fetch') && !error.message.includes('(')) {
        throw new Error('Failed to fetch logs: Cannot connect to backend server. Make sure the backend is running on http://localhost:8080');
      }
      throw error;
    }
    throw new Error('Failed to fetch logs: Unknown error');
  }
}

// Create a new log entry
export async function createLog(log: LogCreate): Promise<Log> {
  const response = await fetch(`${API_BASE_URL}/api/logs`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'X-Telegram-Init-Data': getInitData(),
    },
    body: JSON.stringify(log),
  });

  if (!response.ok) {
    throw new Error(`Failed to create log: ${response.statusText}`);
  }

  return response.json();
}

// Update an existing log entry
export async function updateLog(id: string, update: LogUpdate): Promise<Log> {
  const response = await fetch(`${API_BASE_URL}/api/logs/${id}`, {
    method: 'PATCH',
    headers: {
      'Content-Type': 'application/json',
      'X-Telegram-Init-Data': getInitData(),
    },
    body: JSON.stringify(update),
  });

  if (!response.ok) {
    throw new Error(`Failed to update log: ${response.statusText}`);
  }

  return response.json();
}

// Delete a log entry
export async function deleteLog(id: string): Promise<void> {
  const response = await fetch(`${API_BASE_URL}/api/logs/${id}`, {
    method: 'DELETE',
    headers: {
      'X-Telegram-Init-Data': getInitData(),
    },
  });

  if (!response.ok) {
    throw new Error(`Failed to delete log: ${response.statusText}`);
  }
}
