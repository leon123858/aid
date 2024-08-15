// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

import { FluentProvider, webLightTheme } from "@fluentui/react-components";
import Chat from "./Chat.tsx";
import Auth from "./Auth.tsx";
import styles from "./App.module.css";
import { UserProvider } from "./Provider.tsx";

function App() {
  return (
    <FluentProvider theme={webLightTheme}>
      <UserProvider>
        <div className={styles.appContainer}>
          <Chat style={{ flex: 1 }} />
          <Auth side={"right"} />
        </div>
      </UserProvider>
    </FluentProvider>
  );
}

export default App;
