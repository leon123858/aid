import { createContext, ReactNode, useContext, useState } from "react";

// 定義 context 的類型
interface UserContextType {
  uid: string | null;
  setUid: (uid: string | null) => void;
}

// 創建 context
const UserContext = createContext<UserContextType | undefined>(undefined);

// 創建 Provider 組件
export function UserProvider({ children }: { children: ReactNode }) {
  const [uid, setUid] = useState<string | null>("");

  return (
    <UserContext.Provider value={{ uid, setUid }}>
      {children}
    </UserContext.Provider>
  );
}

// 創建自定義 hook 來使用這個 context
// eslint-disable-next-line react-refresh/only-export-components
export function useUser() {
  const context = useContext(UserContext);
  if (context === undefined) {
    throw new Error("useUser must be used within a UserProvider");
  }
  return context;
}
