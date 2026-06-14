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

export const gcpStorageSettings = {
  enabled: isGcsStorageEnabled,
  clientUploads: true,
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
};
