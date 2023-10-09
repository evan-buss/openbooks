import { createSelector, createSlice, PayloadAction } from "@reduxjs/toolkit";
import { RootState } from "./store";

export interface IrcServer {
  name: string;
  address: string;
  channel: string;
  port: number;
  sslPort: number;
}

export interface SelectedServer {
  name: string;
  enableTLS: boolean;
}

export interface ConnectionState {
  servers: IrcServer[];
  selectedServer: SelectedServer;
}

const loadConnection = (): SelectedServer => {
  const defaultConnection = { name: "IRC Highway", enableTLS: true };
  try {
    return JSON.parse(localStorage.getItem("connection")!) ?? defaultConnection;
  } catch (err) {
    return defaultConnection;
  }
};

const initialState = (): ConnectionState => {
  const useMockServers = import.meta.env.VITE_MOCK_IRC_SERVERS;
  if (useMockServers) {
    console.info("Using MOCK IRC servers");
  } else {
    console.warn("Using LIVE IRC servers");
  }
  return {
    selectedServer: loadConnection(),
    servers: [
      {
        name: "IRC Highway",
        address: useMockServers ? "localhost" : "irc.irchighway.net",
        port: 6667,
        sslPort: 6697,
        channel: "ebooks"
      },
      {
        name: "Undernet",
        address: useMockServers ? "localhost" : "irc.undernet.org",
        port: 6667,
        sslPort: 6697,
        channel: "bookz"
      }
    ]
  };
};

const connectionSlice = createSlice({
  name: "connection",
  initialState: initialState(),
  reducers: {
    setIrcServer(state, action: PayloadAction<SelectedServer>) {
      state.selectedServer = action.payload;
    }
  }
});

const selectConnection = (state: RootState) => state.connection;

const selectActiveIrcServer = createSelector(
  selectConnection,
  ({ servers, selectedServer }) => {
    const server = servers.find((x) => x.name === selectedServer.name)!;
    return { ...server, enableTLS: selectedServer.enableTLS };
  }
);

export const { setIrcServer } = connectionSlice.actions;
export { connectionSlice, selectActiveIrcServer };
export default connectionSlice.reducer;
