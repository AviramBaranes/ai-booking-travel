import { NotFoundContent } from "@/shared/components/NotFoundContent";
import { getNotFoundData } from "../../not-found";

export default async function NotFound() {
  const lang = "he"// if using headers pages will become dynamic.
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
