import { Navigate, Outlet } from "react-router-dom";
import { useLocation } from "react-router-dom";
import { useAuth } from "../hooks";

export const RequireAuth = () => {
  const { accessToken } = useAuth();
  let location = useLocation();

  if (!accessToken) {
    return <Navigate to='/login' state={{ from: location }} />
  }

  return <Outlet />
}