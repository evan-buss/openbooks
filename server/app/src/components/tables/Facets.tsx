import cx from "clsx";
import { Button, Indicator, Text } from "@mantine/core";
import { FacetEntryProps } from "./DataTable/ToolbarFacetFilter";
import { useGetServersQuery } from "../../state/api";
import React from "react";
import classes from "./Facets.module.css";
import { conditionalAttribute } from "../../utils/attribute-helper";

export function ServerFacetEntry({
  entry,
  onClick,
  selected,
  style
}: FacetEntryProps) {
  const { data: servers } = useGetServersQuery(null);
  const serverOnline = servers && servers.includes(entry);

  return (
    <Button
      variant="subtle"
      tabIndex={0}
      {...conditionalAttribute("selected", selected)}
      className={classes.listEntry}
      style={style}
      onClick={() => onClick(entry)}>
      <div className={cx({ [classes.indicator]: selected })}></div>

      <Text span fz={12} fw="normal" c="dark" style={{ marginLeft: 20 }}>
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
  return (
    <Button
      variant="subtle"
      tabIndex={0}
      {...conditionalAttribute("selected", selected)}
      className={classes.listEntry}
      style={style}
      onClick={() => onClick(entry)}>
      <div className={cx({ [classes.indicator]: selected })}></div>

      <Text fz={12} fw="normal" c="dark">
        {entry}
      </Text>
    </Button>
  );
}
