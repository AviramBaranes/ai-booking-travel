import { SearchFormCheckbox } from "./SearchFormCheckbox";

interface SearchFormOptionsProps {
  label: string;
  isDropoffDifferentLoc: boolean;
  setIsDropoffDifferentLoc: (value: boolean) => void;
}

export function DifferentLocCheckbox({
  label,
  isDropoffDifferentLoc,
  setIsDropoffDifferentLoc,
}: SearchFormOptionsProps) {
  return (
    <SearchFormCheckbox
      label={label}
      value={isDropoffDifferentLoc}
      setValue={setIsDropoffDifferentLoc}
      id="dropoff-different-loc"
      name="dropoff-different-loc"
    />
  );
}
