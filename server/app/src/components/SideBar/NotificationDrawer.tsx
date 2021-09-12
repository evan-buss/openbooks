import { SideSheet, Alert, IntentTypes } from "evergreen-ui";
import React from "react";
import { useDispatch, useSelector } from "react-redux";
import {
  dismissNotification,
  toggleDrawer
} from "../../state/notificationSlice";
import { RootState } from "../../state/store";
import { NotificationType } from "../../state/models";

const NotificationDrawer = () => {
  const { isOpen, notifications } = useSelector(
    (store: RootState) => store.notifications
  );
  const dispatch = useDispatch();

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
      <div className="p-2 border-b border-gray-300">
        <h2 className="text-2xl font-medium">Notifications</h2>
      </div>
      <div className="flex flex-col p-2 gap-2">
        {notifications.map((notif) => (
          <div key={notif.timestamp}>
            <p className="text-gray-600 font-semibold text-xs mb-1">
              {new Date(notif.timestamp).toLocaleTimeString("en-US", {
                timeStyle: "short"
              })}
            </p>
            <Alert
              intent={getIntent(notif.type)}
              title={notif.title}
              onRemove={() => dispatch(dismissNotification(notif))}
              isRemoveable>
              {notif.detail}
            </Alert>
          </div>
        ))}
      </div>
    </SideSheet>
  );
};

export default NotificationDrawer;
