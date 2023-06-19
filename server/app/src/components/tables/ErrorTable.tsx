import {
  useDebouncedValue,
  useDisclosure,
  useTextSelection
} from "@mantine/hooks";
import { ActionIcon, Popover, Text } from "@mantine/core";
import { createColumnHelper, Table } from "@tanstack/react-table";
import React, { useEffect } from "react";
import { ParseError } from "../../state/messages";
import DataTable, { ColumnWidthFunc } from "./DataTable/DataTable";
import { DataTableColumnHeader } from "./DataTable/ColumnHeader";
import ToolbarFacetFilter from "./DataTable/ToolbarFacetFilter";
import { StandardFacetEntry } from "./Facets";
import { Info } from "@phosphor-icons/react";

const columnHelper = createColumnHelper<ParseError>();

interface ErrorTableProps {
  errors: ParseError[];
  setSearchQuery: (query: string) => void;
}

export default function ErrorTable({
  errors,
  setSearchQuery
}: ErrorTableProps) {
  const selection = useTextSelection();
  const [selectionText] = useDebouncedValue(selection?.toString() ?? "", 250);

  useEffect(() => {
    setSearchQuery(selectionText);
  }, [selectionText]);

  const columns = (cols: ColumnWidthFunc) => [
    columnHelper.accessor("line", {
      header: (props) => (
        <DataTableColumnHeader column={props.column} title="Line" />
      ),
      cell: (props) => <code style={{ margin: 0 }}>{props.getValue()}</code>,
      size: cols(9)
    }),
    columnHelper.accessor("error", {
      header: (props) => (
        <DataTableColumnHeader column={props.column} title="Error" />
      ),
      size: cols(3),
      filterFn: "arrIncludesSome"
    })
  ];

  const [opened, { close, open }] = useDisclosure(false);

  const facetFilters = (table: Table<ParseError>) => [
    <ToolbarFacetFilter
      key="error"
      placeholder="Error"
      column={table.getColumn("error")!}
      Entry={StandardFacetEntry}
    />,
    <Popover
      key="info-text"
      width={300}
      position="bottom"
      shadow="md"
      opened={opened}>
      <Popover.Target>
        <ActionIcon
          color="brand"
          variant="light"
          onMouseEnter={open}
          onMouseLeave={close}>
          <Info size={18} />
        </ActionIcon>
      </Popover.Target>
      <Popover.Dropdown sx={{ pointerEvents: "none" }} w={400}>
        <Text size="sm">
          {errors?.length === 1 ? "This result" : "These results"} could not be
          parsed to due to {errors?.length === 1 ? "its" : "their"} non-standard
          format. To download, copy the line up to the <code>::INFO::</code> or
          file size at the end and paste into the text box above.
        </Text>
      </Popover.Dropdown>
    </Popover>
  ];

  return (
    <>
      <DataTable columns={columns} data={errors} facetFilters={facetFilters} />
    </>
  );
}
