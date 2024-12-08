/*
** EPITECH PROJECT, 2024
** AREA
** File description:
** App
*/

import React from 'react';

import { Routes, Route } from 'react-router-dom';

import Home from "./routes/home/Home";
import Login from "./routes/login/Login";
import Register from "./routes/register/Register";
import NotFound from "./routes/notfound/NotFound";

function App() {

  return (
    <div>
      <Routes>
        <Route path="/" element={<Home />} />
        <Route path="/login" element={<Login />} />
        <Route path="/register" element={<Register />} />
        <Route path="*" element={<NotFound />} />
      </Routes>
    </div>
  );
}

export default App;
