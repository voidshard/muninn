import React from "react";
import PropTypes from 'prop-types';


class LinkRow extends React.Component { // ToDo: investigate rolling Links, Resources & Attrs into one class
    constructor(props) {
        super(props);
    }

    render() {
        let r = this.props.value;
        return (<tr className="clickable" onClick={(event) => {this.props.onLinkClicked(event, r)}}>
            <td>{r.name}</td>
            <td>{r.class}</td>
            <td>{r.subclass}</td>
        </tr>)
    }
}

LinkRow.propTypes = {
    onLinkClicked: PropTypes.func,
    value: PropTypes.object
};


// Represents a single row in the table
export class Links extends React.Component {
    constructor(props) {
        super();
        this.state = {links: []};
    }

    setData(data) {
        this.setState({links: data});
    }

    render() {
        let links = this.state.links;
        if (links.length < 1 || links === undefined) {
            return (<div></div>)
        }

        let rows = [];
        links.forEach((r) => {
            rows.push(<LinkRow onLinkClicked={this.props.onLinkClicked.bind(this)} value={r} />);
        });

        return (
            <div>
                <label>Links:</label> <br />
                <div className="dentright">
                    <table>
                        {rows}
                    </table>
                </div>
            </div>
        )
    }
}
