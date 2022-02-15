import { Button, majorScale, Pane, SearchInput, Text } from "evergreen-ui";
import React, { FormEvent, useEffect, useState } from "react";
import { useDispatch, useSelector } from "react-redux";
import { BooksGrid } from "../components/BooksGrid/BooksGrid";
import { sendMessage, sendSearch } from "../state/stateSlice";
import { RootState } from "../state/store";
import ErrorsGrid from "../components/ErrorsGrid/ErrorsGrid";
import { MessageType } from "../state/messages";
import { Warning } from "phosphor-react";
import image from "../assets/reading.svg";

export default function SearchPage() {
  const dispatch = useDispatch();
  const activeItem = useSelector((store: RootState) => store.state.activeItem);
  const [searchQuery, setSearchQuery] = useState("");
  const [showErrors, setShowErrors] = useState(false);

  useEffect(() => {
    setShowErrors(false);
  }, [activeItem]);

  const hasErrors = () =>
    activeItem?.errors?.length && activeItem?.errors?.length > 0;
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
      return <ErrorsGrid errors={activeItem?.errors} />;
    } else {
      return <BooksGrid books={activeItem?.results} />;
    }
  };

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
        <SearchInput
          disabled={activeItem !== null && !activeItem.results}
          value={searchQuery}
          onChange={(e: any) => setSearchQuery(e.target.value)}
          placeholder={
            errorMode() ? "Download a book manually." : "Search for a book."
          }
          height={majorScale(5)}
          className="rounded-md"
          marginRight={majorScale(4)}
          width="100%"
        />
        <Button
          type="submit"
          className="rounded-md"
          width={majorScale(16)}
          height={majorScale(5)}
          appearance="primary"
          disabled={!validInput()}>
          {errorMode() ? "Download" : "Search"}
        </Button>
      </form>
      {hasErrors() && (
        <Button
          appearance={errorMode() ? "primary" : "minimal"}
          onClick={() => setShowErrors((show) => !show)}
          className="mb-2 self-start"
          iconBefore={<Warning size={18} />}
          marginRight={12}
          size="small">
          {activeItem?.errors?.length} Parsing{" "}
          {activeItem?.errors?.length === 1 ? "Error" : "Errors"}
        </Button>
      )}
      {renderBody()}
    </Pane>
  );
}
