import { Column } from "@tanstack/react-table";
import { Button, createStyles, Menu, Text } from "@mantine/core";
import {
  SortAscending,
  SortDescending,
  CaretUpDown,
  XCircle
} from "@phosphor-icons/react";
import React from "react";

const useStyles = createStyles((theme) => ({
  button: {
    color:
      theme.colorScheme === "dark" ? theme.colors.dark[0] : theme.colors.gray[7]
  }
}));

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
  const { classes } = useStyles();

  if (!column.getCanSort()) {
    return (
      <Text mx={20} transform="uppercase">
        {title}
      </Text>
    );
  }

  return (
    <Menu shadow="md">
      <Menu.Target>
        <Button
          color="dark"
          leftIcon={icon}
          uppercase
          className={classes.button}
          size="xs"
          variant="subtle"
          rightIcon={
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
          icon={<SortAscending size={14} />}
          onClick={() => column.toggleSorting(false)}>
          Ascending
        </Menu.Item>
        <Menu.Item
          icon={<SortDescending size={14} />}
          onClick={() => column.toggleSorting(true)}>
          Descending
        </Menu.Item>
        {column.getIsSorted() !== false && (
          <Menu.Item
            icon={<XCircle size={14} />}
            onClick={() => column.clearSorting()}>
            Clear Sort
          </Menu.Item>
        )}
      </Menu.Dropdown>
    </Menu>
  );
}
