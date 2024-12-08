/*
** EPITECH PROJECT, 2024
** AREA
** File description:
** LoginBox
*/

import React, { useState } from "react"
import axios from "axios"

function LoginTexts() {
    return (
        <div className="text-center">
            <p className="text-white text-4xl sm:text-5xl md:text-6xl lg:text-7xl font-spartan font-bold">
                LOGIN
            </p>
            <p className="text-violet-600 font-bold text-xl sm:text-2xl md:text-3xl">
                Nice to see you again
            </p>
        </div>
    )
}

export function LoginTextField ( {text, id, handleChangeField} ) {
    return (
        <div className="pt-5 justify-center flex">
            <input type={id} id={id} className="bg-gray-500 border border-gray-700 text-white 
                                          text-lg sm:text-xl md:text-2xl 
                                          w-11/12 sm:w-4/5 md:w-3/4 lg:w-2/3 
                                          rounded-full focus:ring-blue-500 focus:border-blue-500 block p-3 sm:p-4" 
                   placeholder={text} required onChange={handleChangeField}/>
        </div>
    )
}

export function LoginTextFieldsBox( {text1, text2, handleChangeField} ) {
    return (
        <div className="pt-10">
            <LoginTextField text={text1} id="email" handleChangeField={handleChangeField}/>
            <LoginTextField text={text2} id="password" handleChangeField={handleChangeField}/>
        </div>
    )
}

export function Button( {text, handleClick} ) {
    return(
        <div className="flex justify-center pt-10">
            <button className="bg-white hover:bg-gray-300 text-black 
                              text-base sm:text-lg md:text-xl lg:text-2xl 
                              font-bold py-2 sm:py-3 px-8 sm:px-10 rounded-full" onClick={handleClick}>
                {text}
            </button>
        </div>
    )
}

export default function LoginBox () {

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
        axios.post("http://localhost:42000/api/login", {
            email: email,
            password: password
        }, {
            headers: {
                "Content-Type": "application/json"
            }
        })
        .then((response) => {
            console.log('Success:', response.data);
            window.location.href = "/";
        })
        .catch((error) => {
            console.error('Error:', error);
        });
    }

    return (
        <div className="bg-gradient-to-b from-zinc-700 to-gray-800 flex flex-col justify-center 
                        w-3/4 sm:w-3/4 md:w-2/3 lg:w-1/2 xl:w-2/3 
                        h-4/6 sm:h-3/4 md:h-2/3 lg:h-3/4 rounded-md">
            <LoginTexts />
            <LoginTextFieldsBox text1="Email" text2="Password" handleChangeField={handleChange}/>
            <div className="text-center pt-8 sm:pt-10 text-white text-sm sm:text-base md:text-lg">
                You don't have an account ? 
                <a href="/register" className="font-bold text-white dark:text-white hover:underline"> Register here!</a>
            </div>
            <Button text="Login" handleClick={handleSubmit}/>
            <Button text="Connect with Discord" />  
        </div>
    )
}
