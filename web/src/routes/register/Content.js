/*
** EPITECH PROJECT, 2024
** AREA
** File description:
** Content
*/

import RegisterBox from "../../components/RegisterBox/RegisterBox";

export default function Content( {setToken} ) {
    return (
        <div className="h-screen justify-center items-center flex">
            <RegisterBox setToken={setToken} />
        </div>
    );
}