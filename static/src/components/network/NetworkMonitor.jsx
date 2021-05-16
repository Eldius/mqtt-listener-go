
import "./NetworkMonitor.css"
import React from "react"
import NetworkService from "../../service/network/NetworkService"
import {
    LineChart,
    Line,
    CartesianGrid,
    XAxis,
    YAxis,
    ResponsiveContainer,
    Tooltip,
    Legend,
} from 'recharts';

import moment from 'moment'

class NetworkMonitor extends React.Component {

    constructor(props) {
        super(props);
        this.state = {
            data: []
        };
        this.setState = this.setState.bind(this);
    }

    componentDidMount() {
        NetworkService.networkData().then((data) => {
            this.setState({
                data: data.data.map((d) => {
                    console.log(d.Timestamp + " => d: " + d.Payload.download.bandwidth + " / u: " + d.Payload.upload.bandwidth);
                    return {
                        timestamp: Date.parse(d.Timestamp),
                        download: d.Payload.download.bandwidth,
                        upload: d.Payload.upload.bandwidth
                    }
                })
            })
        }).catch(e => {
            console.log(e);
        });
    }

    render() {
        return (
            <div className="NetworkMonitor" >
                <h2>NetworkMonitor</h2>
                <ResponsiveContainer width='100%' aspect={3.0 / 1.0}>
                    <LineChart data={this.state.data} >
                        <Line type="monotone" dataKey="download" stroke="#00ff00" name="Download" />
                        <Line type="monotone" dataKey="upload" stroke="#ff0000" name="Upload" />
                        <CartesianGrid stroke="#ccc" />
                        <XAxis
                            dataKey="timestamp"
                            type='number'
                            tickFormatter={(unixTime) => moment(unixTime).format('DD/MM HH:mm:ss')}
                            domain={['auto', 'auto']}
                            name='Time'
                        />
                        <YAxis unit="Mbps" />
                        <Tooltip className="networkTooltip" labelFormatter={(unixTime) => "Time: " + moment(unixTime).format('DD/MM HH:mm:ss')} />
                        <Legend />
                    </LineChart>
                </ResponsiveContainer>
            </div>
        )
    }
}

export default NetworkMonitor
