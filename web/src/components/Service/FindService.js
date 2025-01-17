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

export default function FindService({ dark, setService }) {
    const [error, setError] = useState("");
    const [errorMsg, setErrorMsg] = useState("");

    const [search, setSearch] = useState("");
    const [filteredServices, setFilteredServices] = useState([]);

    const [aboutjson, setAboutjson] = useState(null);

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
            const getAboutJson = async () => {
                axios.get(`${process.env.REACT_APP_BACKEND_URL}/about.json`, { headers: { "Content-Type": "application/json" } })
                    .then((response) => {
                        setAboutjson(response.data);
                    })
                    .catch((error) => {
                        setError(true);
                        setErrorMsg("Error when trying to get about.json: " + error);
                    });
            };
            getAboutJson();
    
        }, [aboutjson]);
    

    useEffect(() => {
        if (!aboutjson) return;

        if (search === "") {
            setFilteredServices(aboutjson.server.services);
            return;
        }

        let fstmp = [];
        aboutjson.server.services.forEach((service) => {
            if (matchPattern(search, service.name)) {
                fstmp.push(service);
            }
        });
        setFilteredServices(fstmp);
    }, [search, setFilteredServices, aboutjson]);

    return (
        <div className="h-full flex flex-col justify-start">
            {error && <Notification error={true} msg={errorMsg} setError={setError}/>}

            <SearchInput
                placeholder={"Search for a service"}
                setText={setSearch}
    
                bgColor={mode.bgColor}
                txtColor={mode.txtColor}
                iconColor={mode.iconColor}
                borderColor={mode.borderColor}
            />
            
            <div className="mt-5 overflow-y-auto max-h-[calc(85vh-4rem-64px)]">
                <ServiceKit
                    services={filteredServices}
                    gap={"gap-2"}
                    centered={true}
                    bgColor={dark ? "bg-[#1d1d1d]" : ""}
                    setService={setService}
                    />
            </div>
        </div>
    )
}
