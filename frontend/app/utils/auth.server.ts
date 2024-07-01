import { redirect } from '@remix-run/react';

export async function authenticate(request: Request): Promise<UserSession> {
  const response = await fetch(`${process.env.BASE_URL}/v1/session`, {
    credentials: 'include',
    headers: request.headers,
    method: 'GET',
  });

  if (response.status !== 200) {
    throw redirect('/login');
  }

  return response.json();
}

export type UserSession = {
  email: string;
  id: string;
  name: string;
  phone: string;
  username: string;
};
