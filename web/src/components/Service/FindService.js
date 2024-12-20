/*
** EPITECH PROJECT, 2024
** AREA
** File description:
** FindService
*/

import axios from "axios";
import { useEffect, useState } from "react";

import SearchInput from '../SearchInputBox'
import ServiceKit from "./ServiceKit";
import Notification from '../Notification'

export default function FindService({ dark }) {
    const [services, setServices] = useState([]);
    const [error, setError] = useState("");
    const [search, setSearch] = useState("");

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

    return (
        <div clasName="flex flex-wrap justify-center items-center">
            {error && <Notification error={true} msg={error} setError={setError}/>}

            <SearchInput
                placeholder={"Search for a service"}
                setText={setSearch}
    
                bgColor={mode.bgColor}
                txtColor={mode.txtColor}
                iconColor={mode.iconColor}
                borderColor={mode.borderColor}
            />
            <ServiceKit services={services}/>
        </div>
    )
}
