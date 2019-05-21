import React from 'react';
import $ from 'jquery';

import BaseComponent from './BaseComponent';
import Spinner from './Spinner';

class HelloWorld extends BaseComponent {
    constructor(props) {
        super(props);
        this.state = {};
    }

    render() {
        return (
            <div>Hello World</div>
        )
    }
}

export default HelloWorld;

