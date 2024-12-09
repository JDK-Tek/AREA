/*
** EPITECH PROJECT, 2024
** AREA
** File description:
** Content
*/

import { useEffect, useState } from "react";

import Button from "../../components/Button";
import AppletKit from "./../../components/Applet/AppletKit";
import ServiceKit from "./../../components/Service/ServiceKit";

import AppletData from "./../../data/AppletData";

async function fetchData(url) {
    const request = {
        method: "GET"
    };

    try {
        const res = await fetch(url, request);
        if (!res.ok) {
            throw new Error(`Response status: ${res.status}`);
        }
    
        const json = await res.json();
        return { success: true, data: json };
    } catch (err) {
        return { success: false, error: err };
    }
}

export default function Content({ data }) {
    const [error, setError] = useState(null);
    const [services, setServices] = useState([]);
    const [applets, setApplets] = useState([]);

    useEffect(() => {
    
        fetchData("http://localhost:42000/api/services").then(({ success, data, error }) => {
            if (!success) { setError("Error while fetching services data", error);
            } else { setServices(data.res);}
        });

        fetchData("http://localhost:42000/api/applets").then(({ success, data, error }) => {
            if (!success) { setError("Error while fetching applets data", error);
            } else { setApplets(data.res); console.log(data.res); }
        });
    
    }, []);
    
    return (
        <div className="pb-14">
            <label className="text-1xl font-bold text-red-900">{error}</label>
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
