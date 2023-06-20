import { Table } from "@tanstack/react-table";
import { Button, Group, Space, TextInput } from "@mantine/core";
import { X } from "@phosphor-icons/react";
import React from "react";
import { DataTableViewOptions } from "./DataTableViewOptions";

interface DataTableToolbarProps<TData> {
  table: Table<TData>;
  searchQuery: string;
  setSearchQuery: (searchQuery: string) => void;
  facetFilters?: (table: Table<TData>) => React.ReactNode[];
}

export default function DataTableToolbar<TData>({
  table,
  searchQuery,
  setSearchQuery,
  facetFilters
}: DataTableToolbarProps<TData>) {
  const isFiltered =
    table.getPreFilteredRowModel().rows.length >
    table.getFilteredRowModel().rows.length;

  return (
    <Group mb={10} w="100%" position="apart">
      <Group spacing={8}>
        <TextInput
          w={{ xs: "150px", sm: "195px" }}
          variant={"filled"}
          size="xs"
          value={searchQuery}
          onChange={(e: any) => setSearchQuery(e.currentTarget.value)}
          placeholder="Filter results..."
          radius="sm"
          type="search"
        />
        {facetFilters &&
          facetFilters(table)?.map((facetFilter, index) => facetFilter)}
        {isFiltered && (
          <Button
            variant="subtle"
            size="xs"
            onClick={() => {
              table.resetColumnFilters();
              table.resetGlobalFilter();
            }}
            rightIcon={<X />}>
            Reset
          </Button>
        )}
      </Group>

      <DataTableViewOptions table={table} />
    </Group>
  );
}
