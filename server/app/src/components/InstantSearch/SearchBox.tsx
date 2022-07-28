import { majorScale, SearchInput, Table } from "evergreen-ui";
import React, { useState } from "react";
import { useSearchBox } from "react-instantsearch-hooks-web";
import { useDebounce } from "react-use";

export function SearchBox() {
  const { refine } = useSearchBox();
  const [search, setSearch] = useState("");

  useDebounce(
    () => {
      refine(search);
    },
    250,
    [search]
  );

  return (
    <SearchInput
      value={search}
      placeholder="Search a Book"
      onChange={(e: any) => setSearch(e.target.value)}
      height={majorScale(5)}
      className="rounded-md"
      marginRight={majorScale(4)}
      width="100%"
    />
  );
}
