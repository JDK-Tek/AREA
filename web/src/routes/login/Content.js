/*
** EPITECH PROJECT, 2024
** AREA
** File description:
** Content
*/

import LoginBox from "../../components/LoginBox/LoginBox";

export default function Content( {setToken} ) {

    return (
        <div className="h-screen justify-center items-center flex">
            <LoginBox setToken={setToken} />
        </div>
    );
}