import { redirect } from "next/navigation";

type Props = {
  params: Promise<{
    lang: string;
  }>;
};

export default async function BlogPage({ params }: Props) {
  const { lang } = await params;

  redirect(`/${lang}/blog/page/1`);
}