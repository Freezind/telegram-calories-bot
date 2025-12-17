import { Dialog } from '@headlessui/react';

interface DeleteConfirmProps {
  isOpen: boolean;
  onClose: () => void;
  onConfirm: () => void;
  logInfo?: string;
}

export function DeleteConfirm({ isOpen, onClose, onConfirm, logInfo }: DeleteConfirmProps) {
  return (
    <Dialog open={isOpen} onClose={onClose} className="dialog-overlay">
      <div className="dialog-backdrop" aria-hidden="true" />

      <div className="dialog-container">
        <Dialog.Panel className="dialog-panel">
          <Dialog.Title className="dialog-title">Confirm Delete</Dialog.Title>

          <div className="confirm-content">
            <p>Are you sure you want to delete this log entry?</p>
            {logInfo && <p className="log-info">{logInfo}</p>}
            <p className="warning">This action cannot be undone.</p>
          </div>

          <div className="form-actions">
            <button type="button" onClick={onClose} className="btn-cancel">
              Cancel
            </button>
            <button type="button" onClick={onConfirm} className="btn-delete">
              Delete
            </button>
          </div>
        </Dialog.Panel>
      </div>
    </Dialog>
  );
}
