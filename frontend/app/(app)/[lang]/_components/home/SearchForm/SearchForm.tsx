"use client";

import { useState } from "react";
import { LocationCombobox } from "./LocationCombobox";
import { SearchFormOptions } from "./SearchFormOptions";

export function SearchForm() {
  const [isReturnDifferentLoc, setIsReturnDifferentLoc] = useState(false);

  return (
    <form className="flex flex-col w-10/12 mx-auto mt-4">
      <SearchFormOptions
        isReturnDifferentLoc={isReturnDifferentLoc}
        setIsReturnDifferentLoc={setIsReturnDifferentLoc}
      />
      <div className="bg-white w-full rounded-l-xl rounded-br-xl">
        <div className="flex w-1/2 my-5">
          <LocationCombobox placeholder="מהיכן תאספו את הרכב?" />
          {isReturnDifferentLoc && (
            <LocationCombobox placeholder="לאן תחזירו את הרכב?" />
          )}
        </div>
      </div>
    </form>
  );
}
