
import "./NetworkMonitor.css"
import React from "react"
import NetworkService from "../../service/network/NetworkService"
import { LineChart, Line, CartesianGrid, XAxis, YAxis, ResponsiveContainer, ScatterChart, Scatter } from 'recharts';
import moment from 'moment'

class NetworkMonitor extends React.Component {

    constructor(props) {
        super(props);
        this.state = {
            data: []
        };
        this.setState = this.setState.bind(this);
    }

    componentWillMount() {
        NetworkService.networkData().then((data) => {
            this.setState({
                downData: data.data.map((d) => {
                    console.log(d.Timestamp + " => " + d.Payload.download.bandwidth);
                    return {
                        timestamp: Date.parse(d.Timestamp),
                        value: d.Payload.download.bandwidth
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


                <LineChart width={600} height={300} data={this.state.downData}>
                    <Line type="monotone" dataKey="value" stroke="#8884d8" />
                    <CartesianGrid stroke="#ccc" />
                    <XAxis
                        dataKey="timestamp"
                        type='number'
                        tickFormatter={(unixTime) => moment(unixTime).format('HH:mm:ss')}
                        domain={['auto', 'auto']}
                        name='Time' />
                    <YAxis />
                </LineChart>
                <ResponsiveContainer>
                    <ScatterChart width={600} height={300} >
                        <XAxis
                            dataKey='timestamp'
                            domain={['auto', 'auto']}
                            name='Time'
                            tickFormatter={(unixTime) => moment(unixTime).format('HH:mm:ss')}
                            type='number'
                        />
                        <YAxis dataKey='value' name='Value' />
                        <Scatter
                            data={this.state.downData}
                            line={{ stroke: '#eee' }}
                            lineJointType='monotoneX'
                            lineType='joint'
                            name='Values'
                        />
                    </ScatterChart>
                </ResponsiveContainer>
            </div>
        )
    }
}

export default NetworkMonitor
