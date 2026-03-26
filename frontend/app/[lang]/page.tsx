import { tempLangTest } from "@/shared/api/locations";

export default async function Home() {
  const res = await tempLangTest();

  return (
    <>
      <h1>Welcome to AI Booking Travel lang is {res?.lang}</h1>
    </>
  );
}
