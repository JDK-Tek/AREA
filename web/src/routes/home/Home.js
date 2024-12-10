/*
** EPITECH PROJECT, 2024
** AREA
** File description:
** Home
*/

import { useState } from "react";

import Header from "./../../components/Header/Header";
import Content from "./Content";
import Footer from "./Footer";
import Notification from "./../../components/Notification";


export default function Home() {
    const [error, setError] = useState(null);

    return (
        <div>
            {error && <Notification error={true} msg={error} setError={setError} />}
            <Header />
            <Content setError={setError}/>
            <Footer />
        </div>
    );
}
