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
import { useGetServersQuery } from "../state/api";
import { BookDetail } from "../state/messages";

interface FacetFilterProps {
  placeholder: string;
  column: Column<BookDetail, string>;
  table: Table<BookDetail>;
}

const stringContains = (first: string, second: string): boolean => {
  return first.toLowerCase().includes(second.toLowerCase());
};

const useStyles = createStyles((theme) => {
  const border = `1px solid ${
    theme.colorScheme === "dark" ? theme.colors.dark[3] : theme.colors.gray[3]
  }`;

  return {
    header: {
      padding: 8
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
      overflow: "auto"
    },
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
      ["&:hover, &:focus"]: {
        backgroundColor:
          theme.colorScheme === "dark"
            ? theme.colors.dark[7]
            : theme.colors.gray[0]
      }
    },
    indicator: {
      width: 2,
      height: "100%",
      position: "absolute",
      top: 0,
      left: 0,
      backgroundColor:
        theme.colorScheme === "dark"
          ? theme.colors.brand[3]
          : theme.colors.brand[4],
      borderRadius: theme.radius.xl
    }
  };
});

export default function FacetFilter({
  placeholder,
  column,
  table
}: FacetFilterProps) {
  const [filter, setFilter] = useState("");
  const { classes } = useStyles();
  const { data: onlineServers, refetch } = useGetServersQuery(null);

  const facets = Array.from(column.getFacetedUniqueValues().keys());
  const visible = facets.filter((x) => stringContains(x, filter));
  const [opened, setOpened] = useState(false);
  const listRef = useRef<HTMLDivElement>(null);

  const rowVirtualizer = useVirtualizer({
    count: visible.length,
    getScrollElement: () => listRef.current,
    estimateSize: () => 30,
    overscan: 10
  });

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
          compact
          color="dark.4"
          uppercase
          onClick={() => setOpened((o) => !o)}
          rightIcon={<CaretDown weight="bold" />}>
          {placeholder}
        </Button>
      </Popover.Target>
      <Popover.Dropdown>
        <Group position="apart" className={classes.header}>
          <Text weight="normal" size="xs">
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
        />

        <div ref={listRef} className={classes.container} data-autofocus>
          <div
            style={{
              height: `${rowVirtualizer.getTotalSize()}px`,
              width: "100%",
              position: "relative"
            }}>
            {rowVirtualizer.getVirtualItems().map((virtualRow) => {
              const selected =
                column.getFilterValue() === visible[virtualRow.index];
              return (
                <FacetEntry
                  key={virtualRow.index}
                  style={{
                    position: "absolute",
                    top: 0,
                    left: 0,
                    width: "100%",
                    height: virtualRow.size,
                    transform: `translateY(${virtualRow.start}px)`
                  }}
                  entry={visible[virtualRow.index]}
                  selected={selected}
                  onClick={(entry) =>
                    selected
                      ? column.setFilterValue(null)
                      : column.setFilterValue(entry)
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

interface FacetEntryProps {
  entry: string;
  selected: boolean;
  onClick: (entry: string) => void;
  style: CSSProperties;
}

function FacetEntry({ entry, onClick, selected, style }: FacetEntryProps) {
  const { classes, cx } = useStyles();
  return (
    <Box
      tabIndex={0}
      className={classes.entry}
      style={style}
      onClick={() => onClick(entry)}>
      <div className={cx({ [classes.indicator]: selected })}></div>

      <Text size={12} weight="normal" style={{ marginLeft: 20 }}>
        <Indicator position="middle-start" offset={-16} size={6} color="green">
          {entry}
        </Indicator>
      </Text>
    </Box>
  );
}
