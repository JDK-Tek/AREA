/*
** EPITECH PROJECT, 2024
** AREA
** File description:
** RegisterBox
*/

import React, { useState } from "react"
import axios from "axios"

import { LoginTextField } from "../LoginBox/LoginBox"
import { LoginTextFieldsBox } from "../LoginBox/LoginBox"
import { Button } from "../LoginBox/LoginBox"

function RegisterTexts() {
    return (
        <div className="text-center">
            <p className="text-white text-4xl sm:text-5xl md:text-6xl lg:text-7xl font-spartan font-bold">
                Register
            </p>
            <p className="text-violet-600 font-bold text-xl sm:text-2xl md:text-3xl">
                Welcome to AREA !
            </p>
        </div>
    )
}

export default function RegisterBox () {

    const [email, setEmail] = useState("")
    const [password, setPassword] = useState("")

    const handleChange = (e) => {
        console.log(e.target.id)
        if (e.target.id === "email") {
            setEmail(e.target.value)
        }
        if (e.target.id === "password") {
            setPassword(e.target.value)
        }
    }

    const handleSubmit = (e) => {
        e.preventDefault();
        axios.post("http://localhost:42000/api/register", {
            email: email,
            password: password
        }, {
            headers: {
                "Content-Type": "application/json"
            }
        })
        .then((response) => {
            console.log('Success:', response.data);
            window.location.href = "/login";
        })
        .catch((error) => {
            console.error('Error:', error);
        });
    }

    return (
        <div className="bg-gradient-to-b from-zinc-700 to-gray-800 flex flex-col justify-center 
                        w-3/4 sm:w-3/4 md:w-2/3 lg:w-1/2 xl:w-2/3 
                        h-4/6 sm:h-3/4 md:h-2/3 lg:h-3/4 rounded-md">
            <RegisterTexts />
            <LoginTextFieldsBox text1="Email" text2="Password" handleChangeField={handleChange}/>
            <div className="text-center pt-8 sm:pt-10 text-white text-sm sm:text-base md:text-lg">
                You already have an account? 
                <a href="/login" className="font-bold text-white dark:text-white hover:underline"> Login here!</a>
            </div>
            <Button text="Register" handleClick={handleSubmit}/>    
        </div>
    )
}
