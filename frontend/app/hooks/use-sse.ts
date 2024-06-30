import { useEffect, useState } from 'react';

export function useSSE<Payload extends Record<string, unknown>>(
  url: string
): Message<Payload>[] {
  const [messages, setMessages] = useState<Message<Payload>[]>([]);

  useEffect(() => {
    if (!url) return;

    const source = new EventSource(url, {
      withCredentials: true,
    });

    source.onmessage = (event) => {
      const message = JSON.parse(event.data) as Message<Payload>;
      console.log(message);
      setMessages((prevState) => [message, ...prevState]);
    };

    source.onopen = (event) => {
      console.log('connected', event);
    };

    source.onerror = (err) => {
      console.log('Coguld not connect to SSE', err);
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
