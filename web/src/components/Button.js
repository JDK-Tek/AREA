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

export function LRButton( {text, handleClick} ) {
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