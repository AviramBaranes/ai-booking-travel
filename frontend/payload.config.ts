import sharp from "sharp";
import { gcsStorage } from "@payloadcms/storage-gcs";
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
import { Pages } from "./collections/Pages";
import { SharedSections } from "./collections/SharedSections";
import { AddonImagesGlobal } from "./globals/AddonImages";
import { BookingSettings } from "./globals/BookingSettings";
import { Homepage } from "./globals/Homepage";
import { Header } from "./globals/Header";
import { Footer } from "./globals/Footer";
import { NotFoundConfig } from "./globals/NotFound";
import { SuppliersGallery } from "./globals/SupplierGallery";

const storageEnv = process.env.PAYLOAD_STORAGE_ENV ?? "local";

if (!["local", "dev", "prod"].includes(storageEnv)) {
  throw new Error("PAYLOAD_STORAGE_ENV must be one of: local, dev, prod");
}

const isGcsStorageEnabled = storageEnv === "dev" || storageEnv === "prod";
const gcsBucket =
  storageEnv === "dev"
    ? process.env.PAYLOAD_GCS_DEV_BUCKET
    : storageEnv === "prod"
      ? process.env.PAYLOAD_GCS_PROD_BUCKET
      : undefined;
const gcsPrivateKey = process.env.PAYLOAD_GCS_PRIVATE_KEY?.replace(
  /\\n/g,
  "\n",
);

if (isGcsStorageEnabled && !gcsBucket) {
  throw new Error(
    `Missing ${storageEnv === "dev" ? "PAYLOAD_GCS_DEV_BUCKET" : "PAYLOAD_GCS_PROD_BUCKET"}`,
  );
}

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
    theme: "light",
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
    livePreview: {
      url: (doc) => {
        const url = `${process.env.NEXT_PUBLIC_PAYLOAD_URL}/${doc.locale.code}/${doc.data.slug}?payload_preview=1`;
        return url;
      },
      collections: ["pages"],
    },
  },

  routes: {
    admin: "/cms",
  },
  // Define and configure your collections in this array
  collections: [Admins, Media, Pages, SharedSections],

  globals: [
    Header,
    Footer,
    NotFoundConfig,
    Homepage,
    BookingSettings,
    SuppliersGallery,
    AddonImagesGlobal,
  ],

  plugins: [
    seoPlugin({
      collections: ["pages"],
      globals: ["homepage"],
      uploadsCollection: "media",
      tabbedUI: true,
      generateTitle: ({ doc }) => doc?.title ?? "",
      generateDescription: ({ doc }) => doc?.excerpt ?? "",
    }),
    gcsStorage({
      enabled: isGcsStorageEnabled,
      collections: {
        media: {
          prefix: `${storageEnv}/media`,
        },
      },
      bucket: gcsBucket ?? "local-disabled",
      options: {
        projectId: process.env.PAYLOAD_GCS_PROJECT_ID,
        credentials:
          process.env.PAYLOAD_GCS_CLIENT_EMAIL && gcsPrivateKey
            ? {
                client_email: process.env.PAYLOAD_GCS_CLIENT_EMAIL,
                private_key: gcsPrivateKey,
              }
            : undefined,
      },
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

