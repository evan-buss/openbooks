import { CornerDialog } from "evergreen-ui";
import { Hit as AlgoliaHit } from "instantsearch.js/es/types";
import React from "react";
import { InstantSearch } from "react-instantsearch-hooks-web";
import { useLocalStorage } from "react-use";
import TypesenseAdapter from "typesense-instantsearch-adapter";
import { Grid as CustomGrid } from "../components/InstantSearch/Grid";
import { Pagination as CustomPagination } from "../components/InstantSearch/Pagination";
import { SearchBox } from "../components/InstantSearch/SearchBox";

const typesenseInstantsearchAdapter = new TypesenseAdapter({
  server: {
    apiKey: "5neo5QvlAxxtrxX0YEfoq5S9xUkHlmmF", // Be sure to use an API key that only allows search operations
    // m8yIAmgBGBLDo4MGgGiY5aAbYtbgKqDA
    nodes: [
      {
        host: "localhost",
        port: 7189,
        protocol: "https"
      }
    ],
    cacheSearchResultsForSeconds: 2 * 60 // Cache search results from server. Defaults to 2 minutes. Set to 0 to disable caching.
  },
  // The following parameters are directly passed to Typesense's search API endpoint.
  //  So you can pass any parameters supported by the search endpoint below.
  //  query_by is required.
  additionalSearchParameters: {
    query_by: "title,author",
    per_page: 14
  }
});
const searchClient = typesenseInstantsearchAdapter.searchClient;

export type BookHit = AlgoliaHit<{
  title: string;
  author: string;
  format: string;
  server: string;
  id: string;
  size: string;
  full: string;
}>;

export function LiveSearchPage() {
  return (
    <div className="w-full h-full m-8">
      <InstantSearch searchClient={searchClient} indexName="books">
        <div className="mb-4">
          <SearchBox />
        </div>
        <div className="mb-1">
          <CustomPagination />
        </div>
        <CustomGrid />
      </InstantSearch>
      <FirstTimePopup />
    </div>
  );
}

export function FirstTimePopup() {
  const [isShown, setIsShown] = useLocalStorage("show-live-search-popup", true);

  return (
    <>
      <CornerDialog
        title="Welcome to the Live Book Search"
        isShown={isShown}
        onCloseComplete={() => setIsShown(false)}>
        IRC Highway's eBook lists are parsed and indexed to provide a live
        search experience. Searches return results quickly and don't require
        sending a request to the IRC server.
      </CornerDialog>
    </>
  );
}
