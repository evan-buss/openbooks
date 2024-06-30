import { useAutoAnimate } from "@formkit/auto-animate/react";
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
import { Eye, EyeSlash, MagnifyingGlass, Trash } from "@phosphor-icons/react";
import {
  HistoryItem,
  removeResults,
  setActiveItem
} from "../../state/historySlice";
import { AppDispatch, useAppDispatch, useAppSelector } from "../../state/store";
import { conditionalAttribute } from "../../utils/attribute-helper";
import classes from "./SidebarButton.module.css";

export default function History() {
  const history = useAppSelector((store) => store.history.items);
  const activeTS = useAppSelector((store) => store.history.active);
  const dispatch = useAppDispatch();
  const [parent] = useAutoAnimate(/* optional config */);

  return (
    <Stack gap="xs" ref={parent}>
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
        <Center>
          <Text c="dimmed" size="sm">
            History is a mystery.
          </Text>
        </Center>
      )}
    </Stack>
  );
}

type Props = {
  activeTS: number | undefined;
  item: HistoryItem;
  dispatch: AppDispatch;
};

function HistoryCard({ activeTS, item, dispatch }: Props) {
  const isActive = activeTS === item.timestamp;
  const loading = !item.results?.length && !item.errors?.length;

  return (
    <Menu shadow="md">
      <Menu.Target>
        <Tooltip label={item.query} openDelay={1_000}>
          <Button
            classNames={{
              root: classes.root,
              inner: classes.inner,
              label: classes.label
            }}
            {...conditionalAttribute("active", isActive)}
            radius="sm"
            variant="outline"
            fullWidth
            leftSection={<MagnifyingGlass size={18} weight="bold" />}
            rightSection={
              loading ? (
                <Loader color="blue" size="xs" />
              ) : (
                <Badge color="blue" radius="sm" size="sm" variant="light">
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
            leftSection={<Eye size={18} weight="bold" />}
            onClick={() => dispatch(setActiveItem(item))}>
            Show Results
          </Menu.Item>
        ) : (
          <Menu.Item
            leftSection={<EyeSlash size={18} weight="bold" />}
            onClick={() => dispatch(setActiveItem(null))}>
            Hide Results
          </Menu.Item>
        )}

        <Menu.Item
          color="red"
          leftSection={<Trash size={18} weight="bold" />}
          onClick={() => dispatch(removeResults(item.timestamp))}>
          Delete item
        </Menu.Item>
      </Menu.Dropdown>
    </Menu>
  );
}
