import {
  Button,
  majorScale,
  Pane,
  Spinner,
  StatusIndicator,
  Table,
  Text,
  Tooltip
} from "evergreen-ui";
import React, { useEffect, useMemo, useState } from "react";
import { useDispatch } from "react-redux";
import { BookDetail, MessageType } from "../../models/messages";
import { sendMessage } from "../../state/stateSlice";
import SelectMenuHeader, { makeStatusMenuItem } from "./SelectMenuHeader";
import { useGetServersQuery } from "../../state/api";

const stringContains = (first: string, second: string): boolean => {
  return first.toLowerCase().includes(second.toLowerCase());
};

export interface Props {
  books?: BookDetail[];
}

export const BooksGrid: React.FC<Props> = ({ books }: Props) => {
  const dispatch = useDispatch();
  const { data: onlineServers, refetch } = useGetServersQuery(null);
  const [titleFilter, setTitleFilter] = useState<string>("");
  const [authorFilter, setAuthorFilter] = useState<string>("");
  const [formatFilter, setFormatFilter] = useState<string[]>([]);
  const [serverFilter, setServerFilter] = useState<string[]>([]);

  useEffect(() => {
    refetch();
  }, [books]);

  const filteredBooks = useMemo(() => {
    return (books ?? [])
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
  }, [books, authorFilter, titleFilter, formatFilter, serverFilter]);

  const availableExtensions = useMemo(() => {
    const extensionTypes = books?.map((x) => x.format) ?? [];
    return extensionTypes.filter((x, i, arr) => arr.indexOf(x) === i);
  }, [books]);

  const availableServers = useMemo(() => {
    const servers =
      onlineServers?.concat(books?.map((x) => x.server) ?? []) ?? [];
    return servers.filter((x, i, arr) => arr.indexOf(x) === i);
  }, [onlineServers, books]);

  if (!books) {
    return (
      <Pane>
        <Spinner marginX="auto" marginY={120} />
      </Pane>
    );
  }

  const renderBody = () => {
    // Results length is zero. No results for query
    if (books?.length === 0) {
      return (
        <Text textAlign="center" padding={majorScale(3)}>
          No results for search.
        </Text>
      );
    }

    // FilteredRows length is zero. Filter not valid.
    if (filteredBooks.length === 0) {
      return (
        <Text textAlign="center" padding={majorScale(3)}>
          No books matching filter.
        </Text>
      );
    }

    return filteredBooks.map((book, i) => (
      <Table.Row key={book.full + i} isSelectable>
        <Table.TextCell flexBasis={120} flexGrow={0} flexShrink={0}>
          <Tooltip
            content={
              onlineServers?.includes(book.server) ? "Online" : "Offline"
            }>
            <StatusIndicator
              color={
                onlineServers?.includes(book.server) ? "success" : undefined
              }></StatusIndicator>
          </Tooltip>
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
    <Table
      flex={1}
      display="flex"
      flexDirection="column"
      width="100%"
      className="rounded-lg">
      <Table.Head background="white" className="rounded-t-md">
        <SelectMenuHeader
          flexBasis={120}
          options={availableServers}
          columnTitle="Servers"
          menuTitle="Filter Servers"
          selected={serverFilter}
          setSelected={setServerFilter}
          itemRenderer={makeStatusMenuItem(onlineServers ?? [])}
        />
        <Table.SearchHeaderCell
          onChange={(e) => setAuthorFilter(e)}
          placeholder="AUTHOR"
          flexBasis={250}
          flexGrow={0}
          flexShrink={0}
        />
        <Table.SearchHeaderCell
          onChange={(e) => setTitleFilter(e)}
          placeholder="TITLE"
        />
        <SelectMenuHeader
          columnTitle="Formats"
          menuTitle="Select Book Formats"
          options={availableExtensions}
          selected={formatFilter}
          setSelected={setFormatFilter}
        />
        <Table.TextHeaderCell
          flexBasis={100}
          flexGrow={0}
          flexShrink={0}
          className="font-normal text-gray-700">
          Size
        </Table.TextHeaderCell>
        <Table.TextHeaderCell
          flexBasis={150}
          flexGrow={0}
          flexShrink={0}
          className="font-normal text-gray-700">
          Download?
        </Table.TextHeaderCell>
      </Table.Head>
      <Table.VirtualBody>{renderBody()}</Table.VirtualBody>
    </Table>
  );
};
