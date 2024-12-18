import * as React from 'react';
import { BrowserRouter, Routes, Route } from "react-router-dom";
import { useState, useEffect } from 'react';
import TextField from '@mui/material/TextField';
import Stack from '@mui/material/Stack';
import Button from '@mui/material/Button';
import './App.css';
import Home from './Home';


function Login() {
  return (
    <div>
      <div className='TextField'>
        <Stack spacing={2} direction="column" className='TextField'>
          <TextField label="E-mail" type="email" />
          <TextField label="Password"/>
        </Stack>
      </div>
      <div className='ConfirmButton'>
        <Button variant="contained" onClick={() => {window.location.href = '/home';}}>
          Confirm
        </Button>
      </div>
    </div>
  );
}

function App() {
  const [loadTime, setLoadTime] = useState(0);

  useEffect(() => {
    const [navigation] = performance.getEntriesByType('navigation');
    if (navigation) {
      const loadTime = navigation.domContentLoadedEventEnd - navigation.startTime;
      setLoadTime(loadTime);
    }
  }, []);

  return (
    <BrowserRouter>
    <div className="App">
      <header className="App-header">
        <p>Benchmark AREA - ReactJS</p>
        <p>Elise PIPET - Grégoire LANTIM - Paul PARISOT - Esteban MARQUES - John DE KETTELBUTTER</p>
        <p>Temps de chargement : {loadTime} ms</p>
      </header>
      <Routes>
        <Route path="/" element={<Login/>} />
        <Route path="/home" element={<Home/>} />
      </Routes>
    </div>
    </BrowserRouter>
  );
}

export default App;
