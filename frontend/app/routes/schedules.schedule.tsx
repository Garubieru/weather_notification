import { ActionFunctionArgs } from '@remix-run/node';
import {
  Form,
  MetaFunction,
  useActionData,
  useLoaderData,
} from '@remix-run/react';
import { Button } from '~/components/ui/button';
import { CardContent } from '~/components/ui/card';
import { Input } from '~/components/ui/input';
import { Label } from '~/components/ui/label';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '~/components/ui/select';
import { Switch } from '~/components/ui/switch';
import { padNumber } from '~/lib/date';

export const meta: MetaFunction = () => {
  return [{ name: 'Agendar notificação', title: 'Agendar notificação' }];
};

export async function loader() {
  return {
    hours: Array.from({ length: 24 }).map((_, index) => ({
      value: index,
      label: `${padNumber(index)}:00`,
    })),
  };
}

export default function Schedule() {
  const { hours } = useLoaderData<typeof loader>();
  const actionData = useActionData<typeof action>();

  return (
    <CardContent className="pt-6">
      <Form method="POST" className="flex flex-col gap-3">
        <Input name="city" placeholder="Cidade" />
        <Select name="hour">
          <SelectTrigger className="w-[100%]">
            <SelectValue placeholder="Selecione uma hora" />
          </SelectTrigger>
          <SelectContent>
            {hours.map((hour) => (
              <SelectItem key={hour.value} value={String(hour.value)}>
                {hour.label}
              </SelectItem>
            ))}
          </SelectContent>
        </Select>
        <Input
          placeholder="Intervalo"
          name="interval"
          type="number"
          min={1}
          max={30}
        ></Input>
        <div className="flex content-center items-center gap-2">
          <Switch name="isCoastal" id="isCoastal" value={1} defaultValue={0} />
          <Label htmlFor="isCoastal">Litorânea</Label>
        </div>

        <Input className="hidden" defaultValue="WEB" name="method"></Input>

        <Button type="submit">Agendar</Button>
      </Form>
      {actionData?.error && <p>{actionData.error}</p>}
    </CardContent>
  );
}

export async function action({ request }: ActionFunctionArgs) {
  const formData = await request.formData();

  const headers = new Headers();

  headers.set('content-type', 'application/json');
  headers.set('cookie', request.headers.get('cookie') as string);

  const response = await fetch('http://localhost:3000/v1/account/schedules', {
    body: JSON.stringify({
      hour: Number(formData.get('hour')),
      intervalInDays: Number(formData.get('intervalInDays')),
      cityId: formData.get('city'),
      method: formData.get('method'),
      isCoastalCity: Boolean(Number(formData.get('isCoastal'))),
    }),
    method: 'POST',
    headers: headers,
    cache: 'no-cache',
  });

  if (response.status != 200) {
    return { error: 'Não foi possível agendar' };
  }

  return null;
}
