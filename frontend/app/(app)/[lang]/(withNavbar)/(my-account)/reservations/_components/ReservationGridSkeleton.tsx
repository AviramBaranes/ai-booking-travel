import { Bone } from "@/shared/components/Bone";

function ReservationCardSkeleton() {
  return (
    <div className="p-6 flex flex-col gap-4 rounded-xl bg-white shadow-card">
      {/* Pickup date */}
      <div className="px-6 py-1 flex flex-col gap-2">
        <Bone className="h-3 w-20" />
        <Bone className="h-4 w-32" />
      </div>
      {/* Pickup location */}
      <div className="px-6 py-1 flex flex-col gap-2">
        <Bone className="h-3 w-24" />
        <Bone className="h-4 w-40" />
      </div>
      {/* Driver name */}
      <div className="px-6 py-1 flex flex-col gap-2">
        <Bone className="h-3 w-20" />
        <Bone className="h-4 w-36" />
      </div>
      {/* Booking number */}
      <div className="px-6 py-1 flex flex-col gap-2">
        <Bone className="h-3 w-24" />
        <Bone className="h-4 w-28 font-semibold" />
      </div>
      {/* Status badge */}
      <div className="px-6 py-1 flex flex-col gap-2">
        <Bone className="h-3 w-12" />
        <Bone className="h-6 w-20 rounded-md" />
      </div>
      {/* Action buttons */}
      <div className="flex justify-between px-4 mt-2">
        <Bone className="h-8 w-24 rounded-md" />
        <Bone className="h-8 w-24 rounded-md" />
      </div>
    </div>
  );
}

export function ReservationGridSkeleton() {
  return (
    <div className="grid grid-cols-4 gap-6">
      {Array.from({ length: 8 }).map((_, i) => (
        <ReservationCardSkeleton key={i} />
      ))}
    </div>
  );
}
