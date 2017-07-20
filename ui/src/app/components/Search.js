import React from "react";
import PropTypes from 'prop-types';
import { dataAccessService } from "../dataAccess";

const RgxAlphaNumericSlash = /^([a-zA-Z0-9]|\/)$/g; // Match any number, any letter and '/'
const SearchTip = "Search for an asset via it's 'collection', 'collection/type' or 'collection/type/variant'";

export class Search extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            suggestions: [],
            terms: [],
        };
        this.init();
    }

    init() {
        dataAccessService.suggest().then((data) => {
            if (data.length) {
                let terms = [data[0]];
                this.setState({
                    suggestions: data,
                    terms: terms,
                });
                this.props.onSearch(terms);
            }
        }).catch((err) => {
            console.log("Search init", err);
        });
    }

    onChange(e) { // ToDo: this may not need to be onChange
        if (e.key == "Enter") {
            let unsanitized = e.target.value; // string with who knows what chars
            let sanitized = ""; // string of acceptable chars

            for (let i=0; i < unsanitized.length; i++) {
                let char = unsanitized[i];
                if (!!char.match(RgxAlphaNumericSlash)) {
                    sanitized += char;
                }
            }

            // Trim '/' chars from start and end
            sanitized = sanitized.replace(/^\/+|\/+$/gm,'');

            // Remove duplicate occurences of '/' char
            while (sanitized.indexOf("//") > -1) { // ToDo: there probably is a nicer way..
                sanitized = sanitized.replace("//", "/");
            }

            // Finally, split on remaining '/' symbols and pass to search
            let terms = sanitized.split("/");
            this.setState({terms: terms});
            this.props.onSearch(terms);
        }
    };

    render() {
        let terms = (<div></div>);
        if (this.state.terms) {
            terms = (<div className="dentright">
                <label>Search results for: {this.state.terms.join("/")}</label>
            </div>)
        }

        return (
            <div>
                <div className="inline">
                    <input title={SearchTip}
                           placeholder="Search for assets"
                           onKeyDown={this.onChange.bind(this)} />
                    {terms}
                </div>
            </div>
        );
    };
};

Search.propTypes = {
    query: PropTypes.string,
    onSearch: PropTypes.func
};
