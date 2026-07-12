import type { CollectionConfig } from "payload";

export const Admins: CollectionConfig = {
  slug: "admins",
  admin: {
    useAsTitle: "username",
    hidden: true,
  },
  auth: {
    disableLocalStrategy: true,
    strategies: [
      {
        name: "encore-auth",
        authenticate: async ({ payload, headers }) => {
          // Get cookies from request headers
          const cookieHeader = headers.get("cookie");
          if (!cookieHeader) {
            return { user: null };
          }

          // CSRF guard: the refresh token is an ambient (auto-sent) cookie, so
          // reject any request whose Origin does not match the host serving the
          // CMS. Same-origin admin-UI requests always match; cross-site CSRF
          // requests carry a foreign Origin and are refused. Top-level GET
          // navigations have no Origin header and are allowed (read-only, and
          // their responses are unreadable cross-origin anyway).
          const origin = headers.get("origin");
          const host = headers.get("host");
          if (origin) {
            try {
              if (new URL(origin).host !== host) {
                return { user: null };
              }
            } catch {
              return { user: null };
            }
          }

          try {
            const baseURL = process.env.NEXT_PUBLIC_API_BASE_URL || "http://localhost:4000";
            const sessionResponse = await fetch(`${baseURL}/auth/session`, {
              method: "POST",
              headers: { cookie: cookieHeader },
            });

            if (!sessionResponse.ok) {
              return { user: null };
            }

            const userData = await sessionResponse.json();

            if (!userData || userData.role !== "admin") {
              return { user: null };
            }

            const { docs } = await payload.find({
              collection: "admins",
              where: {
                userId: { equals: userData.id },
              },
            });

            let payloadUser = docs[0];
            
            if (!payloadUser) {
              payloadUser = await payload.create({
                collection: "admins",
                data: {
                  userId: userData.id,
                  username: userData.firstName + " " + userData.lastName,
                },
                draft: false,
              });
            }

            return {
              user: {
                ...payloadUser,
              },
              // responseHeaders
            };
          } catch (error) {
            console.error("Payload auth error:", error);
            return { user: null };
          }
        },
      },
    ],
  },
  fields: [
    {
      name: "userId",
      type: "text",
      required: true,
      unique: true,
    },
    {
      name: "username",
      type: "text",
      required: true,
    },
  ],
};
