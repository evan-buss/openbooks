import { createStyles } from "@mantine/core";

export const useTableStyles = createStyles((theme) => ({
  container: {
    border: `1px solid ${
      theme.colorScheme === "dark" ? theme.colors.dark[3] : theme.colors.gray[3]
    }`,
    borderRadius: theme.radius.md,
    backgroundColor:
      theme.colorScheme === "dark" ? theme.colors.dark[6] : "white",
    height: "100%",
    overflow: "auto",
    width: "100%",
    boxShadow: theme.shadows.xs
  },
  head: {
    position: "sticky",
    top: 0,
    backgroundColor:
      theme.colorScheme === "dark"
        ? theme.colors.dark[5]
        : theme.colors.gray[1],
    zIndex: 1
  },
  headerCell: {
    textTransform: "uppercase",
    position: "relative"
  },
  resizer: {
    position: "absolute",
    right: 0,
    top: 0,
    height: "100%",
    width: "2px",
    background:
      theme.colorScheme === "dark"
        ? theme.colors.dark[3]
        : theme.colors.gray[6],
    cursor: "col-resize",
    userSelect: "none",
    touchAction: "none",
    opacity: 0,

    ["&.isResizing"]: {
      background:
        theme.colorScheme === "dark"
          ? theme.colors.brand[3]
          : theme.colors.brand[4],
      opacity: 1
    },

    ["&:hover"]: {
      opacity: 1
    }
  }
}));
