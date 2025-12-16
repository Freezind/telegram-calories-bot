import { useState, useEffect } from 'react';
import { LogTable } from './components/LogTable';
import { fetchLogs, Log } from './api/logs';
import './App.css';

function App() {
  const [logs, setLogs] = useState<Log[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    loadLogs();
  }, []);

  const loadLogs = async () => {
    try {
      setLoading(true);
      setError(null);
      const data = await fetchLogs();
      setLogs(data);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to load logs');
    } finally {
      setLoading(false);
    }
  };

  if (loading) {
    return <div className="loading">Loading logs...</div>;
  }

  if (error) {
    return <div className="error">Error: {error}</div>;
  }

  return (
    <div className="app">
      <header className="app-header">
        <h1>Calorie Log</h1>
      </header>
      <main className="app-main">
        <LogTable logs={logs} />
      </main>
    </div>
  );
}

export default App;
