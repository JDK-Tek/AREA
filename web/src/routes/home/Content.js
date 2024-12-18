/*
** EPITECH PROJECT, 2024
** AREA
** File description:
** Content
*/


import { useEffect, useState } from "react";

import axios from "axios";
import Button from "../../components/Button";
import AppletKit from "./../../components/Applet/AppletKit";
import ServiceKit from "./../../components/Service/ServiceKit";

export default function Content({ setError }) {
    const [services, setServices] = useState([]);
    const [applets, setApplets] = useState([]);
    
    useEffect(() => {
        const getServices = async () => {
            axios.get(`${process.env.REACT_APP_BACKEND_URL}/api/services`, {headers: {"Content-Type": "application/json"}})
            .then((response) => {
                setServices(response.data.res)
            })
            .catch((error) => {
                setError("Error when trying to get all services: " + error)
            });
        };
    
        const getApplets = async () => {
            axios.get(`${process.env.REACT_APP_BACKEND_URL}/api/applets`, {headers: {"Content-Type": "application/json"}})
            .then((response) => {
                setApplets(response.data.res)
            })
            .catch((error) => {
                setError("Error when trying to get all applets: " + error)
            });
        };

        getServices();
        getApplets();
    
    }, [setApplets, setServices, setError]);

    return (
        <div className="pb-14">
            <AppletKit
                title={"Get started with any Applet"}
                applets={applets}
            />
            <ServiceKit
                title={"or choose from 900+ services"}
                services={services}
                color={"text-chartpurple-200"}
            />
            <div className="flex justify-center items-center mt-8">
                <Button
                    text={"Explore all"}
                    redirect={false}
                    onClick={() => window.location.href = "/explore"}
                    styleClolor={"bg-chartgray-300 text-white hover:bg-chartgray-200 text-2xl"}
                />
            </div>
        </div>
    );
}
