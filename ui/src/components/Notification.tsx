import React, { useState } from 'react';

interface NotificationProps {
  message: string;
  type: 'success' | 'error';
}

const Notification: React.FC<NotificationProps> = ({ message, type }) => {
  const [visible, setVisible] = useState(true);

  if (!visible) {
    return null;
  }

  return (
    <div className={`notification ${type}`}>
      {message}
      <button onClick={() => setVisible(false)}>X</button>
    </div>
  );
};

export default Notification;
