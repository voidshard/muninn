import React from "react";
import PropTypes from 'prop-types';


class AttributeRow extends React.Component {
    constructor(props) {
        super(props);
    }

    render() {
        return (<tr>
            <td>{this.props.name}:</td>
            <td>{this.props.value}</td>
        </tr>)
    }
}

AttributeRow.propTypes = {
    name: PropTypes.string,
    value: PropTypes.string
};


// Represents a single row in the table
export class Attributes extends React.Component {
    constructor(props) {
        super();
        this.state = {attrs: {}};
    }

    setData(data) {
        this.setState({attrs: data});
    }

    render() {
        let attrs = this.state.attrs;
        if (attrs === undefined) {
            return (<div></div>)
        }

        let rows = [];
        for(var somekey in attrs) {
            if (attrs.hasOwnProperty(somekey)) {
                rows.push(<AttributeRow name={somekey} value={attrs[somekey]} />);
            }
        }

        if (rows.length < 1) {
            return (<div></div>)
        }

        return (
            <div>
                <label>Attributes</label> <br />
                <div className="dentright">
                    <table>
                        {rows}
                    </table>
                </div>
            </div>
        )
    }
}
