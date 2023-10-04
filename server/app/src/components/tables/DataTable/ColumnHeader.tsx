import { Column } from "@tanstack/react-table";
import { Button, Menu, Text } from "@mantine/core";
import {
  CaretUpDown,
  SortAscending,
  SortDescending,
  XCircle
} from "@phosphor-icons/react";
import React from "react";

interface DataTableColumnHeaderProps<TData, TValue>
  extends React.HTMLAttributes<HTMLDivElement> {
  column: Column<TData, TValue>;
  title: string;
  icon?: React.ReactNode;
}

export function DataTableColumnHeader<TData, TValue>({
  column,
  title,
  icon
}: DataTableColumnHeaderProps<TData, TValue>) {
  if (!column.getCanSort()) {
    return (
      <Text
        style={{ color: "var(--mantine-color-dark-light-color)" }}
        fw={500}
        mx={20}
        size="xs"
        tt="uppercase">
        {title}
      </Text>
    );
  }

  return (
    <Menu shadow="md">
      <Menu.Target>
        <Button
          color="dark"
          leftSection={icon}
          tt="uppercase"
          size="xs"
          variant="subtle"
          rightSection={
            column.getIsSorted() === "desc" ? (
              <SortDescending className="ml-2 h-4 w-4" />
            ) : column.getIsSorted() === "asc" ? (
              <SortAscending className="ml-2 h-4 w-4" />
            ) : (
              <CaretUpDown className="ml-2 h-4 w-4" />
            )
          }>
          {title}
        </Button>
      </Menu.Target>

      <Menu.Dropdown>
        <Menu.Item
          leftSection={<SortAscending size={14} />}
          onClick={() => column.toggleSorting(false)}>
          Ascending
        </Menu.Item>
        <Menu.Item
          leftSection={<SortDescending size={14} />}
          onClick={() => column.toggleSorting(true)}>
          Descending
        </Menu.Item>
        {column.getIsSorted() !== false && (
          <Menu.Item
            leftSection={<XCircle size={14} />}
            onClick={() => column.clearSorting()}>
            Clear Sort
          </Menu.Item>
        )}
      </Menu.Dropdown>
    </Menu>
  );
}
