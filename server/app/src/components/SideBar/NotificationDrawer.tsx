import { SideSheet, Alert, IntentTypes, IconButton, Text } from "evergreen-ui";
import React from "react";
import { useSelector } from "react-redux";
import {
  clearNotifications,
  dismissNotification,
  toggleDrawer
} from "../../state/notificationSlice";
import { RootState, useAppDispatch } from "../../state/store";
import { BellSimpleSlash } from "phosphor-react";
import { NotificationType } from "../../state/messages";

const NotificationDrawer = () => {
  const { isOpen, notifications } = useSelector(
    (store: RootState) => store.notifications
  );
  const dispatch = useAppDispatch();

  const getIntent = (type: NotificationType): IntentTypes => {
    switch (type) {
      case NotificationType.NOTIFY:
        return "none";
      case NotificationType.DANGER:
        return "danger";
      case NotificationType.SUCCESS:
        return "success";
      case NotificationType.WARNING:
        return "warning";
    }
  };

  return (
    <SideSheet
      isShown={isOpen}
      width={400}
      onCloseComplete={() => dispatch(toggleDrawer())}>
      <div className="flex flex-row justify-between p-3 border-b border-gray-300">
        <h2 className="text-xl">Notifications</h2>
        <IconButton
          disabled={notifications.length === 0}
          icon={
            <BellSimpleSlash
              color={notifications.length === 0 ? "gray" : "black"}
              weight="bold"
            />
          }
          onClick={() => dispatch(clearNotifications())}></IconButton>
      </div>
      <div className="flex flex-col gap-2 p-2">
        {notifications.length === 0 ? (
          <Text marginX="auto" marginY={16} color="muted">
            No notifications.
          </Text>
        ) : (
          notifications.map((notif) => (
            <div key={notif.timestamp}>
              <p className="mb-1 text-xs font-semibold text-gray-600">
                {new Date(notif.timestamp).toLocaleTimeString("en-US", {
                  timeStyle: "short"
                })}
              </p>
              <Alert
                intent={getIntent(notif.appearance)}
                title={notif.title}
                onRemove={() => dispatch(dismissNotification(notif))}
                isRemoveable>
                {notif.detail}
              </Alert>
            </div>
          ))
        )}
      </div>
    </SideSheet>
  );
};

export default NotificationDrawer;
