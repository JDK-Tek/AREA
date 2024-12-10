/*
** EPITECH PROJECT, 2024
** AREA
** File description:
** Content
*/

import { useEffect, useState } from "react";
import fetchData from "../../utils/fetchData";

import Button from "../../components/Button";
import AppletKit from "./../../components/Applet/AppletKit";
import ServiceKit from "./../../components/Service/ServiceKit";

export default function Content({ setError }) {
    const [services, setServices] = useState([]);
    const [applets, setApplets] = useState([]);


    useEffect(() => {
        const fetchServices = async () => {
            const { success, data, error } = await fetchData("http://localhost:42000/api/services");
            if (!success) {
                setError("Error while fetching services data: " + error);
            } else {
                setServices(data.res);
            }
        };

        const fetchApplets = async () => {
            const { success, data, error } = await fetchData("http://localhost:42000/api/applets");
            if (!success) {
                setError("Error while fetching applets data: " + error);
            } else {
                setApplets(data.res);
            }
        };

        if (services.length === 0) {
            fetchServices();
        }
        if (applets.length === 0) {
            fetchApplets();
        }
    
    }, [applets, services, setError]);
    
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
