import { toaster } from "evergreen-ui";
import { Notification, NotificationType } from "./messages";

export const getWebsocketURL = (): URL => {
  const websocketURL = new URL(window.location.href + "ws");
  if (websocketURL.protocol.startsWith("https")) {
    websocketURL.protocol = websocketURL.protocol.replace("https", "wss");
  } else {
    websocketURL.protocol = websocketURL.protocol.replace("http", "ws");
  }

  if (import.meta.env.DEV) {
    websocketURL.port = "5228";
  }

  return websocketURL;
};

export const getApiURL = (): URL => {
  const apiURL = new URL(window.location.href);
  if (import.meta.env.DEV) {
    apiURL.port = "5228";
  }

  return apiURL;
};

export const displayNotification = ({
  appearance = NotificationType.NOTIFY,
  title,
  detail
}: Notification) => {
  switch (appearance) {
    case NotificationType.NOTIFY:
      toaster.notify(title, {
        description: detail
      });
      break;
    case NotificationType.SUCCESS:
      toaster.success(title, {
        description: detail
      });
      break;
    case NotificationType.WARNING:
      toaster.warning(title, {
        description: detail
      });
      break;
    case NotificationType.DANGER:
      toaster.danger(title, {
        description: detail
      });
      break;
  }
};

export const downloadFile = (relativeURL: string) => {
  let url = getApiURL();
  url.pathname += relativeURL;

  let link = document.createElement("a");
  link.download = "";
  link.target = "_blank";
  link.href = url.href;
  link.click();
  link.remove();
};
