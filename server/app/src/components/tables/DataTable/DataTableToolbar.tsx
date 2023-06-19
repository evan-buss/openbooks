import { Table } from "@tanstack/react-table";
import { Button, Group, Space, TextInput } from "@mantine/core";
import { X } from "@phosphor-icons/react";
import React from "react";

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
    <Group mb={8} w="100%" position="apart">
      <Group>
        <TextInput
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
        <Space />
      </Group>
      {isFiltered && (
        <Button
          variant="default"
          size="xs"
          onClick={() => {
            table.resetColumnFilters();
            table.resetGlobalFilter();
          }}
          className="h-8 px-2 lg:px-3"
          leftIcon={<X />}>
          Reset
        </Button>
      )}
    </Group>
  );
}
