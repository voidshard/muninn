import React from "react";
import PropTypes from 'prop-types';

import { dataAccessService } from "../dataAccess";

import { Thumbnail } from "./Thumbnail";
import { Attributes } from "./Attributes";
import { Resources } from "./Resources";
import { Links } from "./Links";


const btnBack = "resources/back.png";


// Represents the display window on the right side of the screen
export class DataView extends React.Component {
    constructor(props) {
        super();
        this.state = {
            asset: undefined,
            previous: [],
        };
    }

    display(assetDesc) {
        dataAccessService.fetch(assetDesc).then((asset) => {
            this.setState({asset: asset});

            this.refs.thumbnail.setThumbnail(asset.thumb);
            this.refs.attributes.setData(asset.attrs);
            this.refs.resources.setData(asset.resources);
            this.refs.links.setData(asset.linked);
        }).catch((err) => {
           console.log("View display", err);
           reject(err);
        });
    }

    onBackClicked(event) {
        let prev = this.state.previous;
        this.display(prev.pop());
        this.setState({previous: prev});
    }

    onLinkClicked(event, link) {
        let prev = this.state.previous;
        prev.push(this.state.asset.description);
        this.setState({previous: prev}); // record the last selected asset
        this.display(link);
    }

    render() {
        let asset = this.state.asset;
        if (asset == undefined) {
            return (<div className="dataview"></div>)
        }

        let prev = this.state.previous;
        let back = (<div><br /><br /></div>);
        if (prev.length > 0) {
            back = (<div>
                <br/> <img onClick={this.onBackClicked.bind(this)}
                           src={btnBack}
                           className="icon clickable"
                           title="Previous asset"/>
            </div>);
        }

        return (
            <div className="dataview right">
                <table>
                    {/* ToDo: probably should use div over a table here .. */}
                    <tr>
                        <td>
                            <Thumbnail ref="thumbnail" />
                        </td>
                        <td>
                            <div className="dentright">
                                {back}
                                <label>
                                    {asset.description.name}
                                    /{asset.description.class}
                                    /{asset.description.subclass}
                                </label>
                                <br />
                                <label>
                                    v{asset.version}
                                </label>
                            </div>
                        </td>
                    </tr>
                </table>
                <br /> <br />
                <div className="subcontainer">
                    <Attributes ref="attributes"/>
                    <br />
                    <Resources ref="resources"/>
                    <br />
                    <Links onLinkClicked={this.onLinkClicked.bind(this)} ref="links"/>
                </div>
            </div>
        )
    }
}

DataView.propTypes = {
    setData: PropTypes.func
};
