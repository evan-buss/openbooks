import { ScrollArea, Table, Text } from "@mantine/core";
import { useDebouncedValue } from "@mantine/hooks";
import {
  ColumnDef,
  FilterFn,
  flexRender,
  getCoreRowModel,
  getFacetedRowModel,
  getFacetedUniqueValues,
  getFilteredRowModel,
  getSortedRowModel,
  Row,
  Table as TableType,
  useReactTable
} from "@tanstack/react-table";

import { useVirtualizer } from "@tanstack/react-virtual";
import { ReactNode, useRef, useState } from "react";
import { BookDetail } from "../../../state/messages";
import Toolbar from "./Toolbar";
import { RankingInfo, rankItem } from "@tanstack/match-sorter-utils";
import classes from "./DataTable.module.css";
import { conditionalAttribute } from "../../../utils/attribute-helper";

declare module "@tanstack/table-core" {
  interface FilterFns {
    fuzzy: FilterFn<unknown>;
  }
  interface FilterMeta {
    itemRank: RankingInfo;
  }
}

const fuzzyFilter: FilterFn<any> = (row, columnId, value, addMeta) => {
  // Rank the item
  const itemRank = rankItem(row.getValue(columnId), value);

  // Store the itemRank info
  addMeta({
    itemRank
  });

  // Return if the item should be filtered in/out
  return itemRank.passed;
};

export type ColumnWidthFunc = (cols: number) => number;

interface DataTableProps<TData, TValue> {
  columns: ColumnDef<TData, TValue>[];
  data: TData[];
  facetFilters?: (table: TableType<TData>) => ReactNode[];
}

export default function DataTable<TData, TValue>({
  columns,
  data,
  facetFilters
}: DataTableProps<TData, TValue>) {
  // Table virtualization stuff
  const virtualizedRef = useRef(null);

  const [searchQuery, setSearchQuery] = useState<string>("");
  const [debouncedSearchQuery] = useDebouncedValue<string>(searchQuery, 250);

  const table = useReactTable({
    data: data,
    columns: columns,
    enableMultiSort: true,
    filterFns: {
      fuzzy: fuzzyFilter
    },
    state: {
      globalFilter: debouncedSearchQuery
    },
    defaultColumn: {
      minSize: 50,
      size: Number.MAX_SAFE_INTEGER,
      maxSize: Number.MAX_SAFE_INTEGER
    },
    onGlobalFilterChange: (value) => setSearchQuery(value ?? ""),
    globalFilterFn: "fuzzy",
    enableFilters: true,
    enableGlobalFilter: true,
    columnResizeMode: "onChange",
    getCoreRowModel: getCoreRowModel(),
    getFilteredRowModel: getFilteredRowModel(),
    getFacetedRowModel: getFacetedRowModel(),
    getFacetedUniqueValues: getFacetedUniqueValues(),
    getSortedRowModel: getSortedRowModel()
  });

  const { rows: tableRows } = table.getRowModel();

  const rowVirtualizer = useVirtualizer({
    count: tableRows.length,
    getScrollElement: () => virtualizedRef.current,
    estimateSize: () => 50,
    overscan: 10
  });

  const virtualItems = rowVirtualizer.getVirtualItems();

  const paddingTop =
    virtualItems.length > 0 ? virtualItems?.[0]?.start || 0 : 0;
  const paddingBottom =
    virtualItems.length > 0
      ? rowVirtualizer.getTotalSize() -
        (virtualItems?.[virtualItems.length - 1]?.end || 0)
      : 0;
  return (
    <>
      <Toolbar
        table={table}
        searchQuery={searchQuery}
        setSearchQuery={setSearchQuery}
        facetFilters={facetFilters}
      />
      <ScrollArea
        viewportRef={virtualizedRef}
        className={classes.container}
        type="hover"
        scrollbarSize={6}
        classNames={{
          thumb: classes.thumb
        }}
        offsetScrollbars={false}>
        <Table highlightOnHover verticalSpacing="0" fz="xs">
          <Table.Thead className={classes.head}>
            {table.getHeaderGroups().map((headerGroup) => (
              <Table.Tr key={headerGroup.id}>
                {headerGroup.headers.map((header) => (
                  <Table.Th
                    key={header.id}
                    pl="lg"
                    py="sm"
                    className={classes.headerCell}
                    style={{
                      width:
                        header.getSize() === Number.MAX_SAFE_INTEGER
                          ? "auto"
                          : header.getSize()
                    }}>
                    {flexRender(
                      header.column.columnDef.header,
                      header.getContext()
                    )}
                    <div
                      onMouseDown={header.getResizeHandler()}
                      onTouchStart={header.getResizeHandler()}
                      {...conditionalAttribute(
                        "resizing",
                        header.column.getIsResizing()
                      )}
                      className={classes.resizer}
                    />
                  </Table.Th>
                ))}
              </Table.Tr>
            ))}
          </Table.Thead>
          <Table.Tbody>
            {paddingTop > 0 && (
              <Table.Tr>
                <Table.Td style={{ height: `${paddingTop}px` }} />
              </Table.Tr>
            )}
            {rowVirtualizer.getVirtualItems().map((virtualRow) => {
              const row = tableRows[
                virtualRow.index
              ] as unknown as Row<BookDetail>;
              return (
                <Table.Tr key={row.id} style={{ height: 50 }}>
                  {row.getVisibleCells().map((cell) => {
                    return (
                      <Table.Td
                        key={cell.id}
                        style={{
                          width:
                            cell.column.getSize() === Number.MAX_SAFE_INTEGER
                              ? "auto"
                              : cell.column.getSize()
                        }}>
                        <Text
                          fz="xs"
                          component="div"
                          pl="xl"
                          lineClamp={1}
                          c="dark">
                          {flexRender(
                            cell.column.columnDef.cell,
                            cell.getContext()
                          )}
                        </Text>
                      </Table.Td>
                    );
                  })}
                </Table.Tr>
              );
            })}
            {paddingBottom > 0 && (
              <Table.Tr>
                <Table.Td style={{ height: `${paddingBottom}px` }} />
              </Table.Tr>
            )}
          </Table.Tbody>
        </Table>
      </ScrollArea>
    </>
  );
}
