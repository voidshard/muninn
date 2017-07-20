import React from "react";
import PropTypes from 'prop-types';


class ResourceRow extends React.Component {
    constructor(props) {
        super(props);
    }

    render() {
        let r = this.props.value;
        return (<tr>
            <td>{r.name}</td>
            <td>{r.class}</td>
            <td>{r.uri}</td>
        </tr>)
    }
}

ResourceRow.propTypes = {value: PropTypes.object};


// Represents a single row in the table
export class Resources extends React.Component {
    constructor(props) {
        super();
        this.state = {resources: []};
    }

    setData(data) {
        this.setState({resources: data});
    }

    render() {
        let resources = this.state.resources;
        if (resources.length < 1 || resources === undefined) {
            return (<div></div>)
        }

        let rows = [];
        resources.forEach((r) => {
            rows.push(<ResourceRow value={r} />);
        });

        return (
            <div>
                <label>Resources:</label> <br />
                <div className="dentright">
                    <table>
                        {rows}
                    </table>
                </div>
            </div>
        )
    }
}
