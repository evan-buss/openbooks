import { ScrollArea, Table, Text } from "@mantine/core";
import { useElementSize, useMergedRef, useTextSelection } from "@mantine/hooks";
import {
  createColumnHelper,
  flexRender,
  getCoreRowModel,
  getFilteredRowModel,
  Row,
  useReactTable
} from "@tanstack/react-table";
import { useVirtualizer } from "@tanstack/react-virtual";
import { MagnifyingGlass, WarningCircle } from "phosphor-react";
import { useEffect, useMemo, useRef } from "react";
import { ParseError } from "../../state/messages";
import { TextFilter } from "./Filters/TextFilter";
import { useTableStyles } from "./styles";

const columnHelper = createColumnHelper<ParseError>();

export default function ErrorTable({
  errors,
  setSearchQuery
}: {
  errors: ParseError[];
  setSearchQuery: (query: string) => void;
}) {
  const selection = useTextSelection();
  const selectionText = selection?.toString() ?? "";

  useEffect(() => {
    setSearchQuery(selectionText);
  }, [selectionText]);

  const { classes, cx, theme } = useTableStyles();
  const { ref: elementSizeRef, height, width } = useElementSize();
  const virtualizerRef = useRef();
  const mergedRef = useMergedRef(elementSizeRef, virtualizerRef);

  const columns = useMemo(() => {
    const cols = (cols: number) => (width / 12) * cols;
    return [
      columnHelper.accessor("line", {
        header: (props) => (
          <TextFilter
            icon={<MagnifyingGlass weight="bold" />}
            placeholder="Line"
            column={props.column}
            table={props.table}
          />
        ),
        cell: (props) => <pre style={{ margin: 0 }}>{props.getValue()}</pre>,
        size: cols(9),
        enableColumnFilter: true
      }),
      columnHelper.accessor("error", {
        header: (props) => (
          <TextFilter
            icon={<WarningCircle weight="bold" />}
            placeholder="Error"
            column={props.column}
            table={props.table}
          />
        ),
        size: cols(3),
        enableColumnFilter: false
      })
    ];
  }, [width]);

  const table = useReactTable({
    data: errors,
    columns: columns,
    enableFilters: true,
    columnResizeMode: "onChange",
    getCoreRowModel: getCoreRowModel(),
    getFilteredRowModel: getFilteredRowModel()
  });

  const { rows: tableRows } = table.getRowModel();

  const rowVirtualizer = useVirtualizer({
    count: tableRows.length,
    getScrollElement: () => virtualizerRef.current,
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
      <Text size="sm" color="dimmed" className="w-full" mb={4}>
        These results could not be parsed to due to their non-standard format.
        To download, copy the line up to the <code>::INFO::</code> or file size
        at the end and paste into the text box above.
      </Text>

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
                    style={{
                      color:
                        theme.colorScheme === "dark"
                          ? theme.colors.dark[3]
                          : theme.colors.dark[1],
                      width: header.getSize(),
                      position: "relative",
                      textTransform: "uppercase"
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
              ] as unknown as Row<ParseError>;
              return (
                <tr key={row.id} style={{ height: 50 }}>
                  {row.getVisibleCells().map((cell) => {
                    return (
                      <td key={cell.id}>
                        <Text lineClamp={1} color="dark">
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
