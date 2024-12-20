/*
** EPITECH PROJECT, 2024
** AREA
** File description:
** FindService
*/

import axios from "axios";
import { useEffect, useState } from "react";

import matchPattern from "../../utils/matchPattern";
import SearchInput from '../SearchInputBox'
import ServiceKit from "./ServiceKit";
import Notification from '../Notification'

export default function FindService({ dark }) {
    const [services, setServices] = useState([]);
    const [error, setError] = useState("");

    const [search, setSearch] = useState("");
    const [filteredServices, setFilteredServices] = useState([]);

    const mode = dark ?
        {
            bgColor: "bg-chartgray-200",
            txtColor: "text-white placeholder-cahrtgray-100",
            iconColor: "text-cahrtgray-100",
            borderColor: "border-chartgray-100 focus:border-blue-500"
        } : 
        {
            bgColor: "bg-gray-50",
            txtColor: "text-gray-900 placeholder-gray-400",
            iconColor: "text-cahrtgray-100",
            borderColor: "border-gray-300 focus:border-blue-500"
        }

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
        getServices();
    }, [setServices, setError]);

    useEffect(() => {
        if (search === "") {
            setFilteredServices(services);
            return;
        }

        let fstmp = [];
        services.forEach((service) => {
            if (matchPattern(search, service.name)) {
                fstmp.push(service);
            }
        });
        setFilteredServices(fstmp);
    }, [search, setFilteredServices, services]);

    return (
        <div className="h-full flex flex-col justify-start">
            {error && <Notification error={true} msg={error} setError={setError}/>}

            <SearchInput
                placeholder={"Search for a service"}
                setText={setSearch}
    
                bgColor={mode.bgColor}
                txtColor={mode.txtColor}
                iconColor={mode.iconColor}
                borderColor={mode.borderColor}
            />
            
            <div className="mt-5 overflow-y-auto max-h-[calc(85vh-4rem-64px)] w-full flex flex-col">
                <ServiceKit
                    services={filteredServices}
                    gap={"gap-1"}
                    centered={false}
                />
            </div>
        </div>
    )
}
