import { Badge } from "@/components/ui/badge";
import { location } from "@/shared/client";

export function BrokerCodeBadge({ code }: { code: location.MergeLocationBroker }) {
  return (
    <Badge
      variant={code.enabled ? "secondary" : "outline"}
      title={code.enabled ? undefined : "מושבת"}
      dir="ltr"
    >
      {code.broker} · {code.broker_location_id}
    </Badge>
  );
}

export default function LocationSummary({
  location: l,
}: {
  location: location.MergeLocation;
}) {
  return (
    <div className="text-sm text-gray-600 space-y-1">
      <div>
        #{l.id} · {l.name} · {l.city ?? "—"} · {l.country} ({l.country_code})
        {l.iata ? ` · ${l.iata}` : ""}
      </div>
      <div className="flex flex-wrap gap-1">
        {l.broker_codes.map((c) => (
          <BrokerCodeBadge key={c.id} code={c} />
        ))}
      </div>
    </div>
  );
}
