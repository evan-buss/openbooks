import { createStyles, ScrollArea, Table, Text } from "@mantine/core";
import {
  useDebouncedValue,
  useElementSize,
  useMergedRef
} from "@mantine/hooks";
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
  SortingFn,
  sortingFns,
  useReactTable,
  Table as TableType
} from "@tanstack/react-table";

import { useVirtualizer } from "@tanstack/react-virtual";
import { ReactNode, useMemo, useRef, useState } from "react";
import { BookDetail } from "../../../state/messages";
import Toolbar from "./Toolbar";
import {
  compareItems,
  RankingInfo,
  rankItem
} from "@tanstack/match-sorter-utils";

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
  console.log("filtering");
  const itemRank = rankItem(row.getValue(columnId), value);

  // Store the itemRank info
  addMeta({
    itemRank
  });

  // Return if the item should be filtered in/out
  return itemRank.passed;
};

const fuzzySort: SortingFn<any> = (rowA, rowB, columnId) => {
  let dir = 0;

  // Only sort by rank if the column has ranking information
  if (rowA.columnFiltersMeta[columnId]) {
    dir = compareItems(
      rowA.columnFiltersMeta[columnId]?.itemRank!,
      rowB.columnFiltersMeta[columnId]?.itemRank!
    );
  }

  // Provide an alphanumeric fallback for when the item ranks are equal
  return dir === 0 ? sortingFns.alphanumeric(rowA, rowB, columnId) : dir;
};

export const useTableStyles = createStyles((theme) => ({
  container: {
    border: `1px solid ${
      theme.colorScheme === "dark" ? theme.colors.dark[3] : theme.colors.gray[3]
    }`,
    borderRadius: theme.radius.md,
    backgroundColor:
      theme.colorScheme === "dark" ? theme.colors.dark[6] : "white",
    height: "100%",
    overflow: "auto",
    width: "100%",
    boxShadow: theme.shadows.xs
  },
  head: {
    position: "sticky",
    top: 0,
    backgroundColor:
      theme.colorScheme === "dark"
        ? theme.colors.dark[5]
        : theme.colors.gray[1],
    zIndex: 1
  },
  headerCell: {
    position: "relative"
  },
  resizer: {
    position: "absolute",
    right: 0,
    top: 0,
    height: "100%",
    width: "2px",
    background:
      theme.colorScheme === "dark"
        ? theme.colors.dark[3]
        : theme.colors.gray[6],
    cursor: "col-resize",
    userSelect: "none",
    touchAction: "none",
    opacity: 0,

    ["&.isResizing"]: {
      background:
        theme.colorScheme === "dark"
          ? theme.colors.brand[3]
          : theme.colors.brand[4],
      opacity: 1
    },

    ["&:hover"]: {
      opacity: 1
    }
  }
}));

export type ColumnWidthFunc = (cols: number) => number;

interface DataTableProps<TData, TValue> {
  columns: (columnFunc: ColumnWidthFunc) => ColumnDef<TData, TValue>[];
  data: TData[];
  facetFilters?: (table: TableType<TData>) => ReactNode[];
}

export default function DataTable<TData, TValue>({
  columns,
  data,
  facetFilters
}: DataTableProps<TData, TValue>) {
  const { classes, cx, theme } = useTableStyles();

  // Table virtualization stuff
  const { ref: elementSizeRef, height, width } = useElementSize();
  const virtualizedRef = useRef();
  const mergedRef = useMergedRef(elementSizeRef, virtualizedRef);

  const columnDefs = useMemo(() => {
    // Use a function to calculate column widths based on the current width of the table
    const cols = (cols: number) => (width / 12) * cols;
    return columns(cols);
  }, [width]);

  const [searchQuery, setSearchQuery] = useState<string>("");
  const [debouncedSearchQuery, setDebouncedSearchQuery] =
    useDebouncedValue<string>(searchQuery, 250);

  const table = useReactTable({
    data: data,
    columns: columnDefs,
    enableMultiSort: true,
    filterFns: {
      fuzzy: fuzzyFilter
    },
    state: {
      globalFilter: debouncedSearchQuery
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
        viewportRef={mergedRef}
        className={classes.container}
        type="hover"
        scrollbarSize={6}
        styles={{ thumb: { ["&::before"]: { minWidth: 4 } } }}
        offsetScrollbars={false}>
        <Table highlightOnHover verticalSpacing="sm" fontSize="xs">
          <thead className={classes.head}>
            {table.getHeaderGroups().map((headerGroup) => (
              <tr key={headerGroup.id}>
                {headerGroup.headers.map((header) => (
                  <th
                    key={header.id}
                    className={classes.headerCell}
                    style={{
                      width: header.getSize()
                    }}>
                    {flexRender(
                      header.column.columnDef.header,
                      header.getContext()
                    )}
                    <div
                      onMouseDown={header.getResizeHandler()}
                      onTouchStart={header.getResizeHandler()}
                      className={cx(classes.resizer, {
                        ["isResizing"]: header.column.getIsResizing()
                      })}
                    />
                  </th>
                ))}
              </tr>
            ))}
          </thead>
          <tbody>
            {paddingTop > 0 && (
              <tr>
                <td style={{ height: `${paddingTop}px` }} />
              </tr>
            )}
            {rowVirtualizer.getVirtualItems().map((virtualRow) => {
              const row = tableRows[
                virtualRow.index
              ] as unknown as Row<BookDetail>;
              return (
                <tr key={row.id} style={{ height: 50 }}>
                  {row.getVisibleCells().map((cell) => {
                    return (
                      <td key={cell.id}>
                        <Text pl={20} lineClamp={1} color="dark">
                          {flexRender(
                            cell.column.columnDef.cell,
                            cell.getContext()
                          )}
                        </Text>
                      </td>
                    );
                  })}
                </tr>
              );
            })}
            {paddingBottom > 0 && (
              <tr>
                <td style={{ height: `${paddingBottom}px` }} />
              </tr>
            )}
          </tbody>
        </Table>
      </ScrollArea>
    </>
  );
}
