import { ReactNode, createContext, useContext, useState } from "react";
import { User } from "../domain/entities";

type State = {
  isAuthenticated: boolean;
  currentUser: User | null;
  accessToken: string;
  refreshToken: string;
  login: (user: User, accessToken: string, refreshToken: string) => void;
  logout: () => void;
};

const AuthContext = createContext<State>({
  isAuthenticated: false,
  currentUser: null,
  accessToken: "",
  refreshToken: "",
  login: () => {},
  logout: () => {},
});

export const useAuth = () => useContext(AuthContext);

type Props = {
  children: ReactNode;
};

const AuthProvider = ({ children }: Props) => {
  const getCurrentUser = () => {
    const user = localStorage.getItem("user");
    if (user) {
      const res: User = JSON.parse(user);
      return res;
    }
    return null;
  };

  const getAccessToken = () => {
    const accessToken = localStorage.getItem("accessToken");
    if (accessToken) {
      return accessToken;
    }
    return "";
  };

  const getRefreshToken = () => {
    const refreshToken = localStorage.getItem("refreshToken");
    if (refreshToken) {
      return refreshToken;
    }
    return "";
  };

  const [state, setState] = useState<State>({
    isAuthenticated:
      !!getCurrentUser() && getAccessToken() !== "" && getRefreshToken() !== "",
    currentUser: getCurrentUser(),
    accessToken: getAccessToken(),
    refreshToken: getRefreshToken(),
    login: (user: User, accessToken: string, refreshToken: string) => {
      localStorage.setItem("user", JSON.stringify(user));
      localStorage.setItem("accessToken", accessToken);
      localStorage.setItem("refreshToken", refreshToken);
      setState({
        ...state,
        isAuthenticated: true,
        currentUser: user,
        accessToken: accessToken,
        refreshToken: refreshToken,
      });
    },
    logout: () => {
      localStorage.removeItem("user");
      localStorage.removeItem("accessToken");
      localStorage.removeItem("refreshToken");
      setState({
        ...state,
        isAuthenticated: false,
        currentUser: null,
        accessToken: "",
        refreshToken: "",
      });
    },
  });

  return <AuthContext.Provider value={state}>{children}</AuthContext.Provider>;
};

export default AuthProvider;
