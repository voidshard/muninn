import React from "react";
import PropTypes from 'prop-types';


// Represents a single row in the table
export class Row extends React.Component {
    render() {
        let item = this.props.item;

        return (
            <tr className="clickable" onClick={(event) => {this.props.onClick(event, item)}}>
                <td>{item.name}</td>
                <td>{item.class}</td>
                <td>{item.subclass}</td>
                <td>{item.description}</td>
            </tr>
        )
    }
}

Row.propTypes = {
    key: PropTypes.string,
    item: PropTypes.object,
    onClick: PropTypes.func
};
