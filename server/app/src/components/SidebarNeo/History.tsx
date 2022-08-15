import {
  Badge,
  Button,
  createStyles,
  Menu,
  Stack,
  Tooltip
} from "@mantine/core";
import { Dispatch } from "@reduxjs/toolkit";
import { Eye, EyeSlash, MagnifyingGlass, Trash } from "phosphor-react";
import React from "react";
import { useSelector } from "react-redux";
import {
  deleteHistoryItem,
  HistoryItem,
  selectHistory
} from "../../state/historySlice";
import { setActiveItem } from "../../state/stateSlice";
import { RootState, useAppDispatch } from "../../state/store";

export default function SearchHistoryNeo() {
  const history = useSelector(selectHistory);
  const activeTS =
    useSelector((store: RootState) => store.state.activeItem?.timestamp) ?? -1;
  const dispatch = useAppDispatch();

  return (
    <Stack spacing="xs" sx={{ width: "278px" }}>
      {history.length > 0 ? (
        history.map((item: HistoryItem) => (
          <HistoryCard
            activeTS={activeTS}
            key={item.timestamp.toString()}
            item={item}
            dispatch={dispatch}
          />
        ))
      ) : (
        <p className="my-4 text-sm text-center text-gray-500">
          History is a mystery.
        </p>
      )}
    </Stack>
  );
}

type Props = {
  activeTS: number;
  item: HistoryItem;
  dispatch: Dispatch<any>;
};

const useHistoryCardStyles = createStyles(
  (theme, { active }: { active: boolean }) => {
    const isDark = theme.colorScheme === "dark";

    return {
      root: {
        "backgroundColor": isDark ? theme.colors.dark[4] : "white",
        "borderColor": active
          ? theme.fn.primaryColor()
          : isDark
          ? theme.colors.gray[8]
          : theme.colors.gray[3],
        "boxShadow": active ? theme.shadows.sm : "none",

        "&:hover": {
          backgroundColor: isDark ? theme.colors.dark[5] : theme.colors.gray[1]
        }
      },
      inner: {
        color: isDark ? "white" : "black",
        fontWeight: "normal",
        justifyContent: "space-between"
      },
      label: {
        paddingLeft: theme.spacing.sm,
        width: "100%",
        textAlign: "start"
      }
    };
  }
);

const HistoryCard: React.FC<Props> = ({ activeTS, item, dispatch }: Props) => {
  const isActive = activeTS === item.timestamp;
  const { classes } = useHistoryCardStyles({ active: isActive });

  //   const loading = !item.results?.length && !item.errors?.length;
  const loading = false;

  const badge = loading ? (
    <></>
  ) : (
    <Badge
      sx={(theme) => ({
        color:
          theme.colorScheme === "dark"
            ? theme.colors.brand[2]
            : theme.colors.brand[5]
      })}
      color="brand"
      radius="sm"
      size="sm"
      variant="light">{`${item.results?.length} RESULTS`}</Badge>
  );

  return (
    <Menu shadow="md" width={200}>
      <Menu.Target>
        <Tooltip label={item.query} openDelay={1_000}>
          <Button
            classNames={classes}
            loading={loading}
            radius="sm"
            variant="outline"
            leftIcon={<MagnifyingGlass size={18} weight="bold" />}
            rightIcon={badge}
            loaderPosition="right"
            loaderProps={{ color: "brand" }}>
            {item.query}
          </Button>
        </Tooltip>
      </Menu.Target>

      <Menu.Dropdown>
        {!isActive ? (
          <Menu.Item
            icon={<Eye size={18} weight="bold" />}
            onClick={() => dispatch(setActiveItem(item))}>
            Show Results
          </Menu.Item>
        ) : (
          <Menu.Item
            icon={<EyeSlash size={18} weight="bold" />}
            onClick={() => dispatch(setActiveItem(null))}>
            Hide Results
          </Menu.Item>
        )}

        <Menu.Item
          color="red"
          icon={<Trash size={18} weight="bold" />}
          onClick={() => dispatch(deleteHistoryItem(item.timestamp))}>
          Delete item
        </Menu.Item>
      </Menu.Dropdown>
    </Menu>
  );
};
