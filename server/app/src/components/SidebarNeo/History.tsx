import { Badge, Button, Menu, Stack, Tooltip } from "@mantine/core";
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
import { useSidebarButtonStyle } from "./styles";

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

const HistoryCard: React.FC<Props> = ({ activeTS, item, dispatch }: Props) => {
  const isActive = activeTS === item.timestamp;
  const { classes } = useSidebarButtonStyle({ isActive });

  const loading = !item.results?.length && !item.errors?.length;

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
            rightIcon={
              <Badge color="brand" radius="sm" size="sm" variant="light">
                {`${item.results?.length} RESULTS`}
              </Badge>
            }
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
