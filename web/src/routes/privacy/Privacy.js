/*
** EPITECH PROJECT, 2024
** AREA
** File description:
** Privacy
*/

import React, { useEffect, useState } from 'react';
import ReactMarkdown from 'react-markdown';
import HeaderBar from "./../../components/Header/HeaderBar";
import Logo from "./../../assets/logo.png";
import file from './PrivacyPolicy.md';

export default function Privacy() {

    const [markdown, setMarkdown] = useState('');

    useEffect(() => {
        fetch(file)
            .then((response) => response.text())
            .then((text) => {
                setMarkdown(text);
            });
    }, []);

    return (
        <div>
            <HeaderBar activeBackground={true}/>
            <div className="min-h-screen flex flex-col items-center justify-center">
                <div className="markdown-content">
                    <ReactMarkdown>{markdown}</ReactMarkdown>
                </div>
            </div>
        </div>
    );
}