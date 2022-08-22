import { TextInput } from "@mantine/core";
import { getHotkeyHandler, useDebouncedValue } from "@mantine/hooks";
import { Column, Table } from "@tanstack/react-table";
import { useEffect, useState } from "react";
import { BookDetail } from "../state/messages";

interface TextFilterProps {
  icon?: React.ReactNode;
  placeholder: string;
  column: Column<BookDetail, string>;
  table: Table<BookDetail>;
}

export function TextFilter({
  icon,
  placeholder,
  column,
  table
}: TextFilterProps) {
  const [filterValue, setFilterValue] = useState(
    column.getFilterValue() as string
  );
  const [debounced] = useDebouncedValue(filterValue, 500);

  useEffect(() => {
    column.setFilterValue(debounced);
  }, [debounced]);

  return (
    <TextInput
      icon={icon}
      size="xs"
      sx={(theme) => ({
        ["& input::placeholder"]: {
          color: theme.colors.dark[4],
          fontWeight: "bold",
          textTransform: "uppercase"
        }
      })}
      placeholder={placeholder}
      variant="unstyled"
      onChange={(e) => setFilterValue(e.currentTarget.value)}
      value={filterValue}
      onKeyDown={getHotkeyHandler([["Escape", () => setFilterValue("")]])}
    />
  );
}
