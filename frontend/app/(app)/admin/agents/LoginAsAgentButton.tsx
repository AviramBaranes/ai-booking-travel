import LoginAsUserButton from "@/app/(app)/admin/_components/LoginAsUserButton";

export default function LoginAsAgentButton({ agentId }: { agentId: number }) {
  return <LoginAsUserButton userId={agentId} label="התחבר כסוכן" />;
}
