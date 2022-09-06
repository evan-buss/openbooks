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
import { AnimatePresence, motion } from "framer-motion";
import { Eye, EyeSlash, MagnifyingGlass, Trash } from "phosphor-react";
import { useSelector } from "react-redux";
import {
  deleteHistoryItem,
  HistoryItem,
  selectHistory
} from "../../state/historySlice";
import { setActiveItem } from "../../state/stateSlice";
import { useAppDispatch, useAppSelector } from "../../state/store";
import { defaultAnimation } from "../../utils/animation";
import { useSidebarButtonStyle } from "./styles";

export default function History() {
  const history = useSelector(selectHistory);
  const activeTS =
    useAppSelector((store) => store.state.activeItem?.timestamp) ?? -1;
  const dispatch = useAppDispatch();

  return (
    <Stack spacing="xs">
      <AnimatePresence mode="popLayout">
        {history.length > 0 ? (
          history.map((item: HistoryItem) => (
            <motion.div {...defaultAnimation} key={item.timestamp.toString()}>
              <HistoryCard
                activeTS={activeTS}
                key={item.timestamp.toString()}
                item={item}
                dispatch={dispatch}
              />
            </motion.div>
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

function HistoryCard({ activeTS, item, dispatch }: Props) {
  const isActive = activeTS === item.timestamp;
  const { classes } = useSidebarButtonStyle({ isActive });

  const loading = !item.results?.length && !item.errors?.length;

  return (
    <Menu shadow="md">
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
}
