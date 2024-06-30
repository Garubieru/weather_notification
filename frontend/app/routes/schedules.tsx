import type { LoaderFunctionArgs, MetaFunction } from '@remix-run/node';
import { useLoaderData } from '@remix-run/react';
import {
  Card,
  CardDescription,
  CardHeader,
  CardTitle,
} from '~/components/ui/card';
import { ScheduledNotification } from '~/components/pages/schedules/scheduled-notification';
import { UserSession, authenticate } from '~/utils/auth.server';
import { ScheduleNotificationsContext } from '~/context/schedule-notifications-context';
import { ScheduleList } from '~/components/pages/schedules/schedule-list';
import { Tabs, TabsList, TabsTrigger } from '~/components/ui/tabs';
import { TabsContent } from '@radix-ui/react-tabs';
import { useSSE } from '~/hooks/use-sse';
import { Notifications } from '~/components/pages/schedules/notifications';

export const meta: MetaFunction = () => {
  return [
    { title: 'Notificações Agendadas' },
    { name: 'description', content: 'Notificações Agendadas' },
  ];
};

export async function loader({ request }: LoaderFunctionArgs): Promise<{
  user: UserSession;
  scheduledNotifications: ScheduledNotification[];
  sseUrl: string;
}> {
  const userData = await authenticate(request);

  const response = await fetch('http://localhost:3000/v1/account/schedules', {
    method: 'GET',
    headers: request.headers,
  });

  if (response.status !== 200) {
    return { user: userData, scheduledNotifications: [], sseUrl: '' };
  }

  return {
    user: userData,
    scheduledNotifications: (await response.json()).scheduledNotifications,
    sseUrl: 'http://localhost:3000/v1/stream',
  };
}

export default function Schedules() {
  const { user, scheduledNotifications, sseUrl } =
    useLoaderData<typeof loader>();

  const notifications = useSSE<WeatherNotificationContent>(sseUrl);

  return (
    <main className="font-sans p-4 flex justify-center items-center w-100dvh h-dvh">
      <Tabs defaultValue="list" className="w-[500px]">
        <TabsList className="grid w-full grid-cols-3">
          <TabsTrigger value="list">Agendamentos</TabsTrigger>
          <TabsTrigger value="notifications">Notificações</TabsTrigger>
          <TabsTrigger value="schedule">Agendar</TabsTrigger>
        </TabsList>

        <Card className="p-2 max-w-[500px] w-[100%]">
          <CardHeader>
            <CardTitle>Olá, {user.name}</CardTitle>
            <CardDescription>
              Gerencia as notificações de clima agendadas
            </CardDescription>
          </CardHeader>
          <ScheduleNotificationsContext.Provider
            value={{ notifications: scheduledNotifications }}
          >
            <TabsContent value="list">
              <ScheduleList />
            </TabsContent>

            <TabsContent value="notifications">
              <Notifications notifications={notifications} />
            </TabsContent>

            <TabsContent value="schedule">
              <h1>WIP...</h1>
            </TabsContent>
          </ScheduleNotificationsContext.Provider>
        </Card>
      </Tabs>
    </main>
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

export type WeatherNotificationContent = {
  cityName: string;
  cityStateCode: string;
  predictions: Array<{
    date: string;
    min: number;
    max: number;
    condition: string;
    uvi: number;
  }>;
};
