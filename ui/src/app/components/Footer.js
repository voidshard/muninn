import React from "react";

export class Footer extends React.Component {
    constructor(props) {
        super();
        this.state = {
            message: props.message !== undefined? props.message: "",
        }
    }

    render () {
        return (
            <div className="footer">
                {this.state.message}
            </div>
        );
    }
}
