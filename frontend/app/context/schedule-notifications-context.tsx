import { createContext } from 'react';
import { ScheduledNotification } from '~/routes/schedules';

export const ScheduleNotificationsContext = createContext<{
  notifications: ScheduledNotification[];
}>({ notifications: [] });
