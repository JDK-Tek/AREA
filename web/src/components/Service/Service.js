/*
** EPITECH PROJECT, 2024
** AREA
** File description:
** Service
*/

export default function Service({ service, rounded }) {
    return (
        <div
            className={`select-none relative w-[200px] h-[150px] text-white ${rounded} shadow-md p-6 flex flex-col justify-between items-center cursor-pointer transition-transform duration-200`}
            onClick={() => window.location.href = service.link}
            style={{
                backgroundColor: service.color.normal
            }}
            onMouseEnter={(e) => {
                e.currentTarget.style.backgroundColor = service.color.hover;
                e.currentTarget.style.transform = "scale(1.05)";
            }}
            onMouseLeave={(e) => {
                e.currentTarget.style.backgroundColor = service.color.normal;
                e.currentTarget.style.transform = "scale(1)";
            }}
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

