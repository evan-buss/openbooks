import {
  Button as MButton,
  TextInput,
  useMantineColorScheme
} from "@mantine/core";
import { Button, majorScale, Pane, Text } from "evergreen-ui";
import { MagnifyingGlass, Warning } from "phosphor-react";
import { FormEvent, useEffect, useState } from "react";
import { useSelector } from "react-redux";
import image from "../assets/reading.svg";
import { BooksGrid } from "../components/BooksGrid/BooksGrid";
import ErrorsGrid from "../components/ErrorsGrid/ErrorsGrid";
import { MessageType } from "../state/messages";
import { sendMessage, sendSearch } from "../state/stateSlice";
import { RootState, useAppDispatch } from "../state/store";

export default function SearchPage() {
  const dispatch = useAppDispatch();
  const activeItem = useSelector((store: RootState) => store.state.activeItem);
  const [searchQuery, setSearchQuery] = useState("");
  const [showErrors, setShowErrors] = useState(false);

  useEffect(() => {
    setShowErrors(false);
  }, [activeItem]);

  const hasErrors = () => (activeItem?.errors ?? []).length > 0;
  const errorMode = () => showErrors && activeItem;
  const validInput = () => {
    if (errorMode()) {
      return searchQuery.startsWith("!");
    } else {
      return searchQuery !== "";
    }
  };

  const searchHandler = (event: FormEvent) => {
    event.preventDefault();

    if (errorMode()) {
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

  const renderBody = () => {
    if (activeItem === null) {
      return (
        <>
          <Text fontSize="2em" margin="3em">
            Search a book to get started.
          </Text>
          <img className="w-96 xl:w-1/3" src={image} alt="person reading" />
        </>
      );
    } else if (errorMode()) {
      return (
        <ErrorsGrid
          errors={activeItem?.errors}
          setSearchQuery={setSearchQuery}
        />
      );
    } else {
      return <BooksGrid books={activeItem?.results} />;
    }
  };

  const { colorScheme } = useMantineColorScheme();

  return (
    <Pane
      margin={majorScale(4)}
      width="100%"
      display="flex"
      flexDirection="column"
      alignItems="center">
      <form
        style={{
          width: "100%",
          display: "flex",
          flexFlow: "row nowrap",
          justifyContent: "center",
          alignItems: "center",
          marginBottom: majorScale(4)
        }}
        onSubmit={(e) => searchHandler(e)}>
        <TextInput
          sx={(theme) => ({ width: "100%", marginRight: theme.spacing.xl })}
          disabled={activeItem !== null && !activeItem.results}
          value={searchQuery}
          onChange={(e: any) => setSearchQuery(e.target.value)}
          placeholder={
            errorMode() ? "Download a book manually." : "Search for a book."
          }
          radius="md"
          size="md"
          type="search"
          icon={<MagnifyingGlass weight="bold" size={22} />}
          required
        />

        <MButton
          type="submit"
          color={colorScheme === "dark" ? "brand.2" : "brand"}
          disabled={!validInput()}
          size="md"
          radius="md"
          variant="outline">
          {errorMode() ? "Download" : "Search"}
        </MButton>
      </form>
      {hasErrors() && (
        <Button
          appearance={errorMode() ? "primary" : "minimal"}
          onClick={() => setShowErrors((show) => !show)}
          className="self-start mb-2"
          iconBefore={<Warning size={18} />}
          marginRight={12}
          size="small">
          {activeItem?.errors?.length} Parsing{" "}
          {activeItem?.errors?.length === 1 ? "Error" : "Errors"}
        </Button>

        // <MButton
        //   sx={(theme) => ({
        //     alignSelf: "start",
        //     marginBottom: theme.spacing.xs,
        //     fontWeight: "normal"
        //   })}
        //   color="dark"
        //   variant={errorMode() ? "filled" : "subtle"}
        //   onClick={() => setShowErrors((show) => !show)}
        //   leftIcon={<Warning size={18} />}
        //   size="xs">
        //   {activeItem?.errors?.length} Parsing{" "}
        //   {activeItem?.errors?.length === 1 ? "Error" : "Errors"}
        // </MButton>
      )}
      {renderBody()}
    </Pane>
  );
}
