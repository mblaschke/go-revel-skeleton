import React, { Component } from 'react';
import { BrowserRouter as Router, Route } from 'react-router-dom';
import HelloWorld from './HelloWorld';
import $ from "jquery";

class App extends Component {
    constructor(props) {
        super(props);

        this.state = {};
    }

    render() {
        return (
            <Router>
                <div>
                    <Route path="/" component={HelloWorld} />
                </div>
            </Router>
        )
    }
}

export default App;
