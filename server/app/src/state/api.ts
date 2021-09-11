import { ResultDescription } from "@reduxjs/toolkit/dist/query/endpointDefinitions";
import { createApi, fetchBaseQuery } from "@reduxjs/toolkit/query/react";
import { getApiURL } from "./util";

export interface IrcServer {
  elevatedUsers?: string[];
  regularUsers?: string[];
}

export interface Book {
  name: string;
  downloadLink: string;
  time: string;
}

export const openbooksApi = createApi({
  baseQuery: fetchBaseQuery({
    baseUrl: getApiURL().href,
    credentials: "include",
    mode: "cors"
  }),
  tagTypes: ["books", "servers"],
  endpoints: (builder) => ({
    getServers: builder.query<string[], null>({
      query: () => `servers`,
      transformResponse: (ircServers: IrcServer) => {
        return ircServers.elevatedUsers ?? [];
      }
    }),
    getBooks: builder.query<Book[], null>({
      query: () => `library`,
      providesTags: ["books"]
    }),
    deleteBook: builder.mutation<null, string>({
      query: (book) => ({
        url: `library/${book}`,
        method: "DELETE"
      }),
      invalidatesTags: ["books"]
    })
  })
});

export const { useGetServersQuery, useGetBooksQuery, useDeleteBookMutation } =
  openbooksApi;
