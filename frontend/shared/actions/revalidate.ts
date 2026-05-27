"use server";

import { revalidatePath } from "next/cache";
import { getServerSession } from "next-auth";
import { authOptions } from "@/shared/auth/authOptions";

export async function revalidateHomepage(): Promise<{ ok: boolean; error?: string }> {
  const session = await getServerSession(authOptions);

  if (!session?.user || session.user.role !== "admin") {
    return { ok: false, error: "unauthorized" };
  }

  revalidatePath("/", "layout");

  return { ok: true };
}
