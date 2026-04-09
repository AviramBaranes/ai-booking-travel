import { useParams } from "next/navigation";

export function useDirection() {
  const { lang } = useParams();
  return lang === "he" ? "rtl" : "ltr";
}
