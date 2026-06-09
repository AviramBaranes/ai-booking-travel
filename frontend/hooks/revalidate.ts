import type {
  CollectionAfterChangeHook,
  GlobalAfterChangeHook,
} from "payload";

async function revalidateAllPages(origin: string, payload?: unknown) {
  await fetch(`${origin}/api/revalidate`, {
    method: "POST",
    headers: {
      "content-type": "application/json",
      "x-revalidate-secret": process.env.PAYLOAD_SECRET ?? "",
    },
    body: JSON.stringify(payload ?? {}),
  });
}

export const revalidateOnCollectionChange: CollectionAfterChangeHook = async ({
  req,
  doc,
}) => {
  await revalidateAllPages(req.origin, {
    collection: "pages",
    slug: doc.slug,
  });
};

export const revalidateOnGlobalChange: GlobalAfterChangeHook = async ({
  req,
}) => {
  await revalidateAllPages(req.origin, {
    global: true,
  });
};