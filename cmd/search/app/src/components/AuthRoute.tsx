import { LoadingOverlay } from "@mantine/core";
import { Navigate, useLocation } from "react-router-dom";
import useAuth from "../util/auth";

export function AuthRoute({ children }: { children: React.ReactNode }) {
  const { user, loading } = useAuth();
  const location = useLocation();

  if (loading) {
    return <LoadingOverlay visible={true} overlayBlur={2} />;
  }

  if (!!user) {
    return <>{children}</>;
  }

  return <Navigate to="/login" replace state={{ path: location.pathname }} />;
}
