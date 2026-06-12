import sharp from "sharp";
import { fields, formBuilderPlugin } from "@payloadcms/plugin-form-builder";
import { gcsStorage } from "@payloadcms/storage-gcs";
import {
  FixedToolbarFeature,
  lexicalEditor,
} from "@payloadcms/richtext-lexical";
import { postgresAdapter } from "@payloadcms/db-postgres";
import { buildConfig } from "payload";
import { seoPlugin } from "@payloadcms/plugin-seo";
import { Admins } from "./CMS/collections/Admins";
import { Media } from "./CMS/collections/Media";
import { he } from "@payloadcms/translations/languages/he";
import { Pages } from "./CMS/collections/Pages";
import { SharedSections } from "./CMS/collections/SharedSections/SharedSections";
import { Header } from "./CMS/globals/Header";
import { Footer } from "./CMS/globals/Footer";
import { NotFoundConfig } from "./CMS/globals/NotFound";
import { Homepage } from "./CMS/globals/Homepage";
import { BookingSettings } from "./CMS/globals/BookingSettings";
import { SuppliersGallery } from "./CMS/globals/SupplierGallery";
import { AddonImagesGlobal } from "./CMS/globals/AddonImages";
import { formSettings } from "./CMS/settings/formSettings";
import { gcpStorageSettings } from "./CMS/settings/gcpStorageSettings";
import { emailAdapter } from "./CMS/email/emailAdapter";

export default buildConfig({
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
        const slug = doc.data.slug ?? "";
        const url = `${process.env.NEXT_PUBLIC_PAYLOAD_URL}/${doc.locale.code}/${slug}?payload_preview=1`;
        return url;
      },
      collections: ["pages"],
      globals: ["homepage"],
    },
  },

  routes: {
    admin: "/cms",
  },
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

  email: emailAdapter(),
  
  plugins: [
    formBuilderPlugin(formSettings),
    seoPlugin({
      collections: ["pages"],
      globals: ["homepage"],
      uploadsCollection: "media",
      tabbedUI: true,
      generateTitle: ({ doc }) => doc?.title ?? "",
      generateDescription: ({ doc }) => doc?.excerpt ?? "",
    }),
    gcsStorage(gcpStorageSettings),
  ],
  secret: process.env.PAYLOAD_SECRET || "",
  db: postgresAdapter({
    pool: {
      connectionString: process.env.DATABASE_URL,
    },
  }),
  sharp,
});
