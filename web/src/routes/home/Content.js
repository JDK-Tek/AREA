/*
** EPITECH PROJECT, 2024
** AREA
** File description:
** Content
*/


import { act, useEffect, useState } from "react";

import axios from "axios";
import Button from "../../components/Button";
import AppletKit from "./../../components/Applet/AppletKit";
import ServiceKit from "./../../components/Service/ServiceKit";
import { backendUrl } from "../../App";

export default function Content({ setError }) {
    const [services, setServices] = useState([]);
    const [applets, setApplets] = useState([]);
    const [service, setService] = useState(null);
    
    useEffect(() => {
        const getServices = async () => {
            axios.get(`${backendUrl}/api/services`, {headers: {"Content-Type": "application/json"}})
            .then((response) => {
                const res = response.data.slice(0, 5);
                setServices(res)
            })
            .catch((error) => {
                setError("Error when trying to get all services: " + error)
            });
        };
        getServices();
    }, [setServices, setError]);
    
    useEffect(() => {
        const getApplets = async () => {
            axios.get(`${backendUrl}/api/applets`, {headers: {"Content-Type": "application/json"}})
            .then((response) => {
                const applets = response.data.res;
            
                let res = [];
                for (let i = 0; i < Math.min(5, applets.length); i++) {
                    res.push({
                        id: applets[i].id,
                        name: applets[i].name,
                        action: {
                            service: applets[i].service.name,
                            image: applets[i].service.logo,
                            color: applets[i].service.color.normal,
                        },
                        reaction: {
                            service: "",
                            image: applets[i].service.logopartner,
                            color: applets[i].service.color.hover,
                        },
                        users: applets[i].users,
                    });
                }
                setApplets(res);
            })
            .catch((error) => {
                setError("Error when trying to get all applets: " + error)
            });
        };
        getApplets();
    }, [setApplets, setError]);
    
    useEffect(() => {
        if (service) {
            window.location.href = `/services/${service.name}`;
        }
    }, [service]);
    
    return (
        <div className="pb-14">
            <AppletKit
                title={"Get started with any Applet"}
                applets={applets}
                onClick={() => console.log("clicked")}
            />
            <ServiceKit
                title={"or choose from 900+ services"}
                services={services}
                color={"text-chartpurple-200"}
                gap={"gap-3"}
                rounded={"rounded-xl"}
                setService={setService}
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
