import {
  Badge,
  Box,
  Button,
  Center,
  CloseButton,
  Divider,
  Group,
  Popover,
  Text,
  TextInput,
  useMantineColorScheme
} from "@mantine/core";
import { Column } from "@tanstack/react-table";
import { useVirtualizer } from "@tanstack/react-virtual";
import { MagnifyingGlass, PlusCircle } from "@phosphor-icons/react";
import { CSSProperties, useRef, useState } from "react";
import classes from "./ToolbarFacetFilter.module.css";

const stringContains = (first: string, second: string): boolean => {
  return first.toLowerCase().includes(second.toLowerCase());
};

interface FacetFilterProps<TData> {
  placeholder: string;
  column: Column<TData, unknown>;
  Entry: React.FC<FacetEntryProps>;
}

export interface FacetEntryProps {
  entry: string;
  selected: boolean;
  onClick: (entry: string) => void;
  style: CSSProperties;
}

export default function ToolbarFacetFilter<TData>({
  placeholder,
  column,
  Entry
}: FacetFilterProps<TData>) {
  const [filter, setFilter] = useState("");
  const [opened, setOpened] = useState(false);

  const options = Array.from(column.getFacetedUniqueValues().keys());
  const filteredOptions = options.filter((x) => stringContains(x, filter));

  const listRef = useRef<HTMLDivElement>(null);

  const rowVirtualizer = useVirtualizer({
    count: filteredOptions.length,
    getScrollElement: () => listRef.current,
    estimateSize: () => 30,
    overscan: 10
  });

  const filterValue = (column.getFilterValue() ?? []) as string[];

  const { colorScheme } = useMantineColorScheme();

  const buttonColor =
    colorScheme === "dark"
      ? filterValue.length > 0
        ? "blue.2"
        : "dark.0"
      : filterValue.length > 0
      ? "blue.4"
      : "gray.7";

  return (
    <Popover
      width={200}
      trapFocus
      position="bottom-start"
      withArrow
      shadow="md"
      opened={opened}
      onChange={setOpened}
      styles={{ dropdown: { padding: 0 } }}>
      <Popover.Target>
        <Button
          variant="default"
          className={classes.button}
          size="xs"
          leftSection={<PlusCircle size={16} />}
          color={buttonColor}
          onClick={() => setOpened((o) => !o)}>
          {placeholder}
          {filterValue.length > 0 && (
            <>
              <Divider my={4} ml={8} mr={4} orientation="vertical" />
              <Badge hiddenFrom="md" color={buttonColor} radius="sm" size="xs">
                {filterValue.length}
              </Badge>
              <Box visibleFrom="md">
                {filterValue.length > 2 ? (
                  <Badge color={buttonColor} radius="sm" size="xs">
                    {filterValue.length} selected
                  </Badge>
                ) : (
                  filterValue.map((value, index) => (
                    <Badge
                      key={index}
                      radius="sm"
                      color={buttonColor}
                      size="xs"
                      ml={4}>
                      {value}
                    </Badge>
                  ))
                )}
              </Box>
            </>
          )}
        </Button>
      </Popover.Target>
      <Popover.Dropdown>
        <Group justify="space-between" className={classes.header}>
          <Text fw="normal" size="xs" c="dark">
            Filter {placeholder}
          </Text>
          <CloseButton
            onClick={() => setOpened(false)}
            aria-label="Close modal"
            color="dark"
            iconSize={12}
          />
        </Group>
        <TextInput
          data-autofocus
          leftSection={<MagnifyingGlass weight="bold" />}
          className={classes.search}
          variant="unstyled"
          value={filter}
          size="xs"
          onChange={(e) => setFilter(e.currentTarget.value)}
          placeholder="Filter..."
          rightSection={
            column.getIsFiltered() && (
              <Button
                style={{ marginRight: 25, fontWeight: "normal" }}
                size="compact-xs"
                variant="default"
                color="blue"
                disabled={filterValue.length === 0}
                onClick={() => {
                  column.setFilterValue([]);
                  setFilter("");
                }}>
                Reset
              </Button>
            )
          }
        />

        <div ref={listRef} className={classes.container}>
          <div
            style={{
              height: `${rowVirtualizer.getTotalSize()}px`,
              width: "100%",
              position: "relative"
            }}>
            {rowVirtualizer.getVirtualItems().map((virtualRow) => {
              const selected = filterValue.includes(
                filteredOptions[virtualRow.index]
              );
              return (
                <Entry
                  key={virtualRow.index}
                  style={{
                    position: "absolute",
                    top: 0,
                    left: 0,
                    width: "100%",
                    height: virtualRow.size,
                    transform: `translateY(${virtualRow.start}px)`
                  }}
                  entry={filteredOptions[virtualRow.index]}
                  selected={selected}
                  onClick={(entry) =>
                    selected
                      ? column.setFilterValue(
                          filterValue.filter(
                            (x) => x !== filteredOptions[virtualRow.index]
                          )
                        )
                      : column.setFilterValue([...filterValue, entry])
                  }
                />
              );
            })}
          </div>

          {filteredOptions.length === 0 && (
            <Center py="xs">
              <Text color="dimmed" size="xs">
                No Results
              </Text>
            </Center>
          )}
        </div>
      </Popover.Dropdown>
    </Popover>
  );
}
