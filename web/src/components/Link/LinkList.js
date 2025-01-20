/*
** EPITECH PROJECT, 2024
** AREA
** File description:
** LinkList
*/

export default function LinkList({ title, links }) {

    return (
        <div className="max-w-[500px] pl-20 pr-20 pt-5 pb-5">
            <label className="text-[30px] font-spartan font-bold block mb-7 text-white">{title}</label>
            {links.map((link, index) => (
                <label
                    key={index}
                    className="cursor-pointer block mb-4 text-chartgray-100 hover:text-white"
                    onClick={() => { window.location.href = link.url }}
                    onKeyDown={(event) => { if (event.key === "Enter") { window.location.href = link.url } }}
                    tabIndex={0}
                >
                    {link.title}
                </label>
            ))}
        </div>
    );
}
