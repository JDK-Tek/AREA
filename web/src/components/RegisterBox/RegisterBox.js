/*
** EPITECH PROJECT, 2024
** AREA
** File description:
** RegisterBox
*/

import React, { useState } from "react"
import axios from "axios"

// Login/Register components
import { LRTextFieldsBox } from "../TextFields/TextFields"
import { LRButton } from "../Button"
import LRBox from "../Box/Box"

function RegisterTexts() {
    return (
        <div className="text-center">
            <p className="text-white text-4xl sm:text-5xl md:text-6xl lg:text-7xl font-spartan font-bold">
                REGISTER
            </p>
            <p className="text-violet-600 font-bold text-xl sm:text-2xl md:text-3xl">
                Welcome to AREA !
            </p>
        </div>
    )
}

export default function RegisterBox ( {setToken} ) {

    const backendUrl = process.env.REACT_APP_BACKEND_URL

    const [email, setEmail] = useState("")
    const [password, setPassword] = useState("")

    const handleChange = (e) => {
        if (e.target.id === "email") {
            setEmail(e.target.value)
        }
        if (e.target.id === "password") {
            setPassword(e.target.value)
        }
    }

    const handleSubmit = (e) => {
        e.preventDefault();
        axios.post(`${backendUrl}/api/register`, {
            email: email,
            password: password
        }, {
            headers: {
                "Content-Type": "application/json"
            }
        })
        .then((response) => {
            setToken(response.data.token)
            window.location.href = "/";
        })
        .catch((error) => {
            console.error('Error:', error);
        });
    }

    return (
        <LRBox>
            <RegisterTexts />
            <LRTextFieldsBox text1="Email" text2="Password" handleChangeField={handleChange}/>
            <div className="text-center pt-8 sm:pt-10 text-white text-sm sm:text-base md:text-lg">
                You already have an account? 
                <a href="/login" className="font-bold text-white dark:text-white hover:underline"> Login here!</a>
            </div>
            <LRButton text="Register" handleClick={handleSubmit}/>    
        </LRBox>
    )
}
