"use client";

import { SearchFormOptions } from "./SearchFormOptions";

export function SearchForm() {
  return (
    <form className="flex flex-col w-10/12 mx-auto mt-4">
      <SearchFormOptions />
      <div className="bg-white w-full h-26.5 rounded-l-xl rounded-br-xl"></div>
    </form>
  );
}
