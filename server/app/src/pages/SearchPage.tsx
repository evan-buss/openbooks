import {
  ActionIcon,
  Button,
  Center,
  Group,
  Image,
  Stack,
  TextInput,
  Title,
  useMantineColorScheme
} from "@mantine/core";
import { MagnifyingGlass, Sidebar, Warning } from "@phosphor-icons/react";
import { FormEvent, useEffect, useMemo, useState } from "react";
import image from "../assets/reading.svg";
import BookTable from "../components/tables/BookTable";
import ErrorTable from "../components/tables/ErrorTable";
import { MessageType } from "../state/messages";
import { sendMessage, sendSearch, toggleSidebar } from "../state/stateSlice";
import { useAppDispatch, useAppSelector } from "../state/store";
import classes from "./SearchPage.module.css";
import { conditionalAttribute } from "../utils/attribute-helper";

export default function SearchPage() {
  const dispatch = useAppDispatch();
  const activeItem = useAppSelector((store) => store.state.activeItem);
  const opened = useAppSelector((store) => store.state.isSidebarOpen);

  const { colorScheme } = useMantineColorScheme();

  const [searchQuery, setSearchQuery] = useState("");
  const [showErrors, setShowErrors] = useState(false);

  const hasErrors = (activeItem?.errors ?? []).length > 0;
  const errorMode = showErrors && !!activeItem;
  const validInput = errorMode
    ? searchQuery.startsWith("!")
    : searchQuery !== "";

  useEffect(() => {
    setShowErrors(false);
  }, [activeItem]);

  const searchHandler = (event: FormEvent) => {
    event.preventDefault();

    if (errorMode) {
      dispatch(
        sendMessage({
          type: MessageType.DOWNLOAD,
          payload: { book: searchQuery }
        })
      );
    } else {
      dispatch(sendSearch(searchQuery));
    }

    setSearchQuery("");
  };

  const bookTable = useMemo(
    () => <BookTable books={activeItem?.results ?? []} />,
    [activeItem?.results]
  );

  const errorTable = useMemo(
    () => (
      <ErrorTable
        errors={activeItem?.errors ?? []}
        setSearchQuery={setSearchQuery}
      />
    ),
    [activeItem?.errors]
  );

  return (
    <Stack gap={0} align="flex-start" style={{ flexGrow: 1 }}>
      <form style={{ width: "100%" }} onSubmit={(e) => searchHandler(e)}>
        <Group wrap="nowrap" gap="md" mb="md">
          {!opened && (
            <ActionIcon size="lg" onClick={() => dispatch(toggleSidebar())}>
              <Sidebar weight="bold" size={20}></Sidebar>
            </ActionIcon>
          )}
          <TextInput
            w="100%"
            variant="filled"
            disabled={activeItem !== null && !activeItem.results}
            value={searchQuery}
            onChange={(e: any) => setSearchQuery(e.target.value)}
            placeholder={
              errorMode ? "Download a book manually." : "Search for a book."
            }
            radius="md"
            type="search"
            leftSection={<MagnifyingGlass weight="bold" size={22} />}
            required
          />

          <Button
            type="submit"
            disabled={!validInput}
            color={colorScheme === "dark" ? "blue.2" : "blue"}
            radius="md"
            variant={validInput ? "gradient" : "default"}
            gradient={{ from: "blue.4", to: "blue.3" }}>
            {errorMode ? "Download" : "Search"}
          </Button>
        </Group>
      </form>

      {hasErrors && (
        <Button
          {...conditionalAttribute("error-mode", errorMode)}
          className={classes.errorToggle}
          variant={errorMode ? "filled" : "subtle"}
          onClick={() => setShowErrors((show) => !show)}
          leftSection={<Warning size={18} />}
          size="compact-xs">
          {activeItem?.errors?.length} Parsing{" "}
          {activeItem?.errors?.length === 1 ? "Error" : "Errors"}
        </Button>
      )}
      {!activeItem ? (
        <Center style={{ height: "100%", width: "100%" }}>
          <Stack align="center">
            <Title fw="normal" ta="center">
              Search a book to get started.
            </Title>
            <Image
              hiddenFrom="lg"
              w={300}
              fit="contain"
              src={image}
              alt="person reading"
            />
            <Image
              visibleFrom="lg"
              w={600}
              fit="contain"
              src={image}
              alt="person reading"
            />
          </Stack>
        </Center>
      ) : errorMode ? (
        errorTable
      ) : (
        bookTable
      )}
    </Stack>
  );
}
