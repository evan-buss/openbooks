import { createAsyncThunk, createSlice, PayloadAction } from "@reduxjs/toolkit";
import { RootState } from "./store";
import { sendMessage } from "./stateSlice";
import { MessageType } from "./messages";

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
  console.log("Using mock servers? ", useMockServers);
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
  initialState,
  reducers: {
    selectServer(state, action: PayloadAction<SelectedServer>) {
      state.selectedServer = action.payload;
    }
  }
});

const activeIrcServer = ({
  connection: { servers, selectedServer }
}: RootState): IrcServer & { enableTLS: boolean } => {
  const server = servers.find((x) => x.name === selectedServer.name)!;
  return { ...server, enableTLS: selectedServer.enableTLS };
};

const connectToServer = createAsyncThunk(
  "connection/connect_to_server",
  async (_, { dispatch, getState }) => {
    const { address, port, channel, enableTLS } = activeIrcServer(
      getState() as RootState
    );
    dispatch(
      sendMessage({
        type: MessageType.CONNECT,
        payload: { address: `${address}:${port}`, channel, enableTLS }
      })
    );
  }
);

const { selectServer } = connectionSlice.actions;

export { connectionSlice, selectServer, connectToServer, activeIrcServer };

export default connectionSlice.reducer;
