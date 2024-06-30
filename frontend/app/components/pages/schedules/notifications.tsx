import { Badge } from '~/components/ui/badge';
import { Card, CardContent, CardHeader, CardTitle } from '~/components/ui/card';
import { Message } from '~/hooks/use-sse';
import { format } from 'date-fns';

import { WeatherNotificationContent } from '~/routes/schedules';

export function Notifications(props: NotificationsProps): JSX.Element {
  return (
    <div className="flex flex-col gap-4">
      {props.notifications.map((notification) => (
        <Card className="p-1" key={notification.id}>
          <CardHeader className="flex flex-row items-center justify-between gap-3">
            <CardTitle className="">{notification.payload.cityName}</CardTitle>
            <Badge>{notification.payload.cityStateCode}</Badge>
          </CardHeader>

          <CardContent>
            {notification.payload.predictions.map((prediction) => {
              const date = new Date(prediction.date);

              return (
                <div
                  className="flex border-t-2 border-b-2 pb-2 pt-2"
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
                  {prediction.uvi !== 0 && (
                    <p>Radiação Solar: {prediction.uvi}</p>
                  )}
                </div>
              );
            })}
          </CardContent>
        </Card>
      ))}
    </div>
  );
}

type NotificationsProps = {
  notifications: Message<WeatherNotificationContent>[];
};
