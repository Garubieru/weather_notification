import { Badge } from '~/components/ui/badge';
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '~/components/ui/card';
import { Switch } from '~/components/ui/switch';
import { parseDate } from '~/lib/date';

export function ScheduledNotification(
  props: ScheduledNotificationProps
): JSX.Element {
  return (
    <Card className="relative">
      <CardHeader>
        <CardTitle>{props.city.name}</CardTitle>
        <CardDescription>
          {parseDate(new Date(props.scheduledDate))}
        </CardDescription>
      </CardHeader>
      <CardContent className="flex gap-3">
        <Badge className="text-sm">UF: {props.city.stateCode}</Badge>
        <Badge className="text-sm">
          Intervalo: {props.intervalInDays} dias
        </Badge>
        {props.city.isCoastal && <Badge className="text-sm">Litorânea</Badge>}
        <div className="absolute right-4 top-4">
          <Switch id={`${props.id}-toggle`} checked={props.active} />
        </div>
      </CardContent>
    </Card>
  );
}

type ScheduledNotificationProps = {
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
