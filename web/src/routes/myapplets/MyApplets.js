/*
** EPITECH PROJECT, 2025
** AREA
** File description:
** MyApplets
*/

import axios from "axios";
import Footer from "../home/Footer";
import { useEffect, useState } from "react";
import HeaderBar from "../../components/Header/HeaderBar";
import Notification from "../../components/Notification";
import AppletKit from "../../components/Applet/AppletKit";
import { backendUrl } from "../../App";

export default function MyApplets() {
    const [applets, setApplets] = useState([]);
    const [error, setError] = useState(false);

    const checkConnection = () => {
        axios.get(`${backendUrl}/api/doctor`, {
            headers: {
                "Content-Type": "application/json",
                "Authorization": `Bearer ${sessionStorage.getItem("token")}`,
            }
        })
        .then((res) => {
            if (!res.data.authentificated) {
                window.location.href = "/login";
                sessionStorage.removeItem("token");
            }
        })
        .catch((err) => {
            setError("Impossible to check the authentification: " + err);
        });
    }

    checkConnection();

    useEffect(() => {
        axios
            .get(`${backendUrl}/api/area`, {
                headers: {
                    "Content-Type": "application/json",
                    "Authorization": `Bearer ${sessionStorage.getItem("token")}`
                },
            })
            .then(response => {
                setApplets(response.data);
            })
            .catch(error => {
                setError("Error when trying to get applets: " + error);
            });
    }, [setApplets]);

    console.log(applets);
    return (
        <div>
            <HeaderBar activeBackground={true}/>
            {error && <Notification error={true} setError={setError} msg={error}/>}
            <h1>MyApplets</h1>
            <AppletKit applets={applets}/>
            <Footer/>
        </div>
    );

}