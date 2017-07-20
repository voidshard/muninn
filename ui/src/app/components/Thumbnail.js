import React from "react";


const noThumb = "resources/nothumb.png";


// Represents a single row in the table
export class Thumbnail extends React.Component {
    constructor(props) {
        super(props);
        this.state = {image: undefined};
    }

    setThumbnail(src) {
        this.setState({image: src});
    }

    render() {
        let image = this.state.image;
        if (this.state.image === undefined || this.state.image == "") {
            image = noThumb;
        }

        return (
            <div className="shadow">
                <img className="thumbnail" src={image}/>
            </div>
        )
    }
}
