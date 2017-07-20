import React from "react";
import { render } from "react-dom";

import { Footer } from "./components/Footer";
import { Table } from "./components/Table";


// Main / root component
class App extends React.Component {
    render() {
        return (
            <div className="container">
                <div>
                    <Table/>
                </div>
                <Footer/>
            </div>
        );
    }
}

render(<App/>, window.document.getElementById("app"));
