import * as React from 'react';
import TextField from '@mui/material/TextField'
import Button from '@mui/material/Button'

import './App.css';

function App() {
  return (
    <div className="App">
      <header className="App-header">
        <div className='TextField'>
          <TextField label="E-mail"/>
          <TextField label="Password"/>
          <Button variant="contained" onClick={() => {alert('clicked');}}>
            Confirm
          </Button>
        </div>
      </header>
    </div>
  );
}

export default App;
