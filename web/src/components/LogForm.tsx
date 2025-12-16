import { useState, FormEvent } from 'react';
import { Dialog } from '@headlessui/react';
import { LogCreate, createLog } from '../api/logs';

interface LogFormProps {
  isOpen: boolean;
  onClose: () => void;
  onSuccess: () => void;
}

export function LogForm({ isOpen, onClose, onSuccess }: LogFormProps) {
  const [foodItemsText, setFoodItemsText] = useState('');
  const [calories, setCalories] = useState('');
  const [confidence, setConfidence] = useState<'high' | 'medium' | 'low'>('medium');
  const [errors, setErrors] = useState<string[]>([]);
  const [submitting, setSubmitting] = useState(false);

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault();
    setErrors([]);

    // Client-side validation
    const validationErrors: string[] = [];

    // Parse food items (comma-separated)
    const foodItems = foodItemsText
      .split(',')
      .map((item) => item.trim())
      .filter((item) => item.length > 0);

    if (foodItems.length === 0) {
      validationErrors.push('Food items cannot be empty');
    }

    if (foodItems.length > 10) {
      validationErrors.push('Cannot exceed 10 food items');
    }

    const totalLength = foodItems.join('').length;
    if (totalLength > 1000) {
      validationErrors.push('Total food items text cannot exceed 1000 characters');
    }

    const caloriesNum = parseInt(calories, 10);
    if (isNaN(caloriesNum) || caloriesNum < 0) {
      validationErrors.push('Calories must be a non-negative number');
    }

    if (validationErrors.length > 0) {
      setErrors(validationErrors);
      return;
    }

    // Submit log
    try {
      setSubmitting(true);
      const logData: LogCreate = {
        foodItems,
        calories: caloriesNum,
        confidence,
        timestamp: new Date().toISOString(),
      };

      await createLog(logData);

      // Reset form and close
      setFoodItemsText('');
      setCalories('');
      setConfidence('medium');
      onSuccess();
      onClose();
    } catch (err) {
      setErrors([err instanceof Error ? err.message : 'Failed to create log']);
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <Dialog open={isOpen} onClose={onClose} className="dialog-overlay">
      <div className="dialog-backdrop" aria-hidden="true" />

      <div className="dialog-container">
        <Dialog.Panel className="dialog-panel">
          <Dialog.Title className="dialog-title">Add New Log</Dialog.Title>

          <form onSubmit={handleSubmit} className="log-form">
            {errors.length > 0 && (
              <div className="form-errors">
                {errors.map((error, index) => (
                  <p key={index}>{error}</p>
                ))}
              </div>
            )}

            <div className="form-group">
              <label htmlFor="foodItems">Food Items (comma-separated)</label>
              <input
                type="text"
                id="foodItems"
                value={foodItemsText}
                onChange={(e) => setFoodItemsText(e.target.value)}
                placeholder="Pizza, Salad, Juice"
                required
              />
            </div>

            <div className="form-group">
              <label htmlFor="calories">Calories</label>
              <input
                type="number"
                id="calories"
                value={calories}
                onChange={(e) => setCalories(e.target.value)}
                placeholder="500"
                min="0"
                required
              />
            </div>

            <div className="form-group">
              <label htmlFor="confidence">Confidence Level</label>
              <select
                id="confidence"
                value={confidence}
                onChange={(e) => setConfidence(e.target.value as 'high' | 'medium' | 'low')}
              >
                <option value="high">High</option>
                <option value="medium">Medium</option>
                <option value="low">Low</option>
              </select>
            </div>

            <div className="form-actions">
              <button type="button" onClick={onClose} className="btn-cancel" disabled={submitting}>
                Cancel
              </button>
              <button type="submit" className="btn-submit" disabled={submitting}>
                {submitting ? 'Creating...' : 'Create Log'}
              </button>
            </div>
          </form>
        </Dialog.Panel>
      </div>
    </Dialog>
  );
}
