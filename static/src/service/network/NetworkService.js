//import samples from "./sample"

const baseUrl = process.env.REACT_APP_BACKEND_ENDPOINT;
const fetchQuantity = process.env.REACT_APP_FETCH_COUNT ? process.env.REACT_APP_FETCH_COUNT : 10;


console.log(baseUrl);

class NetworkService {
    async networkData() {
        return fetch(baseUrl + "/last?t=/network/modem&c=" + fetchQuantity)
            .then(d => d.json());
    }
}

/* async function FetchNetworkData() {
    return fetch(baseUrl + "/last?t=/network/modem")
        .then(d => d.json());
} */



export default new NetworkService();
