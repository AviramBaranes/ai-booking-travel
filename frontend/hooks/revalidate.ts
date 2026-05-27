import type {
  CollectionAfterChangeHook,
  GlobalAfterChangeHook,
} from "payload";

async function revalidateAllPages(origin: string) {
  await fetch(`${origin}/api/revalidate`, {
    method: "POST",
    headers: { "x-revalidate-secret": process.env.PAYLOAD_SECRET ?? "" },
  });
}

export const revalidateOnCollectionChange: CollectionAfterChangeHook =
  async ({ req }) => {
    await revalidateAllPages(req.origin);
  };

export const revalidateOnGlobalChange: GlobalAfterChangeHook = async ({
  req,
}) => {
  await revalidateAllPages(req.origin);
};
