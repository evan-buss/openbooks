import { Button, Indicator, Loader, Text, Tooltip } from "@mantine/core";
import { Table, createColumnHelper } from "@tanstack/react-table";
import { useCallback, useMemo, useState } from "react";
import { useSelector } from "react-redux";
import { useDownloadMutation, useGetServersQuery } from "../../state/api";
import { BookDetail } from "../../state/messages";
import { RootState } from "../../state/store";
import { DataTableColumnHeader } from "./DataTable/ColumnHeader";
import DataTable from "./DataTable/DataTable";
import ToolbarFacetFilter from "./DataTable/ToolbarFacetFilter";
import { ServerFacetEntry, StandardFacetEntry } from "./Facets";

const columnHelper = createColumnHelper<BookDetail>();

interface BookTableProps {
  books: BookDetail[];
}

export default function BookTable({ books }: BookTableProps) {
  const { data: servers } = useGetServersQuery(null);

  const columns = useMemo(
    () => [
      columnHelper.accessor("server", {
        header: (props) => (
          <DataTableColumnHeader column={props.column} title="Server" />
        ),
        cell: (props) => {
          const online = servers?.includes(props.getValue());
          return (
            <Tooltip position="top-start" label={online ? "Online" : "Offline"}>
              <Indicator
                zIndex={0}
                position="middle-start"
                offset={-16}
                size={6}
                color={online ? "green.6" : "gray"}>
                <p>{props.getValue()}</p>
              </Indicator>
            </Tooltip>
          );
        },
        size: 100,
        enableGlobalFilter: false,
        filterFn: "arrIncludesSome"
      }),
      columnHelper.accessor("author", {
        header: (props) => (
          <DataTableColumnHeader column={props.column} title="Author" />
        ),
        cell: (props) => (
          <Tooltip label={props.getValue()} position="bottom">
            <p>{props.getValue()}</p>
          </Tooltip>
        ),
        size: 200
      }),
      columnHelper.accessor("title", {
        header: (props) => (
          <DataTableColumnHeader column={props.column} title="Title" />
        ),
        cell: (props) => (
          <Tooltip label={props.getValue()} position="bottom-start">
            <Text fz={12} fw="normal" c="dark">
              {props.getValue()}
            </Text>
          </Tooltip>
        )
      }),
      columnHelper.accessor("format", {
        header: (props) => (
          <DataTableColumnHeader column={props.column} title="Format" />
        ),
        size: 100,
        filterFn: "arrIncludesSome",
        enableGlobalFilter: false
      }),
      columnHelper.accessor("size", {
        header: (props) => (
          <DataTableColumnHeader column={props.column} title="Size" />
        ),
        size: 100,
        enableGlobalFilter: false
      }),
      columnHelper.accessor("full", {
        header: (props) => (
          <DataTableColumnHeader column={props.column} title="Download" />
        ),
        size: 150,
        enableGlobalFilter: false,
        enableSorting: false,
        enableHiding: false,
        cell: (props) => (
          <DownloadButton book={props.getValue()}></DownloadButton>
        )
      })
    ],
    [servers]
  );

  const facetFilters = useCallback(
    (table: Table<BookDetail>) => [
      <ToolbarFacetFilter
        key="server"
        placeholder="Server"
        column={table.getColumn("server")!}
        Entry={ServerFacetEntry}
      />,
      <ToolbarFacetFilter
        key="format"
        placeholder="Format"
        column={table.getColumn("format")!}
        Entry={StandardFacetEntry}
      />
    ],
    []
  );

  return (
    <DataTable columns={columns} data={books} facetFilters={facetFilters} />
  );
}

function DownloadButton({ book }: { book: string }) {
  const [downloadMutation] = useDownloadMutation();

  const [clicked, setClicked] = useState(false);
  const isInFlight = useSelector((state: RootState) =>
    state.state.inFlightDownloads.includes(book)
  );

  // Prevent hitting the same button multiple times
  const onClick = () => {
    if (clicked) return;
    downloadMutation(book);
    setClicked(true);
  };

  return (
    <Button size="compact-xs" radius="sm" onClick={onClick} fw="normal" w={80}>
      {isInFlight ? <Loader type="dots" color="gray" /> : <span>Download</span>}
    </Button>
  );
}
