import { StatusIndicator, Table, Tooltip } from "evergreen-ui";
import React from "react";
import { Highlight } from "react-instantsearch-hooks-web";
import { DownloadButton } from "../DownloadButton";
import type { BookHit } from "../../pages/LiveSearchPage";

const onlineServers = ["test"];

export function GridRow({ book }: { book: BookHit }) {
  return (
    <Table.Row isSelectable height={50}>
      <Table.TextCell flexBasis={120} flexGrow={0} flexShrink={0}>
        <Tooltip
          content={onlineServers?.includes(book.server) ? "Online" : "Offline"}>
          <StatusIndicator
            color={
              onlineServers?.includes(book.server) ? "success" : undefined
            }></StatusIndicator>
        </Tooltip>
        <Highlight hit={book} attribute="server" />
      </Table.TextCell>
      <Table.TextCell flexBasis={250} flexGrow={0} flexShrink={0}>
        <Highlight hit={book} attribute="author" />
      </Table.TextCell>
      <Table.TextCell>
        <Highlight hit={book} attribute="title" />
      </Table.TextCell>
      <Table.TextCell flexBasis={100} flexGrow={0} flexShrink={0}>
        <Highlight hit={book} attribute="format" />
      </Table.TextCell>
      <Table.TextCell flexBasis={100} flexGrow={0} flexShrink={0}>
        <Highlight hit={book} attribute="size" />
      </Table.TextCell>
      <Table.Cell flexBasis={150} flexGrow={0} flexShrink={0}>
        <DownloadButton book={book.full} />
      </Table.Cell>
    </Table.Row>
  );
}
