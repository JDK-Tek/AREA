/*
** EPITECH PROJECT, 2024
** AREA
** File description:
** Home
*/

import HeaderBar from "./../../components/Header/HeaderBar";
import Footer from "../home/Footer";
import Content from './Content';

export default function Register( {setToken} ) {

    return (
        <div className="bg-gradient-to-br from-zinc-900 via-indigo-900 to-violet-900 animate-gradient bg-300%">
            <HeaderBar />
            <Content setToken={setToken} />
            <Footer />
        </div>
    );
}