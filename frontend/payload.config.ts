import sharp from "sharp";
import { lexicalEditor } from "@payloadcms/richtext-lexical";
import { postgresAdapter } from "@payloadcms/db-postgres";
import { buildConfig } from "payload";
import { Admins } from "./collections/Admins";
import { Media } from "./collections/Media";
import { he } from "@payloadcms/translations/languages/he";
import { AddonImages } from "./collections/AddonImages";

export default buildConfig({
  // If you'd like to use Rich Text, pass your editor here
  editor: lexicalEditor(),
  localization: {
    defaultLocale: "he",
    locales: [
      {
        code: "he",
        label: "עברית",
        rtl: true,
      },
    ],
  },
  i18n: {
    fallbackLanguage: "he",
    supportedLanguages: { he },
    translations: {
      he: {
        general: {
          collections: "קולקציות",
        },
      },
    },
  },

  admin: {
    components: {
      header: ["@/shared/components/admin/AdminNavbar"],
      graphics: {
        Icon: "@/shared/components/admin/AdminHomeBtn",
      },
    },
    meta: {
      title: "BT Admin Panel",
      description: "AI Booking Travel Admin Panel",
      icons: [
        {
          rel: "icon",
          type: "image/png",
          url: "/favicon.ico",
        },
      ],
    },
  },

  routes: {
    admin: "/cms",
  },
  // Define and configure your collections in this array
  collections: [Admins, Media, AddonImages],

  // Your Payload secret - should be a complex and secure string, unguessable
  secret: process.env.PAYLOAD_SECRET || "",
  // Whichever Database Adapter you're using should go here
  // Mongoose is shown as an example, but you can also use Postgres
  db: postgresAdapter({
    pool: {
      connectionString: process.env.DATABASE_URL,
    },
  }),
  // If you want to resize images, crop, set focal point, etc.
  // make sure to install it and pass it to the config.
  // This is optional - if you don't need to do these things,
  // you don't need it!
  sharp,
});
