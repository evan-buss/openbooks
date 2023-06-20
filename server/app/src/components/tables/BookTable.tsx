import { Button, Indicator, Loader, Text, Tooltip } from "@mantine/core";
import { createColumnHelper, Table } from "@tanstack/react-table";
import React, { useState } from "react";
import { useSelector } from "react-redux";
import { useGetServersQuery } from "../../state/api";
import { BookDetail } from "../../state/messages";
import { sendDownload } from "../../state/stateSlice";
import { RootState, useAppDispatch } from "../../state/store";
import { DataTableColumnHeader } from "./DataTable/ColumnHeader";
import DataTable, { ColumnWidthFunc } from "./DataTable/DataTable";
import ToolbarFacetFilter from "./DataTable/ToolbarFacetFilter";
import { ServerFacetEntry, StandardFacetEntry } from "./Facets";

const columnHelper = createColumnHelper<BookDetail>();

interface BookTableProps {
  books: BookDetail[];
}

export default function BookTable({ books }: BookTableProps) {
  const { data: servers } = useGetServersQuery(null);

  const columns = (cols: ColumnWidthFunc) => [
    columnHelper.accessor("server", {
      header: (props) => (
        <DataTableColumnHeader column={props.column} title="Server" />
      ),
      cell: (props) => {
        const online = servers?.includes(props.getValue());
        return (
          <Text size={12} weight="normal" color="dark" ml={14}>
            <Tooltip position="top-start" label={online ? "Online" : "Offline"}>
              <Indicator
                zIndex={0}
                position="middle-start"
                offset={-16}
                size={6}
                color={online ? "green.6" : "gray"}>
                {props.getValue()}
              </Indicator>
            </Tooltip>
          </Text>
        );
      },
      size: cols(1),
      enableGlobalFilter: false,
      filterFn: "arrIncludesSome"
    }),
    columnHelper.accessor("author", {
      header: (props) => (
        <DataTableColumnHeader column={props.column} title="Author" />
      ),
      size: cols(2),
      cell: (props) => (
        <Tooltip label={props.getValue()} position="bottom">
          <Text size={12} weight="normal" color="dark">
            {props.getValue()}
          </Text>
        </Tooltip>
      )
    }),
    columnHelper.accessor("title", {
      header: (props) => (
        <DataTableColumnHeader column={props.column} title="Title" />
      ),
      cell: (props) => (
        <Tooltip label={props.getValue()} position="bottom-start">
          <Text size={12} weight="normal" color="dark">
            {props.getValue()}
          </Text>
        </Tooltip>
      ),
      minSize: 20,
      size: cols(6)
    }),
    columnHelper.accessor("format", {
      header: (props) => (
        <DataTableColumnHeader column={props.column} title="Format" />
      ),
      size: cols(1),
      filterFn: "arrIncludesSome",
      enableGlobalFilter: false
    }),
    columnHelper.accessor("size", {
      header: (props) => (
        <DataTableColumnHeader column={props.column} title="Size" />
      ),
      size: cols(1),
      enableGlobalFilter: false
    }),
    columnHelper.accessor("full", {
      header: (props) => (
        <DataTableColumnHeader column={props.column} title="Download" />
      ),
      size: cols(1),
      enableGlobalFilter: false,
      enableSorting: false,
      enableHiding: false,
      cell: (props) => <DownloadButton book={props.getValue()}></DownloadButton>
    })
  ];

  const facetFilters = (table: Table<BookDetail>) => [
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
  ];

  return (
    <DataTable columns={columns} data={books} facetFilters={facetFilters} />
  );
}

function DownloadButton({ book }: { book: string }) {
  const dispatch = useAppDispatch();

  const [clicked, setClicked] = useState(false);
  const isInFlight = useSelector((state: RootState) =>
    state.state.inFlightDownloads.includes(book)
  );

  // Prevent hitting the same button multiple times
  const onClick = () => {
    if (clicked) return;
    dispatch(sendDownload(book));
    setClicked(true);
  };

  return (
    <Button
      compact
      size="xs"
      radius="sm"
      onClick={onClick}
      sx={{ fontWeight: "normal", width: 80 }}>
      {isInFlight ? (
        <Loader variant="dots" color="gray" />
      ) : (
        <span>Download</span>
      )}
    </Button>
  );
}
