import type {
  CollectionAfterChangeHook,
  GlobalAfterChangeHook,
} from "payload";

export type RevalidateRequest = {
  collection?: string;
  slug?: string;
  global?: boolean;
};

async function revalidateAllPages(
  origin: string,
  payload?: RevalidateRequest,
) {
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
  collection,
}) => {
  await revalidateAllPages(req.origin, {
    // Must come from the firing collection — blog posts live under
    // /{lang}/blog/{slug}, not /{lang}/{slug}.
    collection: collection.slug,
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
