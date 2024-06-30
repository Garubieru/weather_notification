

export function parseDate(date: Date): string {
  const day = padNumber(date.getDate());
  const month = padNumber(date.getMonth());
  const year = padNumber(date.getFullYear());

  const hour = padNumber(date.getHours());
  const minutes = padNumber(date.getMinutes());

  return `${day}/${month}/${year} às ${hour}:${minutes}`;
}

export function padNumber(num: number): string {
  return String(num).padStart(2, '0');
}
