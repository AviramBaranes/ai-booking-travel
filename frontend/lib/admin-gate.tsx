import { cookies } from "next/headers";
import { redirect } from "next/navigation";
import { Button } from "@/components/ui/button";
import { Card } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";

async function hashPassword(password: string): Promise<string> {
  const encoder = new TextEncoder();
  const data = encoder.encode(password);
  const buffer = await crypto.subtle.digest("SHA-256", data);
  const hashArray = Array.from(new Uint8Array(buffer));
  return hashArray.map((b) => b.toString(16).padStart(2, "0")).join("");
}

async function verify(path: string, formData: FormData) {
  "use server";

  const submitted = formData.get("password");
  if (typeof submitted !== "string") {
    const c = await cookies();
    c.set("admin_gate_error", "1", { maxAge: 5, path: "/" });
    redirect(path);
  }

  const expected = process.env.ADMIN_GATE_PASSWORD;
  if (!expected) {
    console.error("ADMIN_GATE_PASSWORD not set");
    const c = await cookies();
    c.set("admin_gate_error", "1", { maxAge: 5, path: "/" });
    redirect(path);
  }

  const submittedHash = await hashPassword(submitted);
  const expectedHash = await hashPassword(expected);

  if (submittedHash === expectedHash) {
    const c = await cookies();
    c.set("admin_gate", expectedHash, {
      httpOnly: true,
      sameSite: "lax",
      secure: process.env.NODE_ENV === "production",
      path: "/",
      maxAge: 60 * 60 * 24,
    });
    c.delete("admin_gate_error");
    redirect(path);
  } else {
    const c = await cookies();
    c.set("admin_gate_error", "1", { maxAge: 5, path: "/" });
    redirect(path);
  }
}

export async function adminGate(path: string): Promise<React.ReactNode | null> {
  const c = await cookies();
  const expected = process.env.ADMIN_GATE_PASSWORD;

  if (!expected) {
    console.warn("ADMIN_GATE_PASSWORD not set, skipping gate");
    return null;
  }

  const expectedHash = await hashPassword(expected);
  const cookieValue = c.get("admin_gate")?.value;

  if (cookieValue === expectedHash) {
    return null;
  }

  const hasError = c.has("admin_gate_error");

  return (
    <div className="flex min-h-screen items-center justify-center bg-gray-50 px-4">
      <Card className="w-full max-w-sm">
        <div className="p-6">
          <h1 className="mb-6 text-xl font-bold text-gray-900">סיסמה דרושה</h1>

          <form action={verify.bind(null, path)} className="space-y-4">
            <div className="space-y-2">
              <Label htmlFor="password" className="text-sm font-medium">
                הזן סיסמה
              </Label>
              <Input
                id="password"
                name="password"
                type="password"
                autoComplete="off"
                autoFocus
                required
                dir="rtl"
              />
            </div>

            {hasError && (
              <p className="text-sm text-red-600">סיסמה שגויה, נסה שוב</p>
            )}

            <Button type="submit" className="w-full">
              כניסה
            </Button>
          </form>
        </div>
      </Card>
    </div>
  );
}
