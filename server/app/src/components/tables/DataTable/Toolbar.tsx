import { Table } from "@tanstack/react-table";
import { Button, Group, TextInput } from "@mantine/core";
import { X } from "@phosphor-icons/react";
import React from "react";
import { ToolbarViewOptions } from "./ToolbarViewOptions";

interface ToolbarProps<TData> {
  table: Table<TData>;
  searchQuery: string;
  setSearchQuery: (searchQuery: string) => void;
  facetFilters?: (table: Table<TData>) => React.ReactNode[];
}

export default function Toolbar<TData>({
  table,
  searchQuery,
  setSearchQuery,
  facetFilters
}: ToolbarProps<TData>) {
  const isFiltered =
    table.getPreFilteredRowModel().rows.length >
    table.getFilteredRowModel().rows.length;

  const noData = table.getPreFilteredRowModel().rows.length === 0;

  return (
    <Group mb={10} w="100%" justify="space-between">
      <Group gap={8}>
        <TextInput
          w={{ xs: "150px", sm: "195px" }}
          variant={"filled"}
          size="xs"
          disabled={noData}
          value={searchQuery}
          onChange={(e) => setSearchQuery(e.currentTarget.value)}
          placeholder="Filter results..."
          radius="sm"
          type="search"
        />
        {facetFilters && facetFilters(table)}
        {isFiltered && (
          <Button
            variant="subtle"
            size="xs"
            onClick={() => {
              table.resetColumnFilters();
              table.resetGlobalFilter();
            }}
            rightSection={<X />}>
            Reset
          </Button>
        )}
      </Group>

      <ToolbarViewOptions table={table} />
    </Group>
  );
}
