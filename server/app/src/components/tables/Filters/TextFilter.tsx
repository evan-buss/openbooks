import { Box, TextInput } from "@mantine/core";
import { getHotkeyHandler } from "@mantine/hooks";
import { Column, Table } from "@tanstack/react-table";
import { useEffect, useState } from "react";

interface TextFilterProps {
  icon?: React.ReactNode;
  placeholder: string;
  column: Column<any, string>;
  table: Table<any>;
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

  useEffect(() => {
    column.setFilterValue(filterValue);
  }, [filterValue]);

  const styledIcon = (
    <Box
      component="span"
      sx={(theme) => ({
        display: "flex",
        color:
          theme.colorScheme === "dark"
            ? filterValue
              ? theme.colors.brand[3]
              : theme.colors.dark[3]
            : filterValue
            ? theme.colors.brand[4]
            : theme.colors.dark[1]
      })}>
      {icon}
    </Box>
  );

  return (
    <TextInput
      icon={styledIcon}
      size="xs"
      placeholder={placeholder}
      styles={(theme) => ({
        input: {
          ["&::placeholder"]: {
            color:
              theme.colorScheme === "dark"
                ? theme.colors.dark[0]
                : theme.colors.gray[7],
            textTransform: "uppercase",
            fontWeight: "bold"
          }
        }
      })}
      variant="unstyled"
      onChange={(e) => setFilterValue(e.currentTarget.value)}
      value={filterValue}
      onKeyDown={getHotkeyHandler([["Escape", () => setFilterValue("")]])}
    />
  );
}
