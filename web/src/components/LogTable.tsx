import { Log } from '../api/logs';

interface LogTableProps {
  logs: Log[];
}

export function LogTable({ logs }: LogTableProps) {
  if (logs.length === 0) {
    return (
      <div className="empty-state">
        <p>No logs yet</p>
        <p className="empty-hint">Start tracking your calories by adding a new log entry.</p>
      </div>
    );
  }

  return (
    <div className="log-table-container">
      <table className="log-table">
        <thead>
          <tr>
            <th>Date/Time</th>
            <th>Food Items</th>
            <th>Calories</th>
            <th>Confidence</th>
            <th>Actions</th>
          </tr>
        </thead>
        <tbody>
          {logs.map((log) => (
            <tr key={log.id}>
              <td className="date-cell">
                {new Date(log.timestamp).toLocaleString()}
              </td>
              <td className="food-items-cell">
                {log.foodItems.join(', ')}
              </td>
              <td className="calories-cell">{log.calories}</td>
              <td className="confidence-cell">
                <span className={`confidence-badge confidence-${log.confidence}`}>
                  {log.confidence}
                </span>
              </td>
              <td className="actions-cell">
                <button className="btn-edit">Edit</button>
                <button className="btn-delete">Delete</button>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}
