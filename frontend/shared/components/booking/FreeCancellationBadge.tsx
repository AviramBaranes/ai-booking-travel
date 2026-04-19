import Image from "next/image";

export const HOURS_BEFORE_PICKUP_TO_ALLOW_CANCELLATION = 48;

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

interface FreeCancellationBadgeProps {
  pickupDate: string;
  pickupTime: string;
  text: string;
}
export function FreeCancellationBadge({
  pickupDate,
  pickupTime,
  text,
}: FreeCancellationBadgeProps) {
  return (
    <>
      {isFutureWithinHours(
        new Date(pickupDate),
        pickupTime,
        HOURS_BEFORE_PICKUP_TO_ALLOW_CANCELLATION,
      ) && (
        <div className="flex gap-1 items-center ">
          <Image
            src="/assets/icons/V.svg"
            alt="Checked Icon"
            width={28}
            height={28}
            className="w-7 h-7"
          />
          <span className="type-label text-success">{text}</span>
        </div>
      )}
    </>
  );
}
