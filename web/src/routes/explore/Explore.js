/*
** EPITECH PROJECT, 2024
** AREA
** File description:
** Explore
*/

import axios from "axios";

import React, { useState, useEffect } from "react";
import Notification from "../../components/Notification";
import HeaderBar from "../../components/Header/HeaderBar";
import FindService from "../../components/Service/FindService";
import { backendUrl } from "../../App";

export default function Explore() {
    const [service, setService] = useState(null);
    const [aboutjson, setAboutjson] = useState(null);
    const [error, setError] = useState(false);

    useEffect(() => {
        const getServices = async () => {
            axios
                .get(`${backendUrl}/about.json`, {
                    headers: { "Content-Type": "application/json" },
                })
                .then(response => {
                    setAboutjson(response.data);
                })
                .catch(error => {
                    setError("Error when trying to get all services: " + error);
                });
        };
        getServices();
    }, [setAboutjson]);

    useEffect(() => {
        if (service) {
            window.location.href = `/service/${service.name}`;
        }
    }, [service]);

    return (
        <div>
            <HeaderBar activeBackground={true}/>
            <h1 className="text-4xl font-spartan font-bold text-center mt-10 text-chartgray-300">Explore</h1>
            <div className="mt-5 p-10">
                {error && <Notification error={true} setError={setError} msg={error} />}
                <FindService
                    dark={false}
                    setService={setService}
                    aboutjson={aboutjson}
                />
            </div>
        </div>
    );
}
