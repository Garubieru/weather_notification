import type { ActionFunctionArgs, MetaFunction } from '@remix-run/node';
import { Form, json, redirect, useActionData } from '@remix-run/react';
import { AlertCircle } from 'lucide-react';
import { Alert, AlertDescription, AlertTitle } from '~/components/ui/alert';
import { Button } from '~/components/ui/button';
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '~/components/ui/card';
import { Input } from '~/components/ui/input';

export const meta: MetaFunction = () => {
  return [{ title: 'Login' }, { name: 'Login', content: 'Login' }];
};

export default function Index() {
  const actionData = useActionData<typeof action>();

  return (
    <main className="font-sans p-4 flex justify-center items-center w-100dvh h-dvh">
      <Card className="p-2 max-w-[400px] w-[100%]">
        <CardHeader>
          <CardTitle>Login</CardTitle>
          <CardDescription>Digite seu nome de usuário e senha</CardDescription>
        </CardHeader>
        <CardContent>
          <Form method="post" className="flex flex-col gap-4">
            <Input placeholder="Nome de usuário" name="username"></Input>
            <Input placeholder="Senha" type="password" name="password"></Input>
            <Button type="submit">Entrar</Button>

            {actionData?.message ? (
              <Alert>
                <AlertCircle></AlertCircle>
                <AlertTitle>Erro ao realizar o login</AlertTitle>
                <AlertDescription>{actionData.message}</AlertDescription>
              </Alert>
            ) : (
              <></>
            )}
          </Form>
        </CardContent>
      </Card>
    </main>
  );
}

export async function action({ request }: ActionFunctionArgs) {
  const formData = await request.formData();

  const headers = new Headers();
  headers.set('content-type', 'application/json');

  const response = await fetch(`${process.env.BASE_URL}/v1/login`, {
    cache: 'no-cache',
    method: 'POST',
    headers,
    body: JSON.stringify({
      username: formData.get('username'),
      password: formData.get('password'),
    }),
    credentials: 'include',
  });

  switch (response.status) {
    case 200: {
      return redirect('/schedules/list', { headers: response.headers });
    }
    case 400: {
      const data = (await response.json()) as ApiError;

      return json({
        message:
          data.code in ErrorMessages
            ? ErrorMessages[data.code]
            : ErrorMessages.Default,
      });
    }
  }
}

const ErrorMessages = {
  InvalidUsernameOrPassword: 'Usuário ou senha inválidos.',
  Default: 'Erro no servidor',
} as const;

type ApiError<Code = keyof typeof ErrorMessages> = {
  code: Code;
  message: string;
};
