/*
** EPITECH PROJECT, 2024
** AREA
** File description:
** Home
*/

import { useState } from "react";

import HeaderBar from "./../../components/Header/HeaderBar";
import Footer from "../home/Footer";
import Content from './Content';
import Notification from "../../components/Notification";

export default function Register( {setToken} ) {
    const [error, setError] = useState(null);

    return (
        <div className="bg-gradient-to-br from-zinc-900 via-indigo-900 to-violet-900 animate-gradient bg-300%">
            {error && <Notification error={true} msg={error} setError={setError} />}
            <HeaderBar />
            <Content setToken={setToken} setError={setError}/>
            <Footer />
        </div>
    );
}