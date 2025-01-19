/*
** EPITECH PROJECT, 2024
** AREA
** File description:
** Button
*/


import { MoveUpRight } from 'lucide-react';

export default function Button({ text, redirect, onClick, styleClolor, icon }) {
    return (
        <div
            className={`cursor-pointer inline-flex justify-center items-center rounded-3xl ${styleClolor} py-2 px-4 m-2`}
            onClick={onClick}
        >
            {icon ?
                icon :
                <label className="text-center cursor-pointer font-spartan font-bold mt-1">{text}</label>}
            {redirect ? 
                <MoveUpRight className="ml-2" /> : 
                null
            }
        </div>
    );
}

export function LRButton( {text, handleClick, img, color} ) {
    return(
        <div
            className={`
                select-none relative
                 text-white shadow-md flex
                 justify-center items-center cursor-pointer
                 transition-transform duration-200
                 rounded-full p-2 m-2
                 lg:w-[400px] md:w-[300px] sm:w-[200px]
                `}
            onClick={(e) => handleClick(e)}
            style={{ backgroundColor: color }}
            onMouseEnter={(e) => { e.currentTarget.style.transform = "scale(1.05)";}}
            onMouseLeave={(e) => { e.currentTarget.style.transform = "scale(1)"; }}
            tabIndex={0}    
        >
            {img && <img className="lg:w-10 lg:h-10
                                    md:w-5 md:h-5
                                    w-4 h-4"
                src={img}
                alt={text}
            />}
            <label className={`${color === "#ffffffff" ? "text-black" : "text-white"}
                text-base sm:text-sm md:text-lg lg:text-xl 
                font-bold font-spartan ml-2`}
            >
                {text}
            </label>
        </div>
    )
}