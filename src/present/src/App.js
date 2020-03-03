import React, { Component } from 'react';
import './App.css';
import example from './grpc/example/example_grpc_web_pb';

class App extends Component {

  render() {
    console.log(example);
    return (
      <div className="App">
        <h2>Example application</h2>
      </div>
    );
  }
}

export default App;
