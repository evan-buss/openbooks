import {
  Badge,
  Button,
  Center,
  Loader,
  Menu,
  Stack,
  Text,
  Tooltip
} from "@mantine/core";
import { Dispatch } from "@reduxjs/toolkit";
import { AnimatePresence } from "framer-motion";
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
import AnimatedElement from "../AnimatedElement";
import { useSidebarButtonStyle } from "./styles";

export default function SearchHistoryNeo() {
  const history = useSelector(selectHistory);
  const activeTS =
    useSelector((store: RootState) => store.state.activeItem?.timestamp) ?? -1;
  const dispatch = useAppDispatch();

  return (
    <Stack spacing="xs">
      <AnimatePresence mode="popLayout">
        {history.length > 0 ? (
          history.map((item: HistoryItem) => (
            <AnimatedElement key={item.timestamp.toString()}>
              <HistoryCard
                activeTS={activeTS}
                key={item.timestamp.toString()}
                item={item}
                dispatch={dispatch}
              />
            </AnimatedElement>
          ))
        ) : (
          <Center>
            <Text color="dimmed" size="sm">
              History is a mystery.
            </Text>
          </Center>
        )}
      </AnimatePresence>
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
            radius="sm"
            variant="outline"
            fullWidth
            leftIcon={<MagnifyingGlass size={18} weight="bold" />}
            rightIcon={
              loading ? (
                <Loader color="brand" size="xs" />
              ) : (
                <Badge color="brand" radius="sm" size="sm" variant="light">
                  {`${item.results?.length} RESULTS`}
                </Badge>
              )
            }>
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
