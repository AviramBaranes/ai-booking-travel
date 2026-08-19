/**
 * Renders a JSON-LD block. Server component — render it inside the page body.
 *
 * `<` is escaped so a stray `</script>` inside CMS content cannot terminate the
 * tag early.
 */
export function JsonLd({ data }: { data: Record<string, unknown> }) {
  return (
    <script
      type="application/ld+json"
      dangerouslySetInnerHTML={{
        __html: JSON.stringify(data).replace(/</g, "\\u003c"),
      }}
    />
  );
}
