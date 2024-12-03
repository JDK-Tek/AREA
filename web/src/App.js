/*
** EPITECH PROJECT, 2024
** AREA
** File description:
** App
*/

import AppletData from "./data/AppletData";
import AppletKit from "./components/AppletKit";

function App() {

  return (
    <div className="App">
      <AppletKit title={"Get started with any Applet"} applets={AppletData}/>
    </div>
  );
}

export default App;
