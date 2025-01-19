/*
** EPITECH PROJECT, 2024
** AREA
** File description:
** Content
*/

import RegisterBox from "../../components/RegisterBox/RegisterBox";

export default function Content( {setToken, setError} ) {
    return (
        <div className="mt-[50px] mb-[150px] h-screen justify-center items-center flex">
            <RegisterBox setToken={setToken} setError={setError} />
        </div>
    );
}