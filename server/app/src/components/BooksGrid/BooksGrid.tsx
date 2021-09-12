import { Button, majorScale, Pane, Spinner, Table, Text } from "evergreen-ui";
import React, { useEffect, useMemo, useState } from "react";
import { useDispatch, useSelector } from "react-redux";
import { MessageType } from "../../models/messages";
import { sendMessage } from "../../state/stateSlice";
import { RootState } from "../../state/store";
import image from "../../assets/reading.svg";
import SelectMenuHeader from "./SelectMenuHeader";
import { useGetServersQuery } from "../../state/api";

const stringContains = (first: string, second: string): boolean => {
  return first.toLowerCase().includes(second.toLowerCase());
};

export const BooksGrid: React.FC = () => {
  const dispatch = useDispatch();
  const activeItem = useSelector((store: RootState) => store.state.activeItem);
  const { data: allServers, refetch } = useGetServersQuery(null);
  const [titleFilter, setTitleFilter] = useState<string>("");
  const [authorFilter, setAuthorFilter] = useState<string>("");
  const [formatFilter, setFormatFilter] = useState<string[]>([]);
  const [serverFilter, setServerFilter] = useState<string[]>([]);

  useEffect(() => {
    refetch();
  }, [activeItem]);

  const filteredBooks = useMemo(() => {
    const books = activeItem?.results ?? [];
    return books
      .filter(
        (x) =>
          serverFilter.length === 0 ||
          serverFilter.findIndex((server) => server === x.server) !== -1
      )
      .filter((x) => stringContains(x.author, authorFilter))
      .filter((x) => stringContains(x.title, titleFilter))
      .filter(
        (x) =>
          formatFilter.length === 0 ||
          formatFilter.findIndex((format) => format === x.format) !== -1
      );
  }, [activeItem, authorFilter, titleFilter, formatFilter, serverFilter]);

  const availableExtensions = useMemo(() => {
    const extensionTypes = activeItem?.results?.map((x) => x.format) ?? [];
    return extensionTypes.filter((x, i, arr) => arr.indexOf(x) === i);
  }, [activeItem]);

  if (activeItem === null) {
    return (
      <>
        <Text fontSize="2em" margin="3em">
          Search a book to get started.
        </Text>
        <img className="w-96 xl:w-1/3" src={image} alt="person reading" />
      </>
    );
  }

  if (activeItem !== null && !activeItem.results) {
    return (
      <Pane>
        <Spinner marginX="auto" marginY={120} />
      </Pane>
    );
  }

  const renderBody = () => {
    // Results length is zero. No results for query
    if (activeItem.results?.length === 0) {
      return (
        <Text textAlign="center" margin={majorScale(3)}>
          No results for search.
        </Text>
      );
    }

    // FilteredRows length is zero. Filter not valid.
    if (filteredBooks.length === 0) {
      return (
        <Text textAlign="center" margin={majorScale(3)}>
          No books matching filter.
        </Text>
      );
    }

    return filteredBooks.map((book, i) => (
      <Table.Row key={book.full + i} isSelectable>
        <Table.TextCell flexBasis={100} flexGrow={0} flexShrink={0}>
          {book.server}
        </Table.TextCell>
        <Table.TextCell flexBasis={250} flexGrow={0} flexShrink={0}>
          {book.author}
        </Table.TextCell>
        <Table.TextCell>{book.title}</Table.TextCell>
        <Table.TextCell flexBasis={100} flexGrow={0} flexShrink={0}>
          {book.format}
        </Table.TextCell>
        <Table.TextCell flexBasis={100} flexGrow={0} flexShrink={0}>
          {book.size}
        </Table.TextCell>
        <Table.Cell flexBasis={150} flexGrow={0} flexShrink={0}>
          <Button
            appearance="primary"
            size="small"
            onClick={() =>
              dispatch(
                sendMessage({
                  type: MessageType.DOWNLOAD,
                  payload: { book: book.full }
                })
              )
            }>
            Download
          </Button>
        </Table.Cell>
      </Table.Row>
    ));
  };

  return (
    <Table flex={1} display="flex" flexDirection="column" width="100%">
      <Table.Head background="white">
        <SelectMenuHeader
          options={allServers ?? []}
          columnTitle="Servers"
          menuTitle="Filter Online Servers"
          selected={serverFilter}
          setSelected={setServerFilter}
        />
        <Table.SearchHeaderCell
          onChange={(e) => setAuthorFilter(e)}
          placeholder="AUTHOR"
          flexBasis={250}
          flexGrow={0}
          flexShrink={0}></Table.SearchHeaderCell>
        <Table.SearchHeaderCell
          onChange={(e) => setTitleFilter(e)}
          placeholder="TITLE"></Table.SearchHeaderCell>
        <SelectMenuHeader
          columnTitle="Formats"
          menuTitle="Select Book Formats"
          options={availableExtensions}
          selected={formatFilter}
          setSelected={setFormatFilter}
        />
        <Table.TextHeaderCell flexBasis={100} flexGrow={0} flexShrink={0}>
          Size
        </Table.TextHeaderCell>
        <Table.TextHeaderCell flexBasis={150} flexGrow={0} flexShrink={0}>
          Download?
        </Table.TextHeaderCell>
      </Table.Head>
      <Table.VirtualBody>{renderBody()}</Table.VirtualBody>
    </Table>
  );
};
