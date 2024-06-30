import { LoaderFunctionArgs } from '@remix-run/node';
import { MetaFunction, useLoaderData } from '@remix-run/react';
import { ScheduledNotification } from '~/components/pages/schedules/scheduled-notification';
import { CardContent } from '~/components/ui/card';

export const meta: MetaFunction = () => {
  return [
    { title: 'Notificações Agendadas' },
    { name: 'description', content: 'Notificações Agendadas' },
  ];
};

export async function loader({ request }: LoaderFunctionArgs) {
  const response = await fetch('http://localhost:3000/v1/account/schedules', {
    method: 'GET',
    headers: request.headers,
  });

  const notifications = (await response.json())
    .scheduledNotifications as ScheduledNotification[];

  if (response.status !== 200) {
    return { scheduledNotifications: [] };
  }

  return {
    scheduledNotifications: notifications,
  };
}

export default function List() {
  const { scheduledNotifications } = useLoaderData<typeof loader>();

  return (
    <CardContent className="flex flex-col gap-2 pt-6">
      {scheduledNotifications.length === 0 && <p>Não há agendamentos</p>}

      {scheduledNotifications.map((scheduledNotification) => (
        <ScheduledNotification
          key={scheduledNotification.id}
          {...scheduledNotification}
        />
      ))}
    </CardContent>
  );
}

export type ScheduledNotification = {
  id: string;
  scheduledDate: string;
  intervalInDays: number;
  hour: number;
  active: boolean;
  city: {
    id: string;
    name: string;
    stateCode: string;
    isCoastal: boolean;
  };
};
