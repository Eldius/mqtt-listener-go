
import "./NetworkMonitor.css"
import React, { useState, useEffect } from "react"
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

const refreshInterval = (process.env.REACT_APP_REFRESH_INTERVAL_SECONDS ? process.env.REACT_APP_REFRESH_INTERVAL_SECONDS : 60) * 1000;

const NetworkMonitor = props => {

    const refreshDataFunc = (cb) => {
        NetworkService.networkData().then((data) => {
            cb({
                data: data.data.map((d) => {
                    console.log(d.Timestamp + " => d: " + d.Payload.download.bandwidth + " / u: " + d.Payload.upload.bandwidth);
                    return {
                        timestamp: Date.parse(d.Timestamp),
                        download: d.Payload.download.bandwidth,
                        upload: d.Payload.upload.bandwidth
                    }
                }),
                updated: moment.now()
            })
        }).catch(e => {
            console.log(e);
            cb({
                data: [],
                updated: moment.now()
            })
        });
    }

    const [data, setData] = useState({
        data: [],
        updated: moment.now()
    })

    useEffect(() => {
        console.log("refresh interval: " + refreshInterval);
        refreshDataFunc(setData);
        const id = setInterval(() => {
            console.log("refresh interval: " + refreshInterval);
            refreshDataFunc(setData);
        }, refreshInterval);
        return () => clearInterval(id);
    }, []);

    return (
        <div className="NetworkMonitor" >
            <h2>NetworkMonitor</h2>
            <div className="NetworkResponsiveContainer" >
                <ResponsiveContainer width='100%' aspect={3.0 / 1.0}>
                    <LineChart data={data.data} >
                        <Line type="monotone" dataKey="download" stroke="#00ff00" name="Download" unit=" Mbps" />
                        <Line type="monotone" dataKey="upload" stroke="#ff0000" name="Upload" unit=" Mbps" />
                        <CartesianGrid stroke="#ccc" />
                        <XAxis
                            dataKey="timestamp"
                            type='number'
                            tickFormatter={(unixTime) => moment(unixTime).format('DD/MM HH:mm:ss')}
                            domain={['auto', 'auto']}
                            name='Time'
                        />
                        <YAxis unit="Mbps" contentStyleType="number" />
                        <Tooltip labelStyle={{
                            color: 'black'
                        }} className="networkTooltip" labelFormatter={(unixTime) => "Time: " + moment(unixTime).format('DD/MM HH:mm:ss')} />
                        <Legend />
                    </LineChart>
                </ResponsiveContainer>
            </div>
            <div>
                <label className="netwotkMonitorFooter">atualizado em {moment(data.updated).format('YYYY-MM-DD HH:mm:ss')}</label>
            </div>
        </div>
    )
}

export default NetworkMonitor
