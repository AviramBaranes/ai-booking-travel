import { getLang } from "@/shared/lang/lang";
import { NotFoundContent } from "@/shared/components/NotFoundContent";
import { getNotFoundData } from "../not-found";

export default async function NotFound() {
  const lang = await getLang();
  const notFoundData = await getNotFoundData(lang);

  return (
    <NotFoundContent
      title={notFoundData.title ?? ""}
      subtitle={notFoundData.subtitle ?? ""}
      buttonText={notFoundData.buttonText ?? ""}
      homepageUrl={`/${lang}`}
    />
  );
}
