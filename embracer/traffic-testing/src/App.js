import React from 'react';
import './App.css';
import TrafficTest from './Traffic';
import { TrafficProvider } from './TrafficContext';

function App() {
  return (
    <div>
      <TrafficProvider>
        <TrafficTest />
      </TrafficProvider>
    </div>
  );
}

export default App;
