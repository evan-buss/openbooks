import { Button, createStyles, Indicator, Text } from "@mantine/core";
import { FacetEntryProps } from "./DataTable/ToolbarFacetFilter";
import { useGetServersQuery } from "../../state/api";
import React from "react";

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
    <Button
      variant="subtle"
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
    </Button>
  );
}

export function StandardFacetEntry({
  entry,
  onClick,
  selected,
  style
}: FacetEntryProps): React.ReactNode {
  const { classes, cx } = useFacetStyles();
  return (
    <Button
      variant="subtle"
      tabIndex={0}
      className={cx(classes.entry, { [classes.entrySelected]: selected })}
      style={style}
      onClick={() => onClick(entry)}>
      <div className={cx({ [classes.indicator]: selected })}></div>

      <Text size={12} weight="normal" color="dark">
        {entry}
      </Text>
    </Button>
  );
}
