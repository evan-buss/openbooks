import { createStyles, Table } from "@mantine/core";
import {
  createColumnHelper,
  flexRender,
  getCoreRowModel,
  useReactTable
} from "@tanstack/react-table";
import { BookDetail } from "../state/messages";

export interface Props {
  books: BookDetail[];
}
const columnHelper = createColumnHelper<BookDetail>();

const columns = [
  columnHelper.accessor("server", {
    header: "Server"
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
  })
  //   columnHelper.accessor("full", {
  //     header: "Status",
  //   })
  //   columnHelper.display({
  //     header: "Download"
  //   })
];

const useStyles = createStyles((theme) => ({
  row: {
    padding: theme.spacing.md,
    margin: theme.spacing.md
  }
}));

export default function BookTable({ books: data }: Props) {
  const { classes } = useStyles();
  const table = useReactTable({
    data,
    columns: columns,
    getCoreRowModel: getCoreRowModel()
  });

  const headers = table.getHeaderGroups().map((headerGroup) => (
    <tr key={headerGroup.id}>
      {headerGroup.headers.map((header) => (
        <th key={header.id}>
          {header.isPlaceholder
            ? null
            : flexRender(header.column.columnDef.header, header.getContext())}
        </th>
      ))}
    </tr>
  ));

  const rows = table.getRowModel().rows.map((row) => (
    <tr key={row.id}>
      {row.getVisibleCells().map((cell) => (
        <td key={cell.id}>
          {flexRender(cell.column.columnDef.cell, cell.getContext())}
        </td>
      ))}
    </tr>
  ));

  return (
    <Table
      highlightOnHover
      verticalSpacing="md"
      sx={(theme) => ({
        borderCollapse: "separate",
        borderRadius: theme.radius.md,
        border: `1px solid ${
          theme.colorScheme === "dark"
            ? theme.colors.dark[5]
            : theme.colors.gray[2]
        }`
      })}>
      <thead style={{ border: "1px solid gray" }}>{headers}</thead>
      <tbody style={{ border: "1px solid gray" }}>{rows}</tbody>
    </Table>
  );
}
