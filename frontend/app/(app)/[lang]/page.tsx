import { Metadata } from "next";

export const metadata: Metadata = {
  title: "Home",
  description: "AI Booking Travel - Find your next car to rent worldwide",
};

export default async function Home() {
  return (
    <>
      <h1>Welcome to AI Booking Travel</h1>
    </>
  );
}
