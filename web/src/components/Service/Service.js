/*
** EPITECH PROJECT, 2024
** AREA
** File description:
** Service
*/

export default function Service({ service }) {
    return (
        <div
            className={`relative w-[200px] h-[150px] text-white rounded-2xl shadow-md p-6 flex flex-col justify-between items-center cursor-pointer`}
            onClick={() => window.location.href = service.link}
            style={{
                backgroundColor: service.color.normal
            }}
            onMouseEnter={(e) => e.currentTarget.style.backgroundColor = service.color.hover}
            onMouseLeave={(e) => e.currentTarget.style.backgroundColor = service.color.normal}
        >
            <img 
                className="w-[50px] h-[50px]" 
                src={service.logo} 
                alt={service.name} 
            />

            <h1 className="font-spartan text-2xl font-bold text-center"> {service.name} </h1>
        </div>
    );
}
