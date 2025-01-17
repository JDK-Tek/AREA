/*
** EPITECH PROJECT, 2024
** AREA
** File description:
** Service
*/

import { useEffect, useState } from "react";

export default function Service({ service, rounded, setService }) {
    const [serviceName, setServiceName] = useState("");
    const [serviceImage, setServiceImage] = useState("");
    const [serviceColor, setServiceColor] = useState("");

    useEffect(() => {
        try {
            setServiceName(service.name);
            setServiceImage(service.image);
            setServiceColor(service.color);
        } catch (error) {
            console.error("Error when trying to set service: ", error);
        }
    }, [service, serviceName, serviceImage, serviceColor]);

    return (
        <div
            className={`select-none relative w-[200px] h-[150px] text-white ${rounded} shadow-md p-6 flex flex-col justify-between items-center cursor-pointer transition-transform duration-200`}
            onClick={() => setService(service)}
            style={{
                backgroundColor: serviceColor
            }}
            onMouseEnter={(e) => {
                e.currentTarget.style.backgroundColor = serviceColor;
                e.currentTarget.style.transform = "scale(1.05)";
            }}
            onMouseLeave={(e) => {
                e.currentTarget.style.backgroundColor = serviceColor;
                e.currentTarget.style.transform = "scale(1)";
            }}
        >
            <img 
                className="w-[50px] h-[50px]" 
                src={serviceImage}
                alt={serviceName} 
            />

            <h1 className="font-spartan text-2xl font-bold text-center"> {serviceName} </h1>
        </div>
    );
}

