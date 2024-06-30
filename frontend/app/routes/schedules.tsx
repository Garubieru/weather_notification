import type { LoaderFunctionArgs, MetaFunction } from '@remix-run/node';
import { NavLink, Outlet, useLoaderData, useLocation } from '@remix-run/react';
import { Card } from '~/components/ui/card';
import { UserSession, authenticate } from '~/utils/auth.server';
import { UserContext } from '~/context/schedule-notifications-context';

export const meta: MetaFunction = () => {
  return [
    { title: 'Notificações Agendadas' },
    { name: 'description', content: 'Notificações Agendadas' },
  ];
};

export async function loader({ request }: LoaderFunctionArgs): Promise<{
  user: UserSession;
  navLinks: Array<{ path: string; name: string }>;
}> {
  const userData = await authenticate(request);

  return {
    user: userData,
    navLinks: [
      { path: '/schedules/list', name: 'Agendamentos' },
      { path: '/schedules/notifications', name: 'Notificações' },
      { path: '/schedules/schedule', name: 'Agendar' },
    ],
  };
}

export default function Schedules() {
  const { user, navLinks } = useLoaderData<typeof loader>();

  const { pathname } = useLocation();

  return (
    <main className="font-sans p-4 flex flex-col justify-center items-center w-100dvh h-dvh">
      <div defaultValue="list" className="w-[500px] max-w-[100%]">
        <h1 className="text-2xl mb-2 text font-sm font-semibold">
          Olá, {user.name}
        </h1>

        <ul className="grid w-full grid-cols-3 list-none bg-neutral-700 rounded-md p-1">
          {navLinks.map((navLink) => (
            <li
              className="header-link"
              data-active={pathname === navLink.path}
              key={navLink.path}
            >
              <NavLink to={navLink.path}>{navLink.name}</NavLink>
            </li>
          ))}
        </ul>

        <Card className="max-w-[500px] w-[100%]">
          <UserContext.Provider value={{ user: user }}>
            <Outlet />
          </UserContext.Provider>
        </Card>
      </div>
    </main>
  );
}
