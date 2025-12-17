import { useState, useEffect } from 'react';
import { LogTable } from './components/LogTable';
import { LogForm } from './components/LogForm';
import { DeleteConfirm } from './components/DeleteConfirm';
import { fetchLogs, deleteLog, Log } from './api/logs';
import './App.css';

function App() {
  const [logs, setLogs] = useState<Log[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [isFormOpen, setIsFormOpen] = useState(false);
  const [formMode, setFormMode] = useState<'create' | 'edit'>('create');
  const [editingLog, setEditingLog] = useState<Log | null>(null);
  const [isDeleteConfirmOpen, setIsDeleteConfirmOpen] = useState(false);
  const [deletingLog, setDeletingLog] = useState<Log | null>(null);

  useEffect(() => {
    // Note: We no longer block if initData is missing, to support dev mode with DEV_FAKE_USER_ID
    // The backend will handle auth via DEV_FAKE_USER_ID or X-Telegram-Init-Data header
    const hasInitData = !!window.Telegram?.WebApp?.initData;
    console.log(`[Frontend] Telegram WebApp initData ${hasInitData ? 'present' : 'missing'} (length: ${window.Telegram?.WebApp?.initData?.length || 0})`);

    if (!hasInitData) {
      console.warn('[Frontend] No initData - assuming dev mode with DEV_FAKE_USER_ID on backend');
    }

    loadLogs();
  }, []);

  const loadLogs = async () => {
    try {
      setLoading(true);
      setError(null);
      const data = await fetchLogs();
      setLogs(data);
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Failed to load logs';
      if (errorMessage.includes('401') || errorMessage.toLowerCase().includes('unauthorized')) {
        setError(errorMessage);
      } else {
        setError(errorMessage);
      }
    } finally {
      setLoading(false);
    }
  };

  if (loading) {
    return <div className="loading">Loading logs...</div>;
  }

  if (error) {
    return (
      <div className="error">
        <p>Error: {error}</p>
        {error.includes('Telegram') && (
          <p className="error-hint">
            {error}
          </p>
        )}
      </div>
    );
  }

  const handleAddLog = () => {
    setFormMode('create');
    setEditingLog(null);
    setIsFormOpen(true);
  };

  const handleEditLog = (log: Log) => {
    setFormMode('edit');
    setEditingLog(log);
    setIsFormOpen(true);
  };

  const handleDeleteLog = (log: Log) => {
    setDeletingLog(log);
    setIsDeleteConfirmOpen(true);
  };

  const confirmDelete = async () => {
    if (!deletingLog) return;

    try {
      await deleteLog(deletingLog.id);
      await loadLogs();
      setIsDeleteConfirmOpen(false);
      setDeletingLog(null);
    } catch (err) {
      console.error('Failed to delete log:', err);
      // TODO: Show error message to user
    }
  };

  return (
    <div className="app">
      <header className="app-header">
        <h1>Calorie Log</h1>
        <button
          className="btn-add-log"
          onClick={handleAddLog}
        >
          Add New Log
        </button>
      </header>
      <main className="app-main">
        <LogTable
          logs={logs}
          onEdit={handleEditLog}
          onDelete={handleDeleteLog}
        />
      </main>

      <LogForm
        isOpen={isFormOpen}
        onClose={() => {
          setIsFormOpen(false);
          setEditingLog(null);
        }}
        onSuccess={loadLogs}
        initialData={editingLog}
        mode={formMode}
      />

      <DeleteConfirm
        isOpen={isDeleteConfirmOpen}
        onClose={() => {
          setIsDeleteConfirmOpen(false);
          setDeletingLog(null);
        }}
        onConfirm={confirmDelete}
        logInfo={deletingLog ? `${deletingLog.foodItems.join(', ')} - ${deletingLog.calories} cal` : undefined}
      />
    </div>
  );
}

export default App;
