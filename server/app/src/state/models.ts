export enum NotificationType {
  NOTIFY,
  SUCCESS,
  WARNING,
  DANGER
}

export interface Notification {
  type: NotificationType;
  title: string;
  detail?: string;
  timestamp: number;
}
