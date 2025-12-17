import { useState, FormEvent, useEffect } from 'react';
import { Dialog } from '@headlessui/react';
import { Log, LogCreate, createLog, updateLog } from '../api/logs';

interface LogFormProps {
  isOpen: boolean;
  onClose: () => void;
  onSuccess: () => void;
  initialData?: Log | null;
  mode?: 'create' | 'edit';
}

export function LogForm({ isOpen, onClose, onSuccess, initialData = null, mode = 'create' }: LogFormProps) {
  const [foodItemsText, setFoodItemsText] = useState('');
  const [calories, setCalories] = useState('');
  const [confidence, setConfidence] = useState<'high' | 'medium' | 'low'>('medium');
  const [errors, setErrors] = useState<string[]>([]);
  const [submitting, setSubmitting] = useState(false);

  // Populate form with initialData when in edit mode
  useEffect(() => {
    if (isOpen && mode === 'edit' && initialData) {
      setFoodItemsText(initialData.foodItems.join(', '));
      setCalories(initialData.calories.toString());
      setConfidence(initialData.confidence);
    } else if (isOpen && mode === 'create') {
      // Reset form for create mode
      setFoodItemsText('');
      setCalories('');
      setConfidence('medium');
    }
  }, [isOpen, mode, initialData]);

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

      if (mode === 'create') {
        const logData: LogCreate = {
          foodItems,
          calories: caloriesNum,
          confidence,
          timestamp: new Date().toISOString(),
        };
        await createLog(logData);
      } else if (mode === 'edit' && initialData) {
        await updateLog(initialData.id, {
          foodItems,
          calories: caloriesNum,
          confidence,
        });
      }

      // Reset form and close
      setFoodItemsText('');
      setCalories('');
      setConfidence('medium');
      onSuccess();
      onClose();
    } catch (err) {
      const action = mode === 'create' ? 'create' : 'update';
      setErrors([err instanceof Error ? err.message : `Failed to ${action} log`]);
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <Dialog open={isOpen} onClose={onClose} className="dialog-overlay">
      <div className="dialog-backdrop" aria-hidden="true" />

      <div className="dialog-container">
        <Dialog.Panel className="dialog-panel">
          <Dialog.Title className="dialog-title">
            {mode === 'create' ? 'Add New Log' : 'Edit Log'}
          </Dialog.Title>

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
                {submitting
                  ? (mode === 'create' ? 'Creating...' : 'Saving...')
                  : (mode === 'create' ? 'Create Log' : 'Save Changes')
                }
              </button>
            </div>
          </form>
        </Dialog.Panel>
      </div>
    </Dialog>
  );
}
