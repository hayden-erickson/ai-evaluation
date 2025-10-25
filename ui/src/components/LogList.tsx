import React from 'react';
import { Log } from '../../types';

interface LogListProps {
  logs: Log[];
  onLogClick: (log: Log) => void;
}

const LogList: React.FC<LogListProps> = ({ logs, onLogClick }) => {
  return (
    <div>
      <h4>Logs</h4>
      {logs.map((log) => (
        <div key={log.id} onClick={() => onLogClick(log)} style={{ cursor: 'pointer' }}>
          <p>{new Date(log.date).toLocaleDateString()}</p>
          <p>{log.notes}</p>
        </div>
      ))}
    </div>
  );
};

export default LogList;
