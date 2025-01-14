/*
** EPITECH PROJECT, 2024
** AREA
** File description:
** Explore
*/

import React, { useState, useEffect } from "react";
import HeaderBar from "../../components/Header/HeaderBar";
import FindService from "../../components/Service/FindService";

export default function Explore() {
    const [service, setService] = useState(null);

    useEffect(() => {
        if (service) {
            window.location.href = `/services/${service.id}`;
        }
    }, [service]);

    return (
        <div>
            <HeaderBar activeBackground={true}/>
            <h1>Explore</h1>
            <FindService dark={false} setService={setService} />
        </div>
    );
}
