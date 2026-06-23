import { create } from "zustand";
import { broker } from "@/shared/client";
import { SelectedFilters } from "@/app/(app)/[lang]/(withNavbar)/(booking)/results/_hooks/useCheckboxFilters";
import { FilterConfig } from "@/app/(app)/[lang]/(withNavbar)/(booking)/_components/_constants/filtersList";

interface BookingSessionState {
  selectedPlanIndex: number;
  isErpSelected: boolean;
  selectedAddons: broker.SelectAddOn[];
  checkboxFilters: SelectedFilters;
  carGroupFilters: Set<string>;
  setSelectedPlanIndex: (index: number) => void;
  setIsErpSelected: (value: boolean) => void;
  setSelectedAddons: (addons: broker.SelectAddOn[]) => void;
  setCheckboxFilters: (filterId: FilterConfig["id"], value: string) => void;
  clearAllCheckboxFilters: () => void;
  clearCarGroupFilters: () => void;
  toggleCarGroupFilter: (groupName: string) => void;
  clearSession: () => void;
}

export const useBookingSessionStore = create<BookingSessionState>((set) => ({
  selectedPlanIndex: 0,
  isErpSelected: false,
  selectedAddons: [],
  checkboxFilters: new Map(),
  carGroupFilters: new Set(),
  setSelectedPlanIndex: (index) => set({ selectedPlanIndex: index }),
  setIsErpSelected: (value) => set({ isErpSelected: value }),
  setSelectedAddons: (addons) => set({ selectedAddons: addons }),
  setCheckboxFilters: (filterId, value) =>
    set((prev) => {
      const next = new Map(prev.checkboxFilters);
      const nextSet = new Set(next.get(filterId) ?? []);

      if (nextSet.has(value)) {
        nextSet.delete(value);
      } else {
        nextSet.add(value);
      }

      if (nextSet.size === 0) {
        next.delete(filterId);
      } else {
        next.set(filterId, nextSet);
      }

      return { checkboxFilters: next };
    }),
  clearAllCheckboxFilters: () => set({ checkboxFilters: new Map() }),
  clearCarGroupFilters: () => set({ carGroupFilters: new Set() }),
  toggleCarGroupFilter: (groupName) =>
    set((prev) => {
      const next = new Set(prev.carGroupFilters);
      if (next.has(groupName)) {
        next.delete(groupName);
      } else {
        next.add(groupName);
      }
      return { carGroupFilters: next };
    }),
  clearSession: () =>
    set({
      selectedPlanIndex: 0,
      isErpSelected: false,
      selectedAddons: [],
      // checkboxFilters: new Map(),
      // carGroupFilters: new Set(),
    }),
}));
