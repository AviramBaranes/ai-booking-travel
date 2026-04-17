export function SelectedCarCardWrapper({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <div className="bg-white shadow-card p-6 flex rounded-2xl flex-col gap-2 justify-between border border-cars-border">
      {children}
    </div>
  );
}
