import { useEffect, useState } from 'react';
import { useToast } from '~/components/ui/use-toast';

export function useSSE<Payload extends Record<string, unknown>>(
  url: string,
  previousNotifications: Array<Message<Payload>>
): Message<Payload>[] {
  const [messages, setMessages] = useState<Message<Payload>[]>(
    previousNotifications
  );

  const { toast } = useToast();

  useEffect(() => {
    if (!url) return;

    const source = new EventSource(url, {
      withCredentials: true,
    });

    source.onmessage = (event) => {
      const message = JSON.parse(event.data) as Message<Payload>;
      toast({ description: 'Nova notificação de clima', title: 'Aviso' });
      setMessages((prevState) => [message, ...prevState]);
    };

    source.onopen = (event) => {
      console.log('Connected', event);
    };

    source.onerror = (err) => {
      console.log('Could not connect to SSE', err);
      source.close();
    };

    window.onbeforeunload = () => {
      source.close();
    };

    return () => {
      source.close();
    };
  }, [url]);

  return messages;
}

export type Message<
  Payload extends Record<string, unknown> = Record<string, unknown>
> = {
  id: string;
  name: string;
  payload: Payload;
};
