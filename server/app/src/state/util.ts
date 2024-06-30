import { showNotification } from "@mantine/notifications";
import { Notification, NotificationType } from "./messages";

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
      showNotification({
        color: "blue",
        title: title,
        message: detail
      });
      break;
    case NotificationType.SUCCESS:
      showNotification({
        color: "green",
        title: title,
        message: detail
      });

      break;
    case NotificationType.WARNING:
      showNotification({
        color: "yellow",
        title: title,
        message: detail
      });
      break;
    case NotificationType.DANGER:
      showNotification({
        color: "red",
        title: title,
        message: detail
      });
      break;
  }
};

export function downloadFile(relativeURL?: string) {
  if (relativeURL === "" || relativeURL === undefined) return;

  const url = getApiURL();
  url.pathname += relativeURL;

  const link = document.createElement("a");
  link.download = "";
  link.target = "_blank";
  link.href = url.href;

  // Ensure the link isn't visible to the user or cause layout shifts.
  link.setAttribute("visibility", "hidden");
  document.body.appendChild(link);
  link.click();
  document.body.removeChild(link);
}
