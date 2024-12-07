/*
** EPITECH PROJECT, 2024
** AREA
** File description:
** Button
*/


import { MoveUpRight } from 'lucide-react';

export default function Button({ text, redirect, onClick, styleClolor }) {
    return (
        <div
            className={`cursor-pointer inline-flex justify-center items-center rounded-3xl ${styleClolor} py-2 px-4 m-2`}
            onClick={onClick}
        >
            <label className="text-center cursor-pointer font-spartan font-bold mt-1">{text}</label>
            {redirect ? 
                <MoveUpRight className="ml-2" /> : 
                null
            }
        </div>
    );
}
