import { Bone } from "@/shared/components/Bone";

export function SelectedCarCardSkeleton() {
  return (
    <div className="bg-white shadow-card p-6 flex rounded-2xl flex-col gap-2 justify-between border border-cars-border">
      {/* Supplier logo + car image */}
      <div className="flex-col flex items-center">
        <div className="mb-12">
          <Bone className="h-8 w-24" />
        </div>
        <Bone className="w-44 h-25" />
      </div>

      {/* Model name + "or similar" */}
      <div className="flex gap-2 flex-col items-start">
        <Bone className="h-5 w-40" />
        <Bone className="h-4 w-24" />
      </div>

      {/* Car details pills */}
      <div className="mb-6 flex gap-2 flex-wrap">
        <Bone className="h-7 w-16 rounded-full" />
        <Bone className="h-7 w-20 rounded-full" />
        <Bone className="h-7 w-14 rounded-full" />
        <Bone className="h-7 w-18 rounded-full" />
      </div>

      {/* Free cancellation badge */}
      <div className="flex gap-1 items-center">
        <Bone className="h-7 w-7 rounded-full" />
        <Bone className="h-4 w-32" />
      </div>
    </div>
  );
}
