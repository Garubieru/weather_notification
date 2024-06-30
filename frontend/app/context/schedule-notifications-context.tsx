import { createContext } from 'react';
import { UserSession } from '~/utils/auth.server';

export const UserContext = createContext<{
  user?: UserSession;
}>({});
