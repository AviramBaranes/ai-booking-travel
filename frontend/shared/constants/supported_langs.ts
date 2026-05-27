export const SUPPORTED_LANGS = ["he", "en"] as const;
export type SupportedLang = (typeof SUPPORTED_LANGS)[number];
