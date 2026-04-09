import { getServerSession } from "next-auth";
import { getTranslations } from "next-intl/server";
import { BackToAdminButton } from "./BackToAdminButton";

export async function BackToAdminBanner() {
  const session = await getServerSession();
  const t = await getTranslations("BackToAdmin");

  if (!session?.user?.isAdminAsAgent) return null;

  return (
    <div className="bg-brand text-white text-sm py-2 px-4 flex items-center justify-center gap-3">
      <span>{t("message")}</span>
      <BackToAdminButton
        accessToken={session.user.accessToken}
        buttonText={t("button")}
      />
    </div>
  );
}
