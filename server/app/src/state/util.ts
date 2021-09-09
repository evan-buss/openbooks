
export const getWebsocketURL = (): URL => {
  const websocketURL = new URL(window.location.href + "ws")
  if (websocketURL.protocol.startsWith("https")) {
    websocketURL.protocol = websocketURL.protocol.replace("https", "wss");
  } else {
    websocketURL.protocol = websocketURL.protocol.replace("http", "ws");
  }

  if (import.meta.env.DEV) {
    websocketURL.port = "5228";
  }

  return websocketURL;
}

export const getApiURL = (): URL => {
  const apiURL = new URL(window.location.href)
  if (import.meta.env.DEV) {
    apiURL.port = "5228";
  }

  return apiURL;
}
