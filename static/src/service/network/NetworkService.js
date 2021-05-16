//import samples from "./sample"

const baseUrl = process.env.REACT_APP_BACKEND_ENDPOINT;


console.log(baseUrl);

class NetworkService {
    async networkData() {
        return fetch(baseUrl + "/last?t=/network/modem")
            .then(d => d.json());
    }
}

/* async function FetchNetworkData() {
    return fetch(baseUrl + "/last?t=/network/modem")
        .then(d => d.json());
} */



export default new NetworkService();
