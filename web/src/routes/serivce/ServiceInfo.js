/*
** EPITECH PROJECT, 2025
** AREA
** File description:
** ServiceInfo
*/

import axios from "axios";
import { useEffect, useState } from "react";
import { useParams } from "react-router-dom";

import Footer from "../home/Footer";
import HeaderBar from "../../components/Header/HeaderBar";
import Notification from "../../components/Notification";
import Button from "../../components/Button";

import { Undo2 } from 'lucide-react';
import { backendUrl } from "../../App";

export default function ServiceInfo() {
    const { service } = useParams();
    const [serviceInfo, setServiceInfo] = useState(null);
    const [error, setError] = useState(false);

    useEffect(() => {
        axios
            .get(`${backendUrl}/api/services/${service}`, {
                headers: { "Content-Type": "application/json" },
            })
            .then(response => {
                setServiceInfo(response.data);
            })
            .catch(error => {
                setError("Error when trying to get service info: " + error);
            });
    }, [service, setServiceInfo]);

    return (
        <div>
            <HeaderBar activeBackground={true} />
            {error && <Notification error={true} setError={setError} msg={error} />}
            <div className="m-10">
                <Button 
                    icon={<Undo2 />}
                    text="Back"
                    onClick={() => window.location.href = "/explore"}
                    styleClolor={`bg-chartpurple-200 hover:bg-chartpurple-100 text-white`}
                />
            </div>
            <div className="mt-10 p-10 flex justify-center">
                <div className="flex flex-col justify-center items-center text-center">
                    <img
                        className="w-[125px] h-[125px] p-5 m-4 rounded-lg"
                        src={`${backendUrl}${serviceInfo?.image}`}
                        alt={service}
                        style={{ backgroundColor: serviceInfo?.color }}
                    />
                    <h1
                        className="font-spartan font-bold text-5xl mt-4 m-4"
                        style={{ color: serviceInfo?.color }}
                    >
                        {service.charAt(0).toUpperCase() + service.slice(1).toLowerCase()}
                    </h1>
                    <div className="lg:w-[900px] md:w-[700px] sm:w-[500px] w-[300px]
                        p-5 m-4 rounded-lg"
                        style={{backgroundColor: serviceInfo?.color}}>
                        <p className="p-5 m-4 font-bold font-spartan text-2xl text-white"
                            style={
                                {
                                    textAlign: "center"
                                }
                            }>
                        Lorem ipsum dolor sit amet, consectetur adipiscing elit. Donec dignissim est lacus, sit amet aliquet nisi ornare quis. 
                        </p>
                    </div>
                </div>
            </div>
            <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 gap-4 mb-[50px]">
                {serviceInfo?.areas.map((area, index) => (
                    <div key={index} className="flex justify-center items-center text-center transition-transform duration-200"
                        onMouseEnter={(e) => { e.currentTarget.style.transform = "scale(1.05)"; }}
                        onMouseLeave={(e) => { e.currentTarget.style.transform = "scale(1)";}}
                        onClick={() => window.location.href = `/create`}
                    >
                        <div className="p-[50px] rounded-lg w-[400px] flex justify-center items-center m-4"
                        style={{ backgroundColor: serviceInfo?.color }}>
                            <label className="
                                font-spartan font-bold text-2xl
                                text-white
                            ">
                                {area.description.charAt(0).toUpperCase() + area.description.slice(1).toLowerCase()}
                            </label>
                        </div>
                    </div>
                ))}
            </div>
            <Footer />
        </div>
    );
}
