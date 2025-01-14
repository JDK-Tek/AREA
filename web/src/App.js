/*
** EPITECH PROJECT, 2024
** AREA
** File description:
** App
*/

import React, { useState } from "react";

import { Routes, Route } from 'react-router-dom';

import listRoutes from "./data/Routes";

import Home from "./routes/home/Home";
import Login from "./routes/login/Login";
import Register from "./routes/register/Register";
import CreateArea from "./routes/create/CreateArea";
import Explore from "./routes/explore/Explore";
import NotFound from "./routes/notfound/NotFound";
// import AreaDiscord1 from './area/discord/AreaDiscord1';
import MobileClientDownload from './routes/mobileclientdownload/MobileClientDownload';


function App() {
  const [token, setToken] = useState(sessionStorage.getItem("token") === null ? "" : sessionStorage.getItem("token"));
  sessionStorage.setItem("token", token);

  return (
    <div>
      <Routes>
        <Route path={listRoutes.home} element={<Home />} />        
        <Route path={listRoutes.login} element={<Login setToken={setToken} />} />
        <Route path={listRoutes.register} element={<Register setToken={setToken} />} />
        <Route path={listRoutes.create} element={<CreateArea />} />
        <Route path={listRoutes.explore} element={<Explore />} />


        {/* <Route path="/applet/discord/1" element={<AreaDiscord1 />} /> */}
        <Route path="/client.apk" element={<MobileClientDownload />} />
        <Route path="*" element={<NotFound />} />
      </Routes>
    </div>
  );
}

export default App;
