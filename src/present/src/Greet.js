import React, { Component } from 'react';
import example from './grpc/example/example_grpc_web_pb';
// import example_message from './grpc/example/example_pb';

class Greet extends Component {
    constructor(props) {
        super(props);
        this.handleSubmit = this.handleSubmit.bind(this);
        this.handleRefresh = this.handleRefresh.bind(this);
        console.log('Example grpc-web stub:', example);
//        console.log('Example message:', example_message);
    }

    render() {
        return (
            <div>
                <h3>Greetings</h3>
                <p><input type='text' id='message' /> <input type='submit' id='submit-greeting' value=' Go! ' onClick={this.handleSubmit}/></p>
                <p><input type='submit' id='refresh-greetings' value=' Refresh! ' onClick={this.handleRefresh}/></p>
                <div id='greetings'><div><i>greetings appear here</i></div></div>
            </div>
        );
    }

    handleSubmit(event) {
        const message = document.getElementById('message').value;
        console.log('Message:', message);
        const request = new example.Greeting();
        console.log('New request:', request);
        request.setMessage(message);
        console.log('Request:', request);
        const client = new example.GreeterServiceClient('http://localhost:3000');
        console.log('Client:', client);
        const response = client.greet(request);
        console.log('Response:', response);
        response.on('data', function(r) {console.log('Greet event:', r);})
    }

    handleRefresh(event) {
        console.log('Handle refresh:', event)

        const client = new example.GreeterServiceClient('http://localhost:3000');
        console.log('Refresh: client:', client);

        const container = document.getElementById('greetings');
        container.innerHTML = '';

        const request = new example.Empty();
        console.log('Refresh: new request:', request);
        const response = client.greetings(request);

        console.log('Refresh: response:', response);
        response.on('data', function(r) {
            console.log('Refresh: greetings event:', r);
            const message = r.getMessage();
            console.log('Refresh: greetings event: message', message);
            const text = document.createTextNode(message);
            const div = document.createElement('div');
            div.appendChild(text);
            container.appendChild(div)
        })
        response.on('status', function(status) {
          console.log('Refresh: stream status: code:', status.code);
          console.log('Refresh: stream status: details:', status.details);
          console.log('Refresh: stream status: metadata:', status.metadata);
        });
        response.on('end', function(end) {
          // stream end signal
        });
    }
}

export default Greet;
