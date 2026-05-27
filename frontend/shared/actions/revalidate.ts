"use server";

import { revalidatePath } from "next/cache";
import { getServerSession } from "next-auth";
import { authOptions } from "@/shared/auth/authOptions";
import { getPayload } from "payload";
import config from "@payload-config";
import { SUPPORTED_LANGS, SupportedLang } from "@/shared/constants/supported_langs";

export async function revalidateHomepage(): Promise<{ ok: boolean; error?: string }> {
  const session = await getServerSession(authOptions);

  if (!session?.user || session.user.role !== "admin") {
    return { ok: false, error: "unauthorized" };
  }

  const payload = await getPayload({ config });

  for (const lang of SUPPORTED_LANGS) {
    revalidatePath(`/${lang}`);

    const result = await payload.find({
      collection: "pages",
      locale: lang as SupportedLang,
      draft: false,
      limit: 1000,
      select: { slug: true },
    });

    for (const page of result.docs) {
      if (page.slug) {
        revalidatePath(`/${lang}/${page.slug}`);
      }
    }
  }

  return { ok: true };
}
