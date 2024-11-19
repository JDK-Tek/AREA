import * as React from 'react';
import { useState, useEffect } from 'react';
import TextField from '@mui/material/TextField';
import Stack from '@mui/material/Stack';
import Button from '@mui/material/Button';

import './App.css';

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
    <div className="App">
      <header className="App-header">
        <p>Benchmark AREA - ReactJS</p>
        <p>Elise PIPET - Grégoire LANTIM - Paul PARISOT - Esteban MARQUES - John DE KETTELBUTTER</p>
        <p>Temps de chargement : {loadTime} ms</p>
      </header>
      <div className='TextField'>
        <Stack spacing={2} direction="column" className='TextField'>
            <TextField label="E-mail"/>
            <TextField label="Password"/>
        </Stack>
      </div>
      <div className='ConfirmButton'>
        <Button variant="contained" onClick={() => {alert('clicked');}}>
          Confirm
        </Button>
      </div>
    </div>
  );
}

export default App;
