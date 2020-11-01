import { toaster } from "evergreen-ui";

export enum NotificationType {
    NOTIFY,
    SUCCESS,
    WARNING,
    DANGER
}

export const displayNotification = (type: NotificationType, message: string) => {
    switch (type) {
        case NotificationType.NOTIFY:
            toaster.notify(message);
            break;
        case NotificationType.SUCCESS:
            toaster.success(message);
            break;
        case NotificationType.WARNING:
            toaster.warning(message);
            break;
        case NotificationType.DANGER:
            toaster.danger(message);
            break;
    }
}
