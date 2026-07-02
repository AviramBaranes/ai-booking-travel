import Image from "next/image";
import { CAR_GROUPS_FILTERS } from "../../../_components/_constants/carGroupsFilters";
import clsx from "clsx";
import { useBookingSessionStore } from "@/shared/store/bookingSessionStore";

interface CarGroupFiltersProps {
  title: string;
}

export function CarGroupsFilter({ title }: CarGroupFiltersProps) {
  const selectedGroups = useBookingSessionStore(
    (state) => state.carGroupFilters,
  );
  const toggleGroup = useBookingSessionStore(
    (state) => state.toggleCarGroupFilter,
  );

  return (
    <div className="lg:mt-12 overflow-x-scroll">
      <h5 className="type-h5 mb-8 text-navy">{title}</h5>
      <div className="max-lg:overflow-x-auto">
        <div className="flex items-center justify-between max-lg:w-max max-lg:flex-nowrap max-lg:gap-2">
          {CAR_GROUPS_FILTERS.map((group) => (
            <div
              onClick={() => toggleGroup(group.name)}
              className={clsx(
                "bg-white text-center rounded-lg shadow-card px-2 py-2 max-lg:my-2 cursor-pointer hover:shadow-card-hover",
                {
                  "border-brand border": selectedGroups.has(group.name),
                },
              )}
              key={group.name}
            >
              <p
                className={clsx("type-paragraph text-navy", {
                  "font-bold": selectedGroups.has(group.name),
                })}
              >
                {group.name}
              </p>
              <Image
                src={group.image}
                alt={group.name}
                width={124}
                height={90}
                className="w-31 h-22.5"
              />
            </div>
          ))}
        </div>
      </div>
    </div>
  );
}
