export function isFutureWithinHours(
  date: Date,
  time: string,
  numOfHours: number,
): boolean {
  const [hours, minutes] = time.split(":").map(Number);

  const target = new Date(date);
  target.setHours(hours, minutes, 0, 0);

  const nowPlus = new Date(Date.now() + numOfHours * 60 * 60 * 1000);

  return nowPlus <= target;
}
