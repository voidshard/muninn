
const apiSuggest = "/api/1/suggest";
const apiSearch = "/api/1/suggest/";
const apiFetch = "/api/1/fetch/";

// DataAccess facilitates fetching server side data
class DataAccess {
    initialSuggestions() {
        return new Promise((resolve, reject) => {
            let url = apiSuggest;
            fetch(url).then((response) => {
                response.json().then((data) => {
                    console.log("[DataAccessAPI Suggest]", url, data);
                    return resolve(data);
                });
            }).catch((error) => {
                console.log("[DataAccessAPI Suggest]", error);
                return reject(error);
            });
        });
    }

    search(terms, page) {
        return new Promise((resolve, reject) => {
            let url = apiSearch + page + "/" + terms.join("/");
            fetch(url).then((response) => {
                response.json().then((data) => {
                    console.log("[DataAccessAPI Search]", page, url, data);
                    let assets = [];
                    data.forEach((i) => {assets.push(new AssetDescription(i))});
                    return resolve(assets);
                });
            }).catch((error) => {
                console.log("[DataAccessAPI Search]", error);
                return reject(error);
            });
        });
    }

    fetch(asset) {
        return new Promise((resolve, reject) => {
            let terms = [asset.name, asset.class, asset.subclass];
            let url = apiFetch + terms.join("/");

            fetch(url).then((response) => {
                response.json().then((data) => {
                    console.log("[DataAccessAPI Fetch]", asset, url, data);
                    return resolve(new AssetData(data));
                });
            }).catch((error) => {
                console.log("[DataAccessAPI Fetch]", error);
                return reject(error);
            });
        });
    }
}

// Rough desc of an asset w/ enough info to identify one uniquely
function AssetDescription(rawData) {
    let self = this;
    self.name = rawData.Name;
    self.class = rawData.Class;
    self.subclass = rawData.Subclass;
    self.description = rawData.Description;
}

// Full asset data, including the 'asset description' info above
function AssetData(rawData) {
    let self = this;
    self.description = new AssetDescription(rawData.Data);

    self.attrs = rawData.Attributes;
    self.version = rawData.Version;
    self.thumb = rawData.Thumbnail;

    let linked = [];
    if (rawData.Linked) {
        rawData.Linked.forEach((i) => {
            linked.push(new AssetDescription(i))
        });
    }
    self.linked = linked;

    let res = [];
    if (rawData.Resources) {
        rawData.Resources.forEach((i) => {
            res.push(new ResourceDescription(i))
        });
    }

    self.resources = res;
}

// Description of some resource attached to our asset
function ResourceDescription(rawData) {
    let self = this;
    self.name = rawData.Name;
    self.class = rawData.Class;
    self.uri = rawData.URI;
}

// Service singleton
const service = new DataAccess();

// Exported obj (holds the public functions of our singleton)
let dataAccessService = {
    suggest: service.initialSuggestions,
    search: service.search,
    fetch: service.fetch,
};

// Export public facing service access
export {dataAccessService};
