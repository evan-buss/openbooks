import { EmptyState, SearchIcon, Table } from "evergreen-ui";
import React from "react";
import { useHits } from "react-instantsearch-hooks-web";
import { BookHit } from "../../pages/LiveSearchPage";
import { useHeight } from "../../state/hooks";
import { GridRow } from "./GridRow";

export function Grid() {
  const { hits } = useHits<BookHit>();
  const height = useHeight();

  const emptyState = () => {
    return (
      <EmptyState
        background="light"
        title="No books found for the search term"
        orientation="horizontal"
        icon={<SearchIcon color="#C1C4D6" />}
        iconBgColor="#EDEFF5"
        description="Keep in mind that not all servers have been indexed yet. If you can't find what you are looking for try searching the IRC channel directly."
      />
    );
  };

  return (
    <Table
      key={height}
      flex={1}
      display="flex"
      flexDirection="column"
      width="100%"
      className="rounded-lg">
      <Table.Head background="white" className="rounded-t-md">
        <Table.TextHeaderCell
          flexBasis={120}
          flexGrow={0}
          flexShrink={0}
          className="font-normal text-gray-700">
          Server
        </Table.TextHeaderCell>
        <Table.TextHeaderCell
          flexBasis={250}
          flexGrow={0}
          flexShrink={0}
          className="font-normal text-gray-700">
          Author
        </Table.TextHeaderCell>
        <Table.TextHeaderCell className="font-normal text-gray-700">
          Title
        </Table.TextHeaderCell>
        <Table.TextHeaderCell
          flexBasis={100}
          flexGrow={0}
          flexShrink={0}
          className="font-normal text-gray-700">
          Format
        </Table.TextHeaderCell>
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
      <Table.Body>
        {hits.length > 0
          ? hits.map((hit) => <GridRow key={hit.id} book={hit} />)
          : emptyState()}
      </Table.Body>
    </Table>
  );
}
