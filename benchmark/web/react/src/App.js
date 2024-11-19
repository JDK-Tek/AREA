import * as React from 'react';
import TextField from '@mui/material/TextField'

import './App.css';

function App() {
  return (
    <div className="App">
      <header className="App-header">
        <div className='TextField'>
          <TextField label="E-mail"/>
          <TextField label="Password"/>
        </div>
      </header>
    </div>
  );
}

export default App;
