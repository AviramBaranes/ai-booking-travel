import sharp from "sharp";
import {
  FixedToolbarFeature,
  lexicalEditor,
} from "@payloadcms/richtext-lexical";
import { postgresAdapter } from "@payloadcms/db-postgres";
import { buildConfig } from "payload";
import { seoPlugin } from "@payloadcms/plugin-seo";
import { Admins } from "./collections/Admins";
import { Media } from "./collections/Media";
import { he } from "@payloadcms/translations/languages/he";
import { AddonImages } from "./collections/AddonImages";
import { Pages } from "./collections/Pages";
import { SharedSections } from "./collections/SharedSections";
import { Homepage } from "./globals/Homepage";
import { Header } from "./globals/Header";
import { Footer } from "./globals/Footer";

export default buildConfig({
  // If you'd like to use Rich Text, pass your editor here
  editor: lexicalEditor({
    features: ({ defaultFeatures }) => [
      ...defaultFeatures,
      FixedToolbarFeature(),
    ],
  }),
  localization: {
    defaultLocale: "he",
    locales: [
      { code: "he", label: "עברית", rtl: true },
      { code: "en", label: "אנגלית", rtl: false },
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
  collections: [Admins, Media, AddonImages, Pages, SharedSections],

  globals: [Header, Footer],

  plugins: [
    seoPlugin({
      collections: ["pages"],
      // globals: ["homepage"],
      uploadsCollection: "media",
      tabbedUI: true,
      generateTitle: ({ doc }) => doc?.title ?? "",
      generateDescription: ({ doc }) => doc?.excerpt ?? "",
    }),
  ],

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
