import {
  Badge,
  Box,
  Button,
  Center,
  CloseButton,
  createStyles,
  Divider,
  Group,
  Popover,
  Text,
  TextInput
} from "@mantine/core";
import { Column } from "@tanstack/react-table";
import { useVirtualizer } from "@tanstack/react-virtual";
import { MagnifyingGlass, PlusCircle } from "@phosphor-icons/react";
import { CSSProperties, useRef, useState } from "react";

const stringContains = (first: string, second: string): boolean => {
  return first.toLowerCase().includes(second.toLowerCase());
};

const useStyles = createStyles((theme) => {
  const border = `1px solid ${
    theme.colorScheme === "dark" ? theme.colors.dark[4] : theme.colors.gray[3]
  }`;

  return {
    header: {
      padding: 6,
      textTransform: "none"
    },
    search: {
      borderTop: border,
      borderBottom: border,
      background:
        theme.colorScheme === "dark"
          ? theme.colors.dark[5]
          : theme.colors.gray[0]
    },
    container: {
      maxHeight: 200,
      overflow: "auto",
      textTransform: "none"
    },
    button: {
      border: `1px dashed ${
        theme.colorScheme === "dark"
          ? theme.colors.dark[3]
          : theme.colors.gray[3]
      }`,
      borderRadius: theme.radius.md,
      backgroundColor:
        theme.colorScheme === "dark" ? theme.colors.dark[6] : "white",
      ["&:hover"]: {
        backgroundColor:
          theme.colorScheme === "dark"
            ? theme.colors.dark[4]
            : theme.colors.gray[0]
      }
    }
  };
});

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

  const { classes, theme } = useStyles();
  const listRef = useRef<HTMLDivElement>(null);

  const rowVirtualizer = useVirtualizer({
    count: filteredOptions.length,
    getScrollElement: () => listRef.current,
    estimateSize: () => 30,
    overscan: 10
  });

  const filterValue = (column.getFilterValue() ?? []) as string[];

  const buttonColor =
    theme.colorScheme === "dark"
      ? filterValue.length > 0
        ? "brand.2"
        : "dark.0"
      : filterValue.length > 0
      ? "brand.4"
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
          variant=""
          className={classes.button}
          size="xs"
          leftIcon={<PlusCircle size={16} />}
          color={buttonColor}
          onClick={() => setOpened((o) => !o)}>
          {placeholder}
          {filterValue.length > 0 && (
            <>
              <Divider my={4} ml={8} mr={4} orientation="vertical" />
              <Box>
                {filterValue.length > 2 ? (
                  <Badge color={buttonColor} size="xs">
                    {filterValue.length} selected
                  </Badge>
                ) : (
                  filterValue.map((value, index) => (
                    <Badge key={index} color={buttonColor} size="xs" ml={4}>
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
        <Group position="apart" className={classes.header}>
          <Text weight="normal" size="xs" color="dark">
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
          icon={<MagnifyingGlass weight="bold" />}
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
                size="xs"
                compact
                variant="default"
                color="brand"
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

        <div ref={listRef} className={classes.container} data-autofocus>
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
