import {
  Box,
  Button,
  CloseButton,
  createStyles,
  Group,
  Indicator,
  Popover,
  Text,
  TextInput
} from "@mantine/core";
import { Column, Table } from "@tanstack/react-table";
import { useVirtualizer } from "@tanstack/react-virtual";
import { CaretDown, MagnifyingGlass } from "phosphor-react";
import { CSSProperties, useRef, useState } from "react";
import { useGetServersQuery } from "../../../state/api";

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
      ["&:hover"]: {
        backgroundColor:
          theme.colorScheme === "dark"
            ? theme.colors.dark[4]
            : theme.colors.gray[0]
      }
    }
  };
});

interface FacetFilterProps {
  placeholder: string;
  column: Column<any, string>;
  table: Table<any>;
  Entry: React.FC<FacetEntryProps>;
}

export default function FacetFilter({
  placeholder,
  column,
  table,
  Entry
}: FacetFilterProps) {
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
      position="bottom"
      withArrow
      shadow="md"
      opened={opened}
      onChange={setOpened}
      styles={{ dropdown: { padding: 0 } }}>
      <Popover.Target>
        <Button
          variant="subtle"
          size="xs"
          className={classes.button}
          compact
          uppercase
          color={buttonColor}
          onClick={() => setOpened((o) => !o)}
          rightIcon={<CaretDown weight="bold" />}>
          {placeholder}
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
                onClick={() => column.setFilterValue([])}>
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
        </div>
      </Popover.Dropdown>
    </Popover>
  );
}

const useFacetStyles = createStyles((theme) => {
  const border = `1px solid ${
    theme.colorScheme === "dark" ? theme.colors.dark[4] : theme.colors.gray[3]
  }`;

  return {
    entry: {
      ...theme.fn.focusStyles(),
      height: 30,
      position: "relative",
      display: "flex",
      alignItems: "center",
      justifyContent: "start",
      padding: "0 8px",
      cursor: "pointer",
      borderBottom: border,
      userSelect: "none",
      ["&:hover, &:focus"]: {
        backgroundColor:
          theme.colorScheme === "dark"
            ? theme.colors.dark[7]
            : theme.colors.gray[0]
      }
    },
    entrySelected: {
      backgroundColor:
        theme.colorScheme === "dark"
          ? theme.colors.dark[7]
          : theme.colors.gray[0]
    },
    indicator: {
      width: 2,
      height: "80%",
      position: "absolute",
      left: 0,
      backgroundColor:
        theme.colorScheme === "dark"
          ? theme.colors.brand[3]
          : theme.colors.brand[4],
      borderRadius: theme.radius.xl
    }
  };
});

export interface FacetEntryProps {
  entry: string;
  selected: boolean;
  onClick: (entry: string) => void;
  style: CSSProperties;
}

export function ServerFacetEntry({
  entry,
  onClick,
  selected,
  style
}: FacetEntryProps) {
  const { classes, cx } = useFacetStyles();
  const { data: servers } = useGetServersQuery(null);
  const serverOnline = servers?.includes(entry) ?? false;

  return (
    <Box
      tabIndex={0}
      className={cx(classes.entry, { [classes.entrySelected]: selected })}
      style={style}
      onClick={() => onClick(entry)}>
      <div className={cx({ [classes.indicator]: selected })}></div>

      <Text size={12} weight="normal" color="dark" style={{ marginLeft: 20 }}>
        <Indicator
          position="middle-start"
          offset={-16}
          size={6}
          color={serverOnline ? "green.6" : "gray"}>
          {entry}
        </Indicator>
      </Text>
    </Box>
  );
}

export function StandardFacetEntry({
  entry,
  onClick,
  selected,
  style
}: FacetEntryProps): JSX.Element {
  const { classes, cx } = useFacetStyles();
  return (
    <Box
      tabIndex={0}
      className={cx(classes.entry, { [classes.entrySelected]: selected })}
      style={style}
      onClick={() => onClick(entry)}>
      <div className={cx({ [classes.indicator]: selected })}></div>

      <Text size={12} weight="normal" color="dark">
        {entry}
      </Text>
    </Box>
  );
}
