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

export const api = createApi({
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
    getBooks: builder.query<Book[], void>({
      query: () => `library`,
      providesTags: ["books"]
    }),
    deleteBook: builder.mutation<void, string>({
      query: (book) => ({
        url: `library/${book}`,
        method: "DELETE"
      }),
      invalidatesTags: ["books"]
    }),
    search: builder.mutation<void, string>({
      query: (searchQuery) => ({
        url: `search`,
        method: "POST",
        params: { query: searchQuery }
      })
    }),
    download: builder.mutation<void, string>({
      query: (book) => ({
        url: `download`,
        params: { book },
        method: "POST"
      })
    })
  })
});

export const {
  useGetServersQuery,
  useGetBooksQuery,
  useDeleteBookMutation,
  useSearchMutation,
  useDownloadMutation
} = api;
