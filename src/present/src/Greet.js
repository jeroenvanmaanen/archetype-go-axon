import React, { Component } from 'react';
import example from './grpc/example/example_grpc_web_pb';
// import example_message from './grpc/example/example_pb';

class Greet extends Component {
    constructor(props) {
        super(props);
        this.handleSubmit = this.handleSubmit.bind(this);
        console.log('Example grpc-web stub:', example);
//        console.log('Example message:', example_message);
    }

    render() {
        return (
            <div>
                <h3>Greetings</h3>
                <p><input type='text' id='message' /> <input type='submit' id='submit-greeting' value=' Go! ' onClick={this.handleSubmit}/></p>
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
        response.on('data', function(e,r) {console.log('Event:', e, r);})
    }
}

export default Greet;
