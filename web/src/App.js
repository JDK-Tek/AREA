/*
** EPITECH PROJECT, 2024
** AREA
** File description:
** App
*/

import AppletData from "./data/AppletData";
import ServiceData from "./data/ServiceData";

import AppletKit from "./components/AppletKit";
import ServiceKit from "./components/ServiceKit";

function App() {

  return (
    <div className="App">
      
      <AppletKit
        title={"Get started with any Applet"}
        applets={AppletData}
      />

      <ServiceKit
        title={"or choose from 900+ services"}
        services={ServiceData}
        color={"text-chartpurple-200"}
      />
      
    </div>
  );
}

export default App;
