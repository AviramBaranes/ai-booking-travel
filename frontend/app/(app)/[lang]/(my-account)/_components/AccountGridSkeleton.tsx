import { Bone } from "@/shared/components/Bone";

function AccountCardSkeleton() {
  return (
    <div className="p-6 flex flex-col gap-4 rounded-xl bg-white shadow-card">
      <div className="px-6 py-1 flex flex-col gap-2">
        <Bone className="h-3 w-20" />
        <Bone className="h-4 w-32" />
      </div>
      <div className="px-6 py-1 flex flex-col gap-2">
        <Bone className="h-3 w-24" />
        <Bone className="h-4 w-40" />
      </div>
      <div className="px-6 py-1 flex flex-col gap-2">
        <Bone className="h-3 w-20" />
        <Bone className="h-4 w-36" />
      </div>
      <div className="px-6 py-1 flex flex-col gap-2">
        <Bone className="h-3 w-24" />
        <Bone className="h-4 w-28 font-semibold" />
      </div>
      <div className="px-6 py-1 flex flex-col gap-2">
        <Bone className="h-3 w-12" />
        <Bone className="h-6 w-20 rounded-md" />
      </div>
      <div className="flex justify-between px-4 mt-2">
        <Bone className="h-8 w-24 rounded-md" />
        <Bone className="h-8 w-24 rounded-md" />
      </div>
    </div>
  );
}

export function AccountGridSkeleton() {
  return (
    <div className="grid grid-cols-4 gap-6">
      {Array.from({ length: 8 }).map((_, i) => (
        <AccountCardSkeleton key={i} />
      ))}
    </div>
  );
}