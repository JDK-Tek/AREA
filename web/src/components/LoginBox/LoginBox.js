/*
** EPITECH PROJECT, 2024
** AREA
** File description:
** LoginBox
*/

import React, { useState } from "react"
import axios from "axios"

// Login/Register components
import { LRTextFieldsBox } from "../TextFields/TextFields"
import { LRButton } from "../Button"
import LRBox from "../Box/Box"

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

export default function LoginBox ( {setToken, setError} ) {

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
        axios.post(`${backendUrl}/api/login`, {
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
            setError(error.response.data.error);
        });
    }

    const handleOauth = (service) => {
        axios.get(`${backendUrl}/api/oauth/${service}`)
            .then((response) => {
                const oauthWindow = window.open(response.data, "_blank");

                const handleMessage = (event) => {
                    if (event.origin !== window.location.origin) {
                        return;
                    }
                    const code = event.data;
                    if (code !== null) {
                        oauthWindow.close();
                        window.removeEventListener('message', handleMessage);
                        axios.post(`${backendUrl}/api/oauth/${service}`, {
                            code: code
                        }, {
                            headers: {
                                "Content-Type": "application/json"
                            }
                        })
                        .then((response) => {
                            setToken(response.data.token);
                            window.location.href = "/";
                        })
                        .catch((error) => {
                            console.error('Error:', error);
                        });
                    }
                };
    
                window.addEventListener('message', handleMessage);
            })
            .catch((error) => {
                console.error('Error:', error);
            });
    };

    return (
        <LRBox>
            <LoginTexts />
            <LRTextFieldsBox text1="Email" text2="Password" handleChangeField={handleChange}/>
            <div className="text-center pt-8 sm:pt-10 text-white text-sm sm:text-base md:text-lg">
                You don't have an account ? 
                <a href="/register" className="font-bold text-white dark:text-white hover:underline"> Register here!</a>
            </div>
            <div className="text-center pt-4 items-center">
                <LRButton text="Login" handleClick={handleSubmit}/>
                    <div className=" justify-center pt-4 overflow-y-auto max-h-[100px]      ">
                        <LRButton color="#5865F2" img="/assets/services/discord.webp" text="Connect with Discord" handleClick={() => handleOauth("discord")} />
                        <LRButton color="#ff4500" img="/assets/services/reddit.webp"  text="Connect with Reddit" handleClick={() => handleOauth("reddit")} />
                        <LRButton color="#24292e" img="/assets/services/github.webp"  text="Connect with Github" handleClick={() => handleOauth("github")} />
                        <LRButton color="#1DB954" img="/assets/services/spotify.png" text="Connect with Spotify" handleClick={() => handleOauth("spotify")} />
                    </div>
            </div>
        </LRBox>
    )
}
