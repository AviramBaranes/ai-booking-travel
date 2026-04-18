import { Bone } from "@/shared/components/Bone";

function RowSkeleton() {
  return (
    <div className="flex justify-between">
      <Bone className="h-4 w-28" />
      <Bone className="h-4 w-20" />
    </div>
  );
}

function SubTitleSkeleton() {
  return (
    <div className="mt-8">
      <Bone className="h-5 w-36 mb-2" />
      <hr />
    </div>
  );
}

export function ReservationSummarySkeleton() {
  return (
    <div className="flex flex-col gap-2 shadow-card rounded-xl p-6 bg-white border border-cars-border">
      {/* Header */}
      <div className="flex items-center justify-between">
        <Bone className="h-6 w-32" />
        <Bone className="h-10 w-28 rounded-md" />
      </div>
      <hr />

      {/* Basic info rows */}
      <RowSkeleton />
      <RowSkeleton />
      <RowSkeleton />
      <RowSkeleton />

      {/* Cost breakdown */}
      <SubTitleSkeleton />
      <RowSkeleton />
      <RowSkeleton />
      <RowSkeleton />

      {/* Total bar */}
      <div className="bg-brand-blue py-3 px-5 flex justify-between items-center rounded-xl mt-8">
        <Bone variant="dark" className="h-4 w-24" />
        <Bone variant="dark" className="h-6 w-20" />
      </div>

      {/* Car details */}
      <SubTitleSkeleton />
      <RowSkeleton />
      <RowSkeleton />
      <RowSkeleton />

      {/* Inclusions */}
      <SubTitleSkeleton />
      <ul className="flex flex-col gap-2 mx-4">
        {Array.from({ length: 4 }).map((_, i) => (
          <li key={i} className="flex items-center gap-2">
            <Bone className="h-7 w-7 rounded-full shrink-0" />
            <Bone className="h-4 w-48" />
          </li>
        ))}
      </ul>
    </div>
  );
}
