import { QueryProvider } from "../_components/providers/QueryProvider";

export default function BookingLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return <QueryProvider>{children}</QueryProvider>;
}
