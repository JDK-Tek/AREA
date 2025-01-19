/*
** EPITECH PROJECT, 2024
** AREA
** File description:
** Content
*/

import LoginBox from "../../components/LoginBox/LoginBox";

export default function Content( {setToken, setError} ) {

    return (
        <div className="mt-[50px] mb-[150px] h-screen justify-center items-center flex">
            <LoginBox setToken={setToken} setError={setError} />
        </div>
    );
}