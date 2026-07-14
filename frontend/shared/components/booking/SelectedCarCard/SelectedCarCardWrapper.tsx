export function SelectedCarCardWrapper({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <div className="bg-white shadow-card p-6 flex rounded-2xl flex-col gap-2 justify-between border max-sm:border-b-0 border-cars-border relative">
      {children}
    </div>
  );
}
