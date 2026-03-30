import { authOptions } from "@/shared/auth/authOptions";
import { getServerSession } from "next-auth/next";
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
        name: "next-auth-encore",
        authenticate: async ({ payload }) => {
          const session = await getServerSession(authOptions);

          if (!session || !session.user || session.user.role !== "admin") {
            return { user: null };
          }

          const { docs } = await payload.find({
            collection: "admins",
            where: {
              userId: { equals: session.user.id },
            },
          });

          let payloadUser = docs[0];

          if (!payloadUser) {
            payloadUser = await payload.create({
              collection: "admins",
              data: {
                userId: session.user.id.toString(),
                username: session.user.username,
                email: "",
              },
              draft: false,
            });
          }

          // 5. Hand the user object back to Payload. They are now authenticated!
          return {
            user: {
              ...payloadUser,
            },
          };
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
