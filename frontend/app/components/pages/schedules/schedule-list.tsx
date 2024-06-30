import { useContext } from 'react';
import { ScheduledNotification } from '~/components/pages/schedules/scheduled-notification';
import { CardContent } from '~/components/ui/card';
import { ScheduleNotificationsContext } from '~/context/schedule-notifications-context';

export function ScheduleList(): JSX.Element {
  const { notifications } = useContext(ScheduleNotificationsContext);

  return (
    <CardContent className="flex flex-col gap-6">
      {notifications.map((scheduledNotification) => (
        <ScheduledNotification
          key={scheduledNotification.id}
          {...scheduledNotification}
        />
      ))}
    </CardContent>
  );
}

export function loader() {
  console.log('oi');
}
