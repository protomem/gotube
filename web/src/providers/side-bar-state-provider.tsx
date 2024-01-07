import {
  ReactNode,
  createContext,
  useContext,
  useEffect,
  useState,
} from "react";

type State = {
  onOpen: () => void;
  onClose: () => void;
  onToggle: () => void;
  isOpen: boolean;
  hide: boolean;
};

const SideBarStateContext = createContext<State>({
  onOpen: () => {},
  onClose: () => {},
  onToggle: () => {},
  isOpen: false,
  hide: false,
});

export const useSideBarState = () => useContext(SideBarStateContext);

type Props = {
  children: ReactNode;
};

const SideBarStateProvider = ({ children }: Props) => {
  const defaultIsOpen = () => {
    const state = localStorage.getItem("stateSideBar");
    if (state) {
      const { isOpen } = JSON.parse(state);
      return isOpen;
    }
    return false;
  };

  const [isOpen, setIsOpen] = useState(defaultIsOpen);
  const onOpen = () => setIsOpen(true);
  const onClose = () => setIsOpen(false);
  const onToggle = () => setIsOpen(!isOpen);

  useEffect(() => {
    localStorage.setItem("stateSideBar", JSON.stringify({ isOpen }));
  }, [isOpen]);

  return (
    <SideBarStateContext.Provider
      value={{ onOpen, onClose, isOpen, hide: false, onToggle }}
    >
      {children}
    </SideBarStateContext.Provider>
  );
};

export default SideBarStateProvider;
