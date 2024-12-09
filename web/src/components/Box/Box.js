/*
** EPITECH PROJECT, 2024
** AREA
** File description:
** Box
*/

import React from 'react';

const LRBox = ({ children }) => {
  return (
    <div className="bg-gradient-to-b from-zinc-700 to-gray-800 flex flex-col justify-center 
                    w-3/4 sm:w-3/4 md:w-2/3 lg:w-1/2 xl:w-2/3 
                    h-4/6 sm:h-3/4 md:h-2/3 lg:h-3/4 rounded-md">
      {children}
    </div>
  );
};

export default LRBox;