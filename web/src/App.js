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
import MobileClientDownload from './routes/mobileclientdownload/MobileClientDownload';
import Connected from './routes/connected/Connected';
import Privacy from './routes/privacy/Privacy';


function App() {
  const [token, setToken] = useState(sessionStorage.getItem("token") === null ? "" : sessionStorage.getItem("token"));
  sessionStorage.setItem("token", token);

  return (
    <div>
      <Routes>
        <Route path={listRoutes.home} element={<Home />} />        
        <Route path={listRoutes.login} element={<Login setToken={setToken} />} />
        <Route path={listRoutes.register} element={<Register setToken={setToken} />} />
        <Route path={listRoutes.create} element={<CreateArea setToken={setToken} />} />
        <Route path={listRoutes.explore} element={<Explore />} />
        <Route path={listRoutes.clientapk} element={<MobileClientDownload />} />
        <Route path={listRoutes.connected} element={<Connected />} />
        <Route path={listRoutes.privacy} element={<Privacy />} />

        <Route path="*" element={<NotFound />} />
      </Routes>
    </div>
  );
}

export default App;
