export function SuccessBadge({ children }: { children: React.ReactNode }) {
  return (
    <div className="bg-success/5 border border-success/30 text-success type-label p-6 rounded-xl">
      {children}
    </div>
  );
}
