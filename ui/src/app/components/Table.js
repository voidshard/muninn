import React from "react";
import PropTypes from 'prop-types'

import { Search } from "./Search";
import { Row } from "./Row";
import { DataView } from "./View";
import { dataAccessService } from "../dataAccess";
import { isWheelUp } from "../utils";

const defaultRowsDisplayed = 35;

function defaultState(rowsDisplayed) {
    // ToDo: dispSize should be set dynamically
    let dispSize = rowsDisplayed? rowsDisplayed: defaultRowsDisplayed;
    return {
        rows: undefined, // map[int][]AssetDescription, ie: {page number: rows}
        searchTerms: undefined, // current search terms (if any)
        dispSize: dispSize, // num rows to display at once
        first: 0, // lowest currently displayed row number
        last: rowsDisplayed, // highest currently displayed row number
        pageSize: 100, // number of rows fetched at once
        rowsPerWheel: 3, // rows to move up / down per scroll of the wheel
    }
}
export class Table extends React.Component {
    constructor(props) {
        super();
        this.state = defaultState(props.rowsDisplayed);
    }

    // Search given some search terms -- this resets our table of data from page 0
    onSearch(terms) {
        if (terms.length < 1) {
            console.log("Warning: Table onSearch: No search terms given");
            return;
        }
        dataAccessService.search(terms, 0).then((data) => {
            let newState = defaultState();
            newState.rows = {0: data};
            newState.searchTerms = terms;
            newState.first = 0;
            newState.last = Math.min(data.length, this.state.dispSize);
            newState.pageSize = data.length;
            this.setState(newState)
        }).catch((err) => {
            console.log("Table onSearch", terms, 0, err);
        });
    }

    // Given some page number, load the page from the server
    loadPage(page) {
        return new Promise((resolve, reject) => {
            if (this.state.searchTerms === undefined) {
                return; // we have no search terms so we can't load more data
            }

            dataAccessService.search(this.state.searchTerms, page).then((data) => {
                let newRows = this.state.rows;
                newRows[page] = data;
                this.setState({rows: newRows});
                resolve(data);
            }).catch((err) => {
                console.log("Table loadPage", page, err);
                reject(err);
            });
        });
    }

    // Given some page number, remove it from our clientside memory
    unloadPage(page) {
        let newRows = this.state.rows;
        if (!newRows.hasOwnProperty(page)) {
            return;
        }
        console.log("Table unloadPage", page);
        delete newRows[page];
        this.setState({rows: newRows});
    }

    // Handler for any currently displayed row being clicked
    onClick(event, i) {
        this.refs.dataview.display(i);
    }

    // Handler for mouse scroll events on the table.
    // Increments the page pointer (to the top row of the table 'first') and decides
    // the 'last' row. Also decides if we need to load / unload pages of data from the server.
    onWheel(e) {
        let up = isWheelUp(e);
        let mult = up? -1: 1;

        let first = this.state.first + (mult * this.state.rowsPerWheel);
        if (first < 0) {
            first = 0;
        }
        let last = first + this.state.dispSize;

        // Some basic math to figure out what page numbers we'll need
        let firstPage = Math.floor(first / this.state.pageSize);
        let lastPage = Math.floor(last / this.state.pageSize);

        // Depending on what direction we're scrolling, we'll need to ensure we have some page number
        let check_page = -1;
        let remove_page = -1;
        if (up) {
            check_page = firstPage;
            remove_page = lastPage + 1
        } else {
            check_page = lastPage;
            remove_page = firstPage - 1
        }

        this.unloadPage(remove_page); // be nice and remove old page data

        if (!this.state.rows.hasOwnProperty(check_page)) {
            this.loadPage(check_page).then((data) => {
                this.setState({
                    first: first,
                    last: last,
                });
            }).catch((err) => {
                console.log("Table onWheel", check_page, err);
                reject(err);
            });
        } else {
            this.setState({
                first: first,
                last: last
            });
        }
    }

    rowsBetween(first, last) {
        let firstRow = first % this.state.pageSize;
        let firstPage = Math.floor(first / this.state.pageSize);
        let lastRow = last % this.state.pageSize;
        let lastPage = Math.floor(last / this.state.pageSize);
        let maxRowNum = firstRow + this.state.dispSize;

        if (this.state.rows === undefined || first === undefined) {
            return [];
        }

        // first set of rows
        let rows = this.state.rows[firstPage].slice(firstRow, Math.min(maxRowNum, this.state.pageSize));

        if (firstPage != lastPage) { // if the two points (first & last) aren't on the same page
            if (this.state.rows.hasOwnProperty(lastPage)) { // assuming we've got another page
                // Add rows from the next page
                rows = rows.concat(this.state.rows[lastPage].slice(0, lastRow));
            }
        }

        // If we don't have enough rows, move first backwards and try again
        if (rows.length < this.state.dispSize && first > 0) {
            first -= this.state.dispSize - rows.length;
            if (first < 0) {
                first = 0;
            }
            last = first + this.state.dispSize;
            this.setState({
                first: first,
                last: last
            });
            return this.rowsBetween(first, last);
        }
        return rows;
    }

    render() {
        let table = (<div> </div>);
        let rows = this.rowsBetween(this.state.first, this.state.last);

        if (rows.length > 0) {// if there is page data
            let items = rows.map((i) => {
                // ToDo: Find better way over adding listener to each row :(
                return (
                    <Row onClick={this.onClick.bind(this)} item={i}/>
                );
            });

            table = (
                <table itemID="datatable" onScroll={this.onWheel.bind(this)} onWheel={this.onWheel.bind(this)}>
                    <thead>
                        <tr>
                            <th>Name</th>
                            <th>Type</th>
                            <th>Subtype</th>
                            <th>Description</th>
                        </tr>
                    </thead>
                    <tbody>
                        {items}
                    </tbody>
                </table>
            )
        }

        return (
            <div>
                <div className="datatable left">
                    <Search onSearch={this.onSearch.bind(this)} />
                    <br />
                    {table}
                </div>
                <DataView ref="dataview" setData={this.onClick.bind(this)}/>
            </div>
        );
    }
}

Table.propTypes = {
    rowsDisplayed: PropTypes.number
};
