import * as React from 'react';
import TextField from '@mui/material/TextField'
import Stack from '@mui/material/Stack';
import Button from '@mui/material/Button'

import './App.css';

function App() {
  return (
    <div className="App">
      <header className="App-header">
        <p>Benchmark AREA - ReactJS</p>
        <p>Elise PIPET - Grégoire LANTIM - Paul PARISOT - Esteban MARQUES - John DE KETTELBUTTER</p>
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
