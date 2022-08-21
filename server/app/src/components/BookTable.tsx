import { Box, Button, createStyles, Table } from "@mantine/core";
import {
  createColumnHelper,
  flexRender,
  getCoreRowModel,
  Row,
  useReactTable
} from "@tanstack/react-table";
import { useVirtualizer } from "@tanstack/react-virtual";
import { useRef } from "react";
import { BookDetail } from "../state/messages";
import "./BookTable.css";

export interface Props {
  books: BookDetail[];
}
const columnHelper = createColumnHelper<BookDetail>();

const columns = [
  columnHelper.accessor("server", {
    header: "Server",
    minSize: 50,
    size: 50
  }),
  columnHelper.accessor("author", {
    header: "Author"
  }),
  columnHelper.accessor("title", {
    header: "Title"
  }),
  columnHelper.accessor("format", {
    header: "Format"
  }),
  columnHelper.accessor("size", {
    header: "Size"
  }),
  columnHelper.display({
    header: "Download",
    cell: (props) => <Button>Download</Button>
  })
];

const useStyles = createStyles((theme) => ({
  container: {
    border: `1px solid ${
      theme.colorScheme === "dark" ? theme.colors.dark[3] : theme.colors.gray[3]
    }`,
    borderRadius: theme.radius.md,
    backgroundColor:
      theme.colorScheme === "dark" ? theme.colors.dark[6] : "white",
    height: "100%",
    overflow: "auto",
    width: "100%"
  },
  table: {},
  head: {
    position: "sticky",
    top: 0,
    backgroundColor:
      theme.colorScheme === "dark" ? theme.colors.dark[6] : "white",
    zIndex: 1,
    // Border bottom doesn't work since the table is border-collapse: collapse.
    boxShadow: `0px 1px 0px 0px ${
      theme.colorScheme === "dark" ? theme.colors.dark[3] : theme.colors.gray[3]
    }`
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
  },
  th: {
    position: "relative"
  },
  row: {}
}));

export default function BookTable({ books: data }: Props) {
  const { classes, cx } = useStyles();

  const tableContainerRef = useRef<HTMLDivElement>(null);
  const rowVirtualizer = useVirtualizer({
    count: data.length,
    getScrollElement: () => tableContainerRef.current,
    estimateSize: () => 57,
    overscan: 10
  });

  const paddingTop =
    rowVirtualizer.getVirtualItems().length > 0
      ? rowVirtualizer.getVirtualItems()?.[0]?.start || 0
      : 0;
  const paddingBottom =
    rowVirtualizer.getVirtualItems().length > 0
      ? rowVirtualizer.getTotalSize() -
        (rowVirtualizer.getVirtualItems()?.[
          rowVirtualizer.getVirtualItems().length - 1
        ]?.end || 0)
      : 0;

  const table = useReactTable({
    data,
    columns: columns,
    columnResizeMode: "onChange",
    getCoreRowModel: getCoreRowModel()
  });

  const { rows: tableRows } = table.getRowModel();

  return (
    <Box ref={tableContainerRef} className={classes.container}>
      <Table highlightOnHover verticalSpacing="xs">
        <thead className={classes.head}>
          {table.getHeaderGroups().map((headerGroup) => (
            <tr key={headerGroup.id}>
              {headerGroup.headers.map((header) => (
                <th
                  className={classes.th}
                  key={header.id}
                  style={{
                    width: header.getSize()
                  }}>
                  {header.isPlaceholder
                    ? null
                    : flexRender(
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
        <tbody style={{ height: "100px", maxHeight: "100px" }}>
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
              <tr key={row.id} className={classes.row}>
                {row.getVisibleCells().map((cell) => {
                  return (
                    <td key={cell.id}>
                      {flexRender(
                        cell.column.columnDef.cell,
                        cell.getContext()
                      )}
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
    </Box>
  );
}
