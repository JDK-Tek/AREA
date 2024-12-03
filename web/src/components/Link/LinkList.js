/*
** EPITECH PROJECT, 2024
** AREA
** File description:
** LinkList
*/

export default function LinkList({ title, links }) {

    return (
        <div className="pl-20 pr-20 pt-5 pb-5">
            <label className="text-[30px] font-spartan font-bold block mb-7 text-white">{title}</label>
            {links.map((link, index) => (
                <label
                    key={index}
                    className="cursor-pointer block mb-4 text-chartgray-100 hover:text-white"
                    onClick={() => { window.location.href = link.url }}
                >
                    {link.title}
                </label>
            ))}
        </div>
    );
}
