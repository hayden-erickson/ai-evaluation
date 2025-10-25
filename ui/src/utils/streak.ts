import { Log } from '../types';

export const calculateStreak = (logs: Log[]): number => {
  if (logs.length === 0) {
    return 0;
  }

  const sortedLogs = logs
    .map((log) => new Date(log.date))
    .sort((a, b) => b.getTime() - a.getTime());

  let streak = 1;
  let lastDate = sortedLogs[0];

  for (let i = 1; i < sortedLogs.length; i++) {
    const currentDate = sortedLogs[i];
    const diffTime = lastDate.getTime() - currentDate.getTime();
    const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24));

    if (diffDays === 1) {
      streak++;
    } else if (diffDays > 2) {
      break;
    }
    lastDate = currentDate;
  }

  return streak;
};
