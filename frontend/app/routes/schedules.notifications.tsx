import { LoaderFunctionArgs } from '@remix-run/node';
import { MetaFunction, redirect, useLoaderData } from '@remix-run/react';
import { format } from 'date-fns';
import { Badge } from '~/components/ui/badge';
import { CardContent, Card, CardHeader, CardTitle } from '~/components/ui/card';
import { Message, useSSE } from '~/hooks/use-sse';
import { WeatherNotificationContent } from '~/protocols/weather-notification-content';

export const meta: MetaFunction = () => {
  return [
    { title: 'Notificações' },
    { name: 'description', content: 'Notificações' },
  ];
};

export async function loader({ request }: LoaderFunctionArgs) {
  const response = await fetch(
    `${process.env.BASE_URL}/v1/account/notifications`,
    { headers: request.headers }
  );

  if (response.status !== 200) {
    return redirect('/schedules/list');
  }

  const notifications = (await response.json()).notifications as Array<
    Message<WeatherNotificationContent>
  >;

  return {
    baseUrl: process.env.BASE_URL,
    notifications,
  };
}

export default function List() {
  const { baseUrl, notifications } = useLoaderData<typeof loader>();

  const newNotifications = useSSE<WeatherNotificationContent>(
    `${baseUrl}/v1/stream`,
    notifications
  );

  return (
    <CardContent className="flex flex-col gap-4 pt-6 max-h-[400px] overflow-auto">
      {notifications.length === 0 && newNotifications.length == 0 && (
        <p>Sem notificações</p>
      )}

      <p>Total de notificações: {newNotifications.length}</p>

      {newNotifications.map((notification) => (
        <Card className="p-1" key={notification.id}>
          <CardHeader className="flex flex-row items-center justify-between gap-3">
            <CardTitle className="">{notification.payload.cityName}</CardTitle>
            <Badge>{notification.payload.cityStateCode}</Badge>
          </CardHeader>

          {notification.payload.prediction.waveConditions != null && (
            <CardContent>
              <p>Agitação das ondas</p>
              <p>
                <b>Manhã:</b>
                {notification.payload.prediction.waveConditions.morning}
              </p>
              <p>
                <b>Tarde:</b>
                {notification.payload.prediction.waveConditions.afternoon}
              </p>
              <p>
                <b>Noite:</b>
                {notification.payload.prediction.waveConditions.evening}
              </p>
            </CardContent>
          )}

          <CardContent>
            {notification.payload.prediction.temperatures.map((prediction) => {
              const date = new Date(prediction.date);

              return (
                <div
                  className="flex border-t-2 pb-2 pt-2 last:border-b-2"
                  key={`${notification.id}-${prediction.date}`}
                >
                  <p>
                    <b>Dia: </b>
                    {format(date, 'dd/MM/yyyy')}
                  </p>
                  <p>
                    <b>Máxima:</b> {prediction.max}°C
                  </p>
                  <p>
                    <b>Mínima:</b> {prediction.min}°C
                  </p>
                  <p>
                    <b>Condição:</b> {prediction.condition}
                  </p>
                </div>
              );
            })}
          </CardContent>
        </Card>
      ))}
    </CardContent>
  );
}
