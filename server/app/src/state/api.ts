import { createApi, fetchBaseQuery } from '@reduxjs/toolkit/query/react'
import { getApiURL } from "./util";

export interface IrcServer {
  elevatedUsers?: string[];
  regularUsers?: string[];
}

export const openbooksApi = createApi({
  baseQuery: fetchBaseQuery({
    baseUrl: getApiURL().href
  }),
  endpoints: (builder) => ({
    getServers: builder.query<string[], null>({
      query: () => `servers`,
      transformResponse: (ircServers: IrcServer) => {
        return ircServers.elevatedUsers ?? [];
      }
    }),
  })
});

export const { useGetServersQuery } = openbooksApi;