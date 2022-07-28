import React from "react";
import { Pagination as EvergreenPagination } from "evergreen-ui";
import { usePagination, PaginationProps } from "react-instantsearch-hooks-web";
import { useConnector } from "react-instantsearch-hooks-web";
import connectStats from "instantsearch.js/es/connectors/stats/connectStats";

import type {
  StatsConnectorParams as UseStatsProps,
  StatsWidgetDescription
} from "instantsearch.js/es/connectors/stats/connectStats";

export function useStats(props?: UseStatsProps) {
  return useConnector<UseStatsProps, StatsWidgetDescription>(
    connectStats,
    props
  );
}

export function Pagination(props: PaginationProps) {
  const { currentRefinement, nbPages, refine } = usePagination(props);
  const { nbHits, processingTimeMS } = useStats();

  return (
    <div className="flex flex-row items-center justify-between">
      <EvergreenPagination
        onNextPage={() => refine(currentRefinement + 1)}
        onPreviousPage={() => refine(currentRefinement - 1)}
        onPageChange={(page) => refine(page - 1)}
        page={currentRefinement + 1}
        totalPages={nbPages}
      />
      <p className="mr-2 text-xs text-gray-600 ">
        {nbHits.toLocaleString("en-US")} results found in {processingTimeMS}ms
      </p>
    </div>
  );
}
