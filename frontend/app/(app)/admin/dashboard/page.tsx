import { adminGate } from "@/lib/admin-gate";
import DashboardShell from "./DashboardShell";

export default async function DashboardPage() {
  const gate = await adminGate("/admin/dashboard");
  if (gate) return gate;
  return <DashboardShell />;
}
